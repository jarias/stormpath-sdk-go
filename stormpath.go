package stormpath

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"github.com/nu7hatch/gouuid"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

//Version is the current SDK Version
const Version = "0.0.1"

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

//DoWithResult executes the given StormpathRequest and serialize the response body into the given expected result,
//it returns an error if any occurred while executing the request or serializing the response
func (client *StormpathClient) DoWithResult(request *StormpathRequest, result interface{}) error {
	var responseData []byte
	var err error
	req, err := request.ToHttpRequest()

	responseData, err = client.execRequestWithCache(req, request.marshalPayload(), request.DontFollowRedirects)

	if err != nil {
		return err
	}
	return unmarshal(responseData, result)
}

//Do executes the StormpathRequest without expecting a response body as a result,
//it returns an error if any occurred while executing the request
func (client *StormpathClient) Do(request *StormpathRequest) error {
	req, err := request.ToHttpRequest()
	if err != nil {
		return err
	}
	_, err = client.execRequest(req, request.marshalPayload(), request.DontFollowRedirects)
	return err
}

//execRequestWithCache executes a request and caches its response if Method == GET and there is a valid cache implementation,
//it would return a byte slice with the raw resoponse data and an error if any occurred
func (client *StormpathClient) execRequestWithCache(req *http.Request, payload []byte, dontfollowRedirects bool) ([]byte, error) {
	var responseData []byte
	var err error

	key := req.URL.String()

	if client.Cache != nil && req.Method == GET && client.Cache.Exists(key) {
		responseData, err = client.Cache.Get(key)
	} else {
		responseData, err = client.execRequest(req, payload, dontfollowRedirects)
	}

	if client.Cache != nil {
		switch req.Method {
		case POST, DELETE, PUT:
			client.Cache.Del(key)
			break
		case GET:
			client.Cache.Set(key, responseData)
		}
	}

	return responseData, err
}

//execRequest executes a request, it would return a byte slice with the raw resoponse data and an error if any occurred
func (client *StormpathClient) execRequest(req *http.Request, payload []byte, dontfollowRedirects bool) ([]byte, error) {
	var resp *http.Response
	var err error

	uuid, _ := uuid.NewV4()
	nonce := uuid.String()

	Authenticate(req, payload, time.Now().In(time.UTC), client.Credentials, nonce)

	if dontfollowRedirects {
		resp, err = client.HTTPClient.Transport.RoundTrip(req)
		if err != nil {
			return []byte{}, err
		}
		//Get the redirect location from the response headers
		req, _ := http.NewRequest(GET, resp.Header.Get(LocationHeader), bytes.NewReader(payload))
		return client.execRequest(req, payload, !dontfollowRedirects)
	} else {
		resp, err = client.HTTPClient.Do(req)
	}

	if err != nil {
		return []byte{}, err
	}
	err = handleStormpathErrors(resp)
	if err != nil {
		return []byte{}, err
	}
	return extractResponseData(resp)
}

//Constants use for the SAuthc1 authentication algorithm
const (
	IDTerminator         = "sauthc1_request"
	AuthenticationScheme = "SAuthc1"
	NL                   = "\n"
)

//Authenticate generates the proper authentication header for the SAuthc1 algorithm use by Stormpath
func Authenticate(req *http.Request, payload []byte, date time.Time, credentials *Credentials, nonce string) {
	timestamp := date.Format("20060102T150405Z0700")
	dateStamp := date.Format("20060102")
	req.Header.Set("Host", req.URL.Host)
	req.Header.Set("X-Stormpath-Date", timestamp)

	canonicalResourcePath := canonicalizeResourcePath(req.URL.Path)
	canonicalQueryString := canonicalizeQueryString(req)
	canonicalHeadersString := canonicalizeHeadersString(req.Header)
	signedHeadersString := signedHeadersString(req.Header)

	requestPayloadHashHex := hex.EncodeToString(hash(payload))

	canonicalRequest :=
		req.Method +
			NL +
			canonicalResourcePath +
			NL +
			canonicalQueryString +
			NL +
			canonicalHeadersString +
			NL +
			signedHeadersString +
			NL +
			requestPayloadHashHex

	id := credentials.Id + "/" + dateStamp + "/" + nonce + "/" + IDTerminator

	canonicalRequestHashHex := hex.EncodeToString(hash([]byte(canonicalRequest)))

	stringToSign :=
		"HMAC-SHA-256" +
			NL +
			timestamp +
			NL +
			id +
			NL +
			canonicalRequestHashHex

	secret := []byte(AuthenticationScheme + credentials.Secret)
	singDate := sing(dateStamp, secret)
	singNonce := sing(nonce, singDate)
	signing := sing(IDTerminator, singNonce)

	signature := sing(stringToSign, signing)
	signatureHex := hex.EncodeToString(signature)

	authorizationHeader :=
		AuthenticationScheme + " " +
			createNameValuePair("sauthc1Id", id) + ", " +
			createNameValuePair("sauthc1SignedHeaders", signedHeadersString) + ", " +
			createNameValuePair("sauthc1Signature", signatureHex)

	req.Header.Set("Authorization", authorizationHeader)
}

func createNameValuePair(name string, value string) string {
	return name + "=" + value
}

func encodeURL(value string, path bool, canonical bool) string {
	if value == "" {
		return ""
	}

	encoded := url.QueryEscape(value)

	if canonical {
		encoded = strings.Replace(encoded, "+", "%20", -1)
		encoded = strings.Replace(encoded, "*", "%2A", -1)
		encoded = strings.Replace(encoded, "%7E", "~", -1)

		if path {
			encoded = strings.Replace(encoded, "%2F", "/", -1)
		}
	}

	return encoded
}

func canonicalizeQueryString(req *http.Request) string {
	stringBuffer := bytes.NewBufferString("")

	queryValues := req.URL.Query()

	keys := sortedMapKeys(queryValues)

	for _, k := range keys {
		key := encodeURL(k, false, true)
		v := queryValues[k]
		for _, vv := range v {
			value := encodeURL(vv, false, true)

			if stringBuffer.Len() > 0 {
				stringBuffer.WriteString("&")
			}

			stringBuffer.WriteString(key + "=" + value)
		}
	}

	return stringBuffer.String()
}

func canonicalizeResourcePath(path string) string {
	if len(path) == 0 {
		return "/"
	}
	return encodeURL(path, true, true)
}

func canonicalizeHeadersString(headers http.Header) string {
	stringBuffer := bytes.NewBufferString("")

	keys := sortedMapKeys(headers)

	for _, k := range keys {
		stringBuffer.WriteString(strings.ToLower(k))
		stringBuffer.WriteString(":")

		first := true

		for _, v := range headers[k] {
			if !first {
				stringBuffer.WriteString(",")
			}
			stringBuffer.WriteString(v)
			first = false
		}
		stringBuffer.WriteString(NL)
	}

	return stringBuffer.String()
}

func signedHeadersString(headers http.Header) string {
	stringBuffer := bytes.NewBufferString("")

	keys := sortedMapKeys(headers)

	first := true
	for _, k := range keys {
		if !first {
			stringBuffer.WriteString(";")
		}
		stringBuffer.WriteString(strings.ToLower(k))
		first = false
	}

	return stringBuffer.String()
}

func sortedMapKeys(m map[string][]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func hash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func sing(data string, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}
