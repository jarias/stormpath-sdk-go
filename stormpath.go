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
const version = "0.1.0-beta.11"

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

//buildExpandParam coverts a slice of expand attributes to a url.Values with
//only one value "expand=attr1,attr2,etc"
func buildExpandParam(expandAttributes []string) url.Values {
	stringBuffer := bytes.NewBufferString("")

	first := true
	for _, expandAttribute := range expandAttributes {
		if !first {
			stringBuffer.WriteString(",")
		}
		stringBuffer.WriteString(expandAttribute)
		first = false
	}

	values := url.Values{}
	expandValue := stringBuffer.String()
	//Should not include the expand query param if the value is empty
	if expandValue != "" {
		values.Add("expand", expandValue)
	}

	return values
}

func requestParams(values ...url.Values) string {
	params := url.Values{}

	for _, v := range values {
		params = appendParams(params, v)
	}

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
			cacheResource(key, result, client.Cache)
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

func cleanCustomData(customData map[string]interface{}) map[string]interface{} {
	// delete illegal keys from data
	// http://docs.stormpath.com/rest/product-guide/#custom-data
	keys := []string{
		"href", "createdAt", "modifiedAt", "meta",
		"spMeta", "spmeta", "ionmeta", "ionMeta",
	}

	for i := range keys {
		delete(customData, keys[i])
	}

	return customData
}
