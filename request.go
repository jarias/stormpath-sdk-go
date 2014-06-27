package stormpath

import (
	"bytes"
	"encoding/json"
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
	Method              string
	URL                 string
	DontFollowRedirects bool
	Payload             interface{}
	PageRequest         *PageRequest
	Filter              Filter
	ExtraParams         url.Values
}

func NewPageRequest(limit int, offset int) PageRequest {
	return PageRequest{Limit: limit, Offset: offset}
}

func NewDefaultPageRequest() PageRequest {
	return PageRequest{Limit: 25, Offset: 0}
}

func (pageRequest PageRequest) ToUrlQueryValues() url.Values {
	val := url.Values{}

	if pageRequest.Offset >= 0 && pageRequest.Limit > 0 {
		val.Add(OFFSET, strconv.Itoa(pageRequest.Offset))
		val.Add(LIMIT, strconv.Itoa(pageRequest.Limit))
	}

	return val
}

func (request *StormpathRequest) marshalPayload() []byte {
	jsonPayload, _ := json.Marshal(request.Payload)
	return jsonPayload
}

func (request *StormpathRequest) ToHttpRequest() (req *http.Request, err error) {
	var query = url.Values{}

	if request.PageRequest != nil {
		pageRequestQuery := request.PageRequest.ToUrlQueryValues()

		for k, v := range pageRequestQuery {
			query[k] = v
		}
	}

	if request.Filter != nil {
		filterQuery := request.Filter.ToUrlQueryValues()

		for k, v := range filterQuery {
			query[k] = v
		}
	}

	if request.ExtraParams != nil {
		for k, v := range request.ExtraParams {
			query[k] = v
		}
	}

	url := request.URL + "?" + query.Encode()

	req, err = http.NewRequest(request.Method, url, bytes.NewReader(request.marshalPayload()))

	if err != nil {
		return nil, err
	}

	if request.Method == POST || request.Method == PUT {
		req.Header.Add("Content-Type", "application/json")
	}

	return
}
