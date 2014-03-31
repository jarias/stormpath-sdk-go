package stormpath

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	OFFSET = "offset"
	LIMIT  = "limit"
)

type PageRequest struct {
	Limit  int
	Offset int
}

type StormpathRequest struct {
	Method          string
	URL             string
	FollowRedirects bool
	Query           url.Values
	Payload         []byte
}

func NewDefaultPageRequest() *PageRequest {
	return &PageRequest{Limit: 25, Offset: 0}
}

func NewStormpathRequest(method string, url string, query url.Values) *StormpathRequest {
	return &StormpathRequest{Method: method, URL: url, Query: query, Payload: []byte(""), FollowRedirects: true}
}

func NewStormpathRequestNoRedirects(method string, url string, query url.Values) *StormpathRequest {
	return &StormpathRequest{Method: method, URL: url, Query: query, Payload: []byte(""), FollowRedirects: false}
}

func (pageRequest *PageRequest) ToUrlQueryValues() url.Values {
	val := url.Values{}

	val.Add(OFFSET, strconv.Itoa(pageRequest.Offset))
	val.Add(LIMIT, strconv.Itoa(pageRequest.Limit))

	return val
}

func (request *StormpathRequest) ToHttpRequest() (req *http.Request, err error) {
	url := request.URL + "?" + request.Query.Encode()
	req, err = http.NewRequest(request.Method, url, bytes.NewReader(request.Payload))

	if err != nil {
		return nil, err
	}

	return
}
