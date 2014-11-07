package stormpath

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"time"

	"github.com/nu7hatch/gouuid"
)

//Version is the current SDK Version
const version = "0.1.0-beta.2"
const baseUrl = "https://api.stormpath.com/v1/"
const followRedirectsHeader = "Stormpath-Go-FollowRedirects"
const locationHeader = "Location"

//Client is default global client variable to execute any Stormpath request
var Client *StormpathClient

//StormpathClient is low level REST client for any Stormpath request,
//it holds the credentials, an the actual http client, and the cache.
//The Cache can be initialize in nil and the client would simply ignore it
//and don't cache any response.
type StormpathClient struct {
	Credentials *Credentials
	HTTPClient  *http.Client
	Cache       Cache
}

//List defines the paged result metadata such as offset and limit
type list struct {
	Href   string `json:"href"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

//Link defines href field in any of the data models this struct is meant to be embeded into other models
type link struct {
	Href string `json:"href"`
}

type stormpathError struct {
	Status           int
	Code             int
	Message          string
	DeveloperMessage string
	MoreInfo         string
}

//NewStormpathClient is a convience constructor for the StormpathClient struct,
//it recieves a pointer to a credentials object a cache implementation and
//returns a pointer to a StormpathClient object, the cache implementation can be nil
func NewStormpathClient(credentials *Credentials, cache Cache) *StormpathClient {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{},
		DisableCompression: true,
	}
	httpClient := &http.Client{Transport: tr}

	return &StormpathClient{credentials, httpClient, cache}
}

func buildRelativeURL(parts ...string) string {
	buffer := bytes.NewBufferString(baseUrl)

	for i, part := range parts {
		buffer.WriteString(part)
		if i+1 < len(parts) {
			buffer.WriteString("/")
		}
	}

	return buffer.String()
}

func buildAbsoluteURL(parts ...string) string {
	buffer := bytes.NewBufferString("")

	for i, part := range parts {
		buffer.WriteString(part)
		if i+1 < len(parts) {
			buffer.WriteString("/")
		}
	}

	return buffer.String()
}

func (client *StormpathClient) newRequestWithoutRedirects(method string, urlStr string, body interface{}) *http.Request {
	req := client.newRequest(method, urlStr, body)
	req.Header.Add(followRedirectsHeader, "false")
	return req
}

func (client *StormpathClient) newRequest(method string, urlStr string, body interface{}) *http.Request {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, urlStr, bytes.NewReader(jsonBody))
	req.Header.Set("User-Agent", "jarias/stormpath-sdk-go/"+version)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	uuid, _ := uuid.NewV4()
	nonce := uuid.String()

	Authenticate(req, jsonBody, time.Now().In(time.UTC), client.Credentials, nonce)
	return req
}

func requestParams(pageRequest url.Values, filter url.Values, extraParams url.Values) string {
	params := url.Values{}

	params = appendParams(params, pageRequest)
	params = appendParams(params, filter)
	params = appendParams(params, extraParams)

	encodedParams := params.Encode()
	if encodedParams != "" {
		return "?" + encodedParams
	}
	return ""
}

func appendParams(params url.Values, toAppend url.Values) url.Values {
	for k, v := range toAppend {
		params[k] = v
	}
	return params
}

func handleResponseError(resp *http.Response, err error) error {
	//Error from the request execution
	if err != nil {
		ERROR.Printf("%s [%s]", err, resp.Request.URL.String())
		return err
	}
	//Check for Stormpath specific errors
	if resp.StatusCode != 200 && resp.StatusCode != 204 && resp.StatusCode != 201 && resp.StatusCode != 302 {
		spError := &stormpathError{}

		err := json.NewDecoder(resp.Body).Decode(spError)
		if err != nil {
			return err
		}

		ERROR.Printf("%s [%s]", spError.Message, resp.Request.URL.String())
		return errors.New(spError.Message)
	}
	//No errors from the request execution
	return nil
}

func emptyPayload() []byte {
	return []byte{}
}

//doWithResult executes the given StormpathRequest and serialize the response body into the given expected result,
//it returns an error if any occurred while executing the request or serializing the response
func (client *StormpathClient) doWithResult(request *http.Request, result interface{}) error {
	var err error
	var response *http.Response

	key := request.URL.String()

	if client.Cache != nil && request.Method == "GET" && client.Cache.Exists(key) {
		err = client.Cache.Get(key, result)
	} else {
		response, err = client.execRequest(request)
		if err != nil {
			return err
		}
		err = json.NewDecoder(response.Body).Decode(result)
	}

	if client.Cache != nil && err == nil {
		switch request.Method {
		case "POST", "DELETE", "PUT":
			client.Cache.Del(key)
			break
		case "GET":
			client.Cache.Set(key, result)
		}
	}

	return err
}

//do executes the StormpathRequest without expecting a response body as a result,
//it returns an error if any occurred while executing the request
func (client *StormpathClient) do(request *http.Request) error {
	_, err := client.execRequest(request)
	return err
}

//execRequest executes a request, it would return a byte slice with the raw resoponse data and an error if any occurred
func (client *StormpathClient) execRequest(req *http.Request) (*http.Response, error) {
	if req.Header.Get(followRedirectsHeader) == "false" {
		req.Header.Del(followRedirectsHeader)
		resp, err := client.HTTPClient.Transport.RoundTrip(req)
		err = handleResponseError(resp, err)
		if err != nil {
			ERROR.Printf("%s [%s]", err, resp.Request.URL.String())
			return nil, err
		}
		//Get the redirect location from the response headers
		newReq := client.newRequest("GET", resp.Header.Get(locationHeader), emptyPayload())
		return client.execRequest(newReq)
	}

	resp, err := client.HTTPClient.Do(req)
	return resp, handleResponseError(resp, err)
}
