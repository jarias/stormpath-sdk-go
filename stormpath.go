package stormpath

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

//BaseURL defines the Stormpath API base URL
var BaseURL = "https://api.stormpath.com/v1/"

//Version is the current SDK Version
const version = "0.1.0-beta.8"

var client *Client

//Client is low level REST client for any Stormpath request,
//it holds the credentials, an the actual http client, and the cache.
//The Cache can be initialize in nil and the client would simply ignore it
//and don't cache any response.
type Client struct {
	Credentials Credentials
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

//Init initializes the underlying client that communicates with Stormpath
func Init(credentials Credentials, cache Cache) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{},
		DisableCompression: true,
	}
	httpClient := &http.Client{Transport: tr}
	httpClient.CheckRedirect = checkRedirect

	client = &Client{credentials, httpClient, cache}

	initLog()
}

//InitWithCustomHTTPClient initializes the underlying client that communicates with Stormpath with a custom http.Client
func InitWithCustomHTTPClient(credentials Credentials, cache Cache, httpClient *http.Client) {
	httpClient.CheckRedirect = checkRedirect
	client = &Client{credentials, httpClient, cache}

	initLog()
}

func (client *Client) post(urlStr string, body interface{}, result interface{}) error {
	return client.execute("POST", urlStr, body, result)
}

func (client *Client) get(urlStr string, body interface{}, result interface{}) error {
	return client.execute("GET", urlStr, body, result)
}

func (client *Client) delete(urlStr string, body interface{}) error {
	return client.do(client.newRequest("DELETE", urlStr, body))
}

func (client *Client) execute(method string, urlStr string, body interface{}, result interface{}) error {
	return client.doWithResult(client.newRequest(method, urlStr, body), result)
}

func buildRelativeURL(parts ...string) string {
	buffer := bytes.NewBufferString(BaseURL)

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

func (client *Client) newRequest(method string, urlStr string, body interface{}) *http.Request {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, urlStr, bytes.NewReader(jsonBody))
	req.Header.Set("User-Agent", fmt.Sprintf("jarias/stormpath-sdk-go/%s (%s; %s)", version, runtime.GOOS, runtime.GOARCH))
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
		Logger.Printf("[ERROR] %s [%s]", err, resp.Request.URL.String())
		return err
	}
	//Check for Stormpath specific errors
	if resp.StatusCode != 200 && resp.StatusCode != 204 && resp.StatusCode != 201 && resp.StatusCode != 302 {
		spError := &stormpathError{}

		err := json.NewDecoder(resp.Body).Decode(spError)
		if err != nil {
			return err
		}

		Logger.Printf("[ERROR] %s", spError)
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
func (client *Client) doWithResult(request *http.Request, result interface{}) error {
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
func (client *Client) do(request *http.Request) error {
	_, err := client.execRequest(request)
	return err
}

//execRequest executes a request, it would return a byte slice with the raw resoponse data and an error if any occurred
func (client *Client) execRequest(req *http.Request) (*http.Response, error) {
	if logLevel == "DEBUG" {
		//Print request
		dump, _ := httputil.DumpRequest(req, true)
		Logger.Printf("[DEBUG] Stormpath request\n%s", dump)
	}
	resp, err := client.HTTPClient.Do(req)
	if logLevel == "DEBUG" {
		//Print response
		dump, _ := httputil.DumpResponse(resp, true)
		Logger.Printf("[DEBUG] Stormpath response\n%s", dump)
	}
	return resp, handleResponseError(resp, err)
}

func checkRedirect(req *http.Request, via []*http.Request) error {
	//Go client defautl behavior is to bail after 10 redirects
	if len(via) > 10 {
		return errors.New("stopped after 10 redirects")
	}
	//No redirect do nothing
	if len(via) == 0 {
		// No redirects
		return nil
	}
	// Re-Authenticate the redirect request
	uuid, _ := uuid.NewV4()
	nonce := uuid.String()

	//We can use an empty payload cause the only redirect is for the current tenant
	//this could change in the future
	Authenticate(req, emptyPayload(), time.Now().In(time.UTC), client.Credentials, nonce)

	return nil
}

func (e stormpathError) String() string {
	return fmt.Sprintf("Stormpath request error \nCode: [ %d ]\nMessage: [ %s ]\nDeveloper Message: [ %s ]\nMore info [ %s ]", e.Code, e.Message, e.DeveloperMessage, e.MoreInfo)
}
