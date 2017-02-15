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
	"strings"
	"sync"
	"time"

	"io/ioutil"

	uuid "github.com/nu7hatch/gouuid"
)

//Version is the current SDK Version
const version = "0.1.0-beta.24"

const (
	Enabled                   = "ENABLED"
	Disabled                  = "DISABLED"
	Unverified                = "UNVERIFIED"
	ApplicationJSON           = "application/json"
	ApplicationFormURLencoded = "application/x-www-form-urlencoded"
	TextPlain                 = "text/plain"
	TextHTML                  = "text/html"
	ContentTypeHeader         = "Content-Type"
	AcceptHeader              = "Accept"
	UserAgentHeader           = "User-Agent"
)

var client *Client
var buffPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

//Client is low level REST client for any Stormpath request,
//it holds the credentials, an the actual http client, and the cache.
//The Cache can be initialize in nil and the client would simply ignore it
//and don't cache any response.
type Client struct {
	ClientConfiguration ClientConfiguration
	HTTPClient          *http.Client
	Cache               Cache
	WebSDKToken         string
}

//Init initializes the underlying client that communicates with Stormpath
func Init(clientConfiguration ClientConfiguration, cache Cache) {
	InitLog()

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{},
		DisableCompression: true,
	}
	httpClient := &http.Client{Transport: tr}
	httpClient.CheckRedirect = checkRedirect

	client = &Client{clientConfiguration, httpClient, nil, ""}

	if clientConfiguration.CacheManagerEnabled && cache == nil {
		client.Cache = NewLocalCache(clientConfiguration.CacheTTL, clientConfiguration.CacheTTI)
	} else if clientConfiguration.CacheManagerEnabled && cache != nil {
		client.Cache = cache
	}
}

//GetClient returns the configured client
func GetClient() *Client {
	return client
}

func (client *Client) postURLEncodedForm(urlStr string, body string, result interface{}) error {
	return client.execute(http.MethodPost, urlStr, []byte(body), result, ApplicationFormURLencoded)
}

func (client *Client) post(urlStr string, body interface{}, result interface{}) error {
	return client.execute(http.MethodPost, urlStr, body, result, ApplicationJSON)
}

func (client *Client) get(urlStr string, result interface{}) error {
	return client.execute(http.MethodGet, urlStr, emptyPayload(), result, ApplicationJSON)
}

func (client *Client) delete(urlStr string) error {
	return client.do(client.newRequest(http.MethodDelete, urlStr, emptyPayload(), ApplicationJSON))
}

func (client *Client) execute(method string, urlStr string, body interface{}, result interface{}, contentType string) error {
	return client.doWithResult(client.newRequest(method, urlStr, body, contentType), result)
}

func buildRelativeURL(parts ...string) string {
	p := append([]string{client.ClientConfiguration.BaseURL}, parts...)
	return buildAbsoluteURL(p...)
}

func buildAbsoluteURL(parts ...string) string {
	buffer := bytes.NewBufferString("")

	for i, part := range parts {
		buffer.WriteString(part)
		if !strings.HasSuffix(part, "/") && i+1 < len(parts) {
			buffer.WriteString("/")
		}
	}

	return buffer.String()
}

func (client *Client) newRequest(method string, urlStr string, body interface{}, contentType string) *http.Request {
	var encodedBody []byte

	if contentType != ApplicationJSON || method == http.MethodGet || method == http.MethodDelete {
		//If content type is not application/json then it is application/x-www-form-urlencoded in which case the body should the encoded params as a []byte
		//Fixes issue #23 if the method is GET then body should also be just the bytes instead of doing a JSON marshaling
		encodedBody = body.([]byte)
	} else {
		if _, ok := body.([]byte); !ok {
			encodedBody, _ = json.Marshal(body)
		}
	}
	req, _ := http.NewRequest(method, urlStr, bytes.NewReader(encodedBody))

	req.Header.Set(UserAgentHeader, strings.TrimSpace(fmt.Sprintf("stormpath-sdk-go/%s %s", version, client.WebSDKToken)))
	req.Header.Set(AcceptHeader, ApplicationJSON)
	req.Header.Set(ContentTypeHeader, contentType)

	uuid, _ := uuid.NewV4()
	nonce := uuid.String()

	Authenticate(req, encodedBody, time.Now().In(time.UTC), client.ClientConfiguration.APIKeyID, client.ClientConfiguration.APIKeySecret, nonce)
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
	buff := buffPool.Get().(*bytes.Buffer)
	buff.Reset()
	defer buffPool.Put(buff)

	first := true
	for _, v := range values {
		encodedValues := v.Encode()

		if buff.Len() > 0 && !first && encodedValues != "" {
			buff.WriteByte('&')
		}
		buff.WriteString(encodedValues)
		first = false
	}

	encodedParams := buff.String()
	if encodedParams != "" {
		return "?" + encodedParams
	}
	return ""
}

func emptyPayload() []byte {
	return []byte{}
}

//doWithResult executes the given StormpathRequest and serialize the response body into the given expected result,
//it returns an error if any occurred while executing the request or serializing the response
func (client *Client) doWithResult(request *http.Request, result interface{}) error {
	var jsonData []byte
	var err error

	key := request.URL.String()

	if client.Cache != nil && request.Method == http.MethodGet && client.Cache.Exists(key) {
		jsonData = client.Cache.Get(key)
	}

	if len(jsonData) == 0 {
		response, err := client.execRequest(request)
		if err != nil {
			return err
		}
		jsonData, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
	}

	if result != nil {
		err = json.NewDecoder(bytes.NewBuffer(jsonData)).Decode(result)
	}

	if client.Cache != nil && err == nil && result != nil {
		switch request.Method {
		case http.MethodPost, http.MethodDelete:
			client.Cache.Del(key)
			break
		case http.MethodGet:
			c, ok := result.(Cacheable)
			if ok &&
				c.IsCacheable() &&
				!strings.Contains(key, "passwordResetTokens") &&
				!strings.Contains(key, "authTokens") {
				client.Cache.Set(key, jsonData)
			}
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
	return resp, handleResponseError(req, resp, err)
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

	//In Go 1.8 the authorization header remains in the redirect request causing auth errors
	req.Header.Del(AuthorizationHeader)

	//We can use an empty payload cause the only redirect is for the current tenant
	//this could change in the future
	Authenticate(req, emptyPayload(), time.Now().In(time.UTC), client.ClientConfiguration.APIKeyID, client.ClientConfiguration.APIKeySecret, nonce)

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
