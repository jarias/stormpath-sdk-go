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
	"sync"
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
	EQ                   = '='
	SPACE                = ' '
	SLASH                = '/'
	AMP                  = '&'
	CS                   = ", "
	COMMA                = ','
	COLON                = ':'
	SemiColon            = ';'
	EMPTY                = ""
)

var sha256Hash = sha256.New()
var buffPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

//Authenticate generates the proper authentication header for the SAuthc1 algorithm use by Stormpath
func Authenticate(req *http.Request, payload []byte, date time.Time, credentials Credentials, nonce string) {
	timestamp := date.Format(TimestampFormat)
	dateStamp := date.Format(DateFormat)
	req.Header.Set(HostHeader, req.URL.Host)
	req.Header.Set(StormpathDateHeader, timestamp)

	sortedHeaderKeys := sortedMapKeys(req.Header)

	signedHeadersString := signedHeadersString(req.Header, sortedHeaderKeys)

	canonicalRequest := buildCanonicalRequest(req, payload, signedHeadersString, sortedHeaderKeys)

	id := buildID(nonce, dateStamp, credentials)

	stringToSign := buildStringToSign(timestamp, id, canonicalRequest)

	secret := []byte(AuthenticationScheme + credentials.Secret)
	singDate := sing(dateStamp, secret)
	singNonce := sing(nonce, singDate)
	signing := sing(IDTerminator, singNonce)

	signature := sing(string(stringToSign), signing)

	req.Header.Set(AuthorizationHeader, buildAuthorizationHeader(id, signedHeadersString, signature))
}

func buildCanonicalRequest(req *http.Request, payload []byte, signedHeadersString string, sortedHeaderKeys []string) string {
	buffer := buffPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer buffPool.Put(buffer)

	buffer.WriteString(req.Method)
	buffer.WriteString(NL)
	canonicalizeResourcePath(buffer, req.URL.Path)
	buffer.WriteString(NL)
	canonicalizeQueryString(buffer, req.URL.Query())
	buffer.WriteString(NL)
	canonicalizeHeadersString(buffer, req.Header, sortedHeaderKeys)
	buffer.WriteString(NL)
	buffer.WriteString(signedHeadersString)
	buffer.WriteString(NL)
	buffer.WriteString(hex.EncodeToString(sha256Sum(payload)))

	return buffer.String()
}

func buildID(nonce string, dateStamp string, credentials Credentials) string {
	buffer := buffPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer buffPool.Put(buffer)

	buffer.WriteString(credentials.ID)
	buffer.WriteByte(SLASH)
	buffer.WriteString(dateStamp)
	buffer.WriteByte(SLASH)
	buffer.WriteString(nonce)
	buffer.WriteByte(SLASH)
	buffer.WriteString(IDTerminator)

	return buffer.String()
}

func buildStringToSign(timestamp string, id string, canonicalRequest string) []byte {
	buffer := buffPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer buffPool.Put(buffer)

	buffer.WriteString(Algorithm)
	buffer.WriteString(NL)
	buffer.WriteString(timestamp)
	buffer.WriteString(NL)
	buffer.WriteString(id)
	buffer.WriteString(NL)
	buffer.WriteString(hex.EncodeToString(sha256Sum([]byte(canonicalRequest))))

	return buffer.Bytes()
}

func buildAuthorizationHeader(id string, signedHeadersString string, signature []byte) string {
	buffer := buffPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer buffPool.Put(buffer)

	buffer.WriteString(AuthenticationScheme)
	buffer.WriteByte(SPACE)
	//SAUTHC1Id
	buffer.WriteString(SAUTHC1Id)
	buffer.WriteByte(EQ)
	buffer.WriteString(id)

	buffer.WriteString(CS)

	//SAUTHC1SignedHeaders
	buffer.WriteString(SAUTHC1SignedHeaders)
	buffer.WriteByte(EQ)
	buffer.WriteString(signedHeadersString)

	buffer.WriteString(CS)
	//SAUTHC1Signature
	buffer.WriteString(SAUTHC1Signature)
	buffer.WriteByte(EQ)
	buffer.WriteString(hex.EncodeToString(signature))

	return buffer.String()
}

func encodeURL(value string, path bool, canonical bool) string {
	if value == EMPTY {
		return EMPTY
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

func canonicalizeQueryString(buffer *bytes.Buffer, queryValues url.Values) {
	keys := sortedMapKeys(queryValues)
	first := true

	for _, k := range keys {
		key := encodeURL(k, false, true)
		v := queryValues[k]

		for _, vv := range v {
			value := encodeURL(vv, false, true)

			if !first {
				buffer.WriteByte(AMP)
			}

			buffer.WriteString(key)
			buffer.WriteByte(EQ)
			buffer.WriteString(value)
			first = false
		}
	}
}

func canonicalizeResourcePath(buffer *bytes.Buffer, path string) {
	if len(path) == 0 {
		buffer.WriteByte(SLASH)
	} else {
		buffer.WriteString(encodeURL(path, true, true))
	}
}

func canonicalizeHeadersString(buffer *bytes.Buffer, headers http.Header, sortedHeaderKeys []string) {
	for _, k := range sortedHeaderKeys {
		buffer.WriteString(strings.ToLower(k))
		buffer.WriteByte(COLON)

		first := true

		for _, v := range headers[k] {
			if !first {
				buffer.WriteByte(COMMA)
			}
			buffer.WriteString(v)
			first = false
		}
		buffer.WriteString(NL)
	}
}

func signedHeadersString(headers http.Header, sortedHeaderKeys []string) string {
	stringBuffer := bytes.NewBufferString(EMPTY)

	for _, k := range sortedHeaderKeys {
		if stringBuffer.Len() > 0 {
			stringBuffer.WriteByte(SemiColon)
		}
		stringBuffer.WriteString(strings.ToLower(k))
	}

	return stringBuffer.String()
}

func sortedMapKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sing(data string, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func sha256Sum(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}
