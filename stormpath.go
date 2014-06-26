package stormpath

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/jarias/stormpath/logger"
	"github.com/nu7hatch/gouuid"
)

const (
	DEFAULT_ALGORITHM      = "SHA256"
	HOST_HEADER            = "Host"
	AUTHORIZATION_HEADER   = "Authorization"
	STORMPATH_DATE_HEADER  = "X-Stormpath-Date"
	ID_TERMINATOR          = "sauthc1_request"
	ALGORITHM              = "HMAC-SHA-256"
	AUTHENTICATION_SCHEME  = "SAuthc1"
	SAUTHC1_ID             = "sauthc1Id"
	SAUTHC1_SIGNED_HEADERS = "sauthc1SignedHeaders"
	SAUTHC1_SIGNATURE      = "sauthc1Signature"
	DATE_FORMAT            = "20060102"
	TIMESTAMP_FORMAT       = "20060102T150405Z0700"
	NL                     = "\n"
)

var (
	Client *StormpathClient
)

type StormpathClient struct {
	Credentials *Credentials
	HttpClient  *http.Client
}

func NewStormpathClient(credentials *Credentials) *StormpathClient {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{},
		DisableCompression: true,
	}
	httpClient := &http.Client{Transport: tr}

	return &StormpathClient{Credentials: credentials, HttpClient: httpClient}
}

func (client *StormpathClient) DoWithResult(request *StormpathRequest, result interface{}) error {
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	return unmarshal(resp, result)
}

func (client *StormpathClient) Do(request *StormpathRequest) (resp *http.Response, err error) {
	req, err := request.ToHttpRequest()

	if err != nil {
		return nil, err
	}

	uuid, _ := uuid.NewV4()
	nonce := uuid.String()

	Authenticate(req, request.marshalPayload(), time.Now().In(time.UTC), client.Credentials, nonce)

	if !request.DontFollowRedirects {
		logger.INFO.Printf("Executing request [%s] following redirects", req.URL)
		resp, err := client.HttpClient.Do(req)

		if err != nil {
			return resp, err
		}

		return resp, handleStormpathErrors(resp)
	} else {
		logger.INFO.Printf("Executing request [%s] without following redirects", req.URL)
		resp, err := client.HttpClient.Transport.RoundTrip(req)

		if err != nil {
			return resp, err
		}

		return resp, handleStormpathErrors(resp)
	}
}

func Authenticate(req *http.Request, payload []byte, date time.Time, credentials *Credentials, nonce string) {
	timestamp := date.Format(TIMESTAMP_FORMAT)
	dateStamp := date.Format(DATE_FORMAT)
	req.Header.Set(HOST_HEADER, req.URL.Host)
	req.Header.Set(STORMPATH_DATE_HEADER, timestamp)

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

	id := credentials.Id + "/" + dateStamp + "/" + nonce + "/" + ID_TERMINATOR

	canonicalRequestHashHex := hex.EncodeToString(hash([]byte(canonicalRequest)))

	stringToSign :=
		ALGORITHM +
			NL +
			timestamp +
			NL +
			id +
			NL +
			canonicalRequestHashHex

	kSecret := []byte(AUTHENTICATION_SCHEME + credentials.Secret)
	kDate := sing(dateStamp, kSecret)
	kNonce := sing(nonce, kDate)
	kSigning := sing(ID_TERMINATOR, kNonce)

	signature := sing(stringToSign, kSigning)
	signatureHex := hex.EncodeToString(signature)

	authorizationHeader :=
		AUTHENTICATION_SCHEME + " " +
			createNameValuePair(SAUTHC1_ID, id) + ", " +
			createNameValuePair(SAUTHC1_SIGNED_HEADERS, signedHeadersString) + ", " +
			createNameValuePair(SAUTHC1_SIGNATURE, signatureHex)

	req.Header.Set(AUTHORIZATION_HEADER, authorizationHeader)
}

func createNameValuePair(name string, value string) string {
	return name + "=" + value
}

func encodeUrl(value string, path bool, canonical bool) string {
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
	var keys []string
	queryValues := req.URL.Query()

	result := ""
	for k, _ := range queryValues {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		key := encodeUrl(k, false, true)
		v := queryValues[k]
		for _, vv := range v {
			value := encodeUrl(vv, false, true)

			if len(result) > 0 {
				result = result + "&"
			}

			result = result + key + "=" + value
		}
	}

	return result
}

func canonicalizeResourcePath(path string) string {
	if len(path) == 0 {
		return "/"
	} else {
		return encodeUrl(path, true, true)
	}
}

func canonicalizeHeadersString(headers http.Header) string {
	var keys []string
	stringBuffer := bytes.NewBufferString("")

	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		stringBuffer.Write([]byte(strings.ToLower(k)))
		stringBuffer.Write([]byte(":"))

		first := true

		for _, v := range headers[k] {
			if !first {
				stringBuffer.Write([]byte(","))
			}
			stringBuffer.Write([]byte(v))
			first = false
		}
		stringBuffer.Write([]byte(NL))
	}

	return stringBuffer.String()
}

func signedHeadersString(headers http.Header) string {
	var keys []string
	stringBuffer := bytes.NewBufferString("")

	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	first := true
	for _, k := range keys {
		if !first {
			stringBuffer.Write([]byte(";"))
		}
		stringBuffer.Write([]byte(strings.ToLower(k)))
		first = false
	}

	return stringBuffer.String()
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
