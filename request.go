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

type Filter interface {
	ToUrlQueryValues() url.Values
}

type PageRequest struct {
	Limit  int
	Offset int
}

type StormpathRequest struct {
	Method          string
	URL             string
	FollowRedirects bool
	Payload         []byte
	PageRequest     PageRequest
	Filter          Filter
}

func NewPageRequest(limit int, offset int) PageRequest {
	return PageRequest{Limit: limit, Offset: offset}
}

func NewDefaultPageRequest() PageRequest {
	return PageRequest{Limit: 25, Offset: 0}
}

func NewStormpathRequest(method string, url string, pageRequest PageRequest, filter Filter) *StormpathRequest {
	return &StormpathRequest{Method: method, URL: url, PageRequest: pageRequest, Filter: filter, Payload: []byte(""), FollowRedirects: true}
}

func NewStormpathRequestNoRedirects(method string, url string, pageRequest PageRequest, filter Filter) *StormpathRequest {
	return &StormpathRequest{Method: method, URL: url, PageRequest: pageRequest, Filter: filter, Payload: []byte(""), FollowRedirects: false}
}

func (pageRequest PageRequest) ToUrlQueryValues() url.Values {
	val := url.Values{}

	if pageRequest.Offset >= 0 && pageRequest.Limit > 0 {
		val.Add(OFFSET, strconv.Itoa(pageRequest.Offset))
		val.Add(LIMIT, strconv.Itoa(pageRequest.Limit))
	}

	return val
}

func (request *StormpathRequest) ToHttpRequest() (req *http.Request, err error) {
	query := request.PageRequest.ToUrlQueryValues()

	filterQuery := request.Filter.ToUrlQueryValues()

	for k, v := range filterQuery {
		query[k] = v
	}

	url := request.URL + "?" + query.Encode()
	req, err = http.NewRequest(request.Method, url, bytes.NewReader(request.Payload))

	if err != nil {
		return nil, err
	}

	return
}
