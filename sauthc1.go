package stormpath

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

//SAuthc1 algorithm constants
const (
	IDTerminator         = "sauthc1_request"
	AuthenticationScheme = "SAuthc1"
	NL                   = "\n"
	HostHeader           = "Host"
	AuthorizationHeader  = "Authorization"
	StormpathDateHeader  = "X-Stormpath-Date"
	Algorithm            = "HMAC-SHA-256"
	SAUTHC1Id            = "sauthc1Id"
	SAUTHC1SignedHeaders = "sauthc1SignedHeaders"
	SAUTHC1Signature     = "sauthc1Signature"
	DateFormat           = "20060102"
	TimestampFormat      = "20060102T150405Z0700"
)

//Authenticate generates the proper authentication header for the SAuthc1 algorithm use by Stormpath
func Authenticate(req *http.Request, payload []byte, date time.Time, credentials Credentials, nonce string) {
	timestamp := date.Format(TimestampFormat)
	dateStamp := date.Format(DateFormat)
	req.Header.Set(HostHeader, req.URL.Host)
	req.Header.Set(StormpathDateHeader, timestamp)

	signedHeadersString := signedHeadersString(req.Header)

	canonicalRequest :=
		req.Method +
			NL +
			canonicalizeResourcePath(req.URL.Path) +
			NL +
			canonicalizeQueryString(req) +
			NL +
			canonicalizeHeadersString(req.Header) +
			NL +
			signedHeadersString +
			NL +
			hex.EncodeToString(hash(payload))

	id := credentials.ID + "/" + dateStamp + "/" + nonce + "/" + IDTerminator

	canonicalRequestHashHex := hex.EncodeToString(hash([]byte(canonicalRequest)))

	stringToSign :=
		Algorithm +
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
			createNameValuePair(SAUTHC1Id, id) + ", " +
			createNameValuePair(SAUTHC1SignedHeaders, signedHeadersString) + ", " +
			createNameValuePair(SAUTHC1Signature, signatureHex)

	req.Header.Set(AuthorizationHeader, authorizationHeader)
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

	for _, k := range keys {
		if stringBuffer.Len() > 0 {
			stringBuffer.WriteString(";")
		}
		stringBuffer.WriteString(strings.ToLower(k))
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
