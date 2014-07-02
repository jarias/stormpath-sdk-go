package stormpath

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

//PageRequest contains the limit and offset values for any paginated Stormpath request
type PageRequest struct {
	Limit  int
	Offset int
}

//Request method constants
const (
	Get    = "GET"
	Post   = "POST"
	Delete = "DELETE"
	Put    = "PUT"
)

//StormpathRequest is an abstraction of an the actual HTTP request for a Stormpath operation
//
//Fields:
//	Method 				(GET,POST,PUT,DELETE,etc.)
//	URL 				the request URL
//	DontFollowRedirects this bool flag is only use for getting the current tenant, the value is always false unless set
//	Payload 			the request payload for PUT and POST requests
//	PageRequest 		if set defines the request pagination
//	Filter 				if set defines the query filter for list based requests
//	ExtraParams 		any aditional optional params needed for a specific request
type StormpathRequest struct {
	Method              string
	URL                 string
	DontFollowRedirects bool
	Payload             interface{}
	PageRequest         PageRequest
	Filter              Filter
	ExtraParams         url.Values
}

//NewPageRequest is a conviniece constructor for a PageRequest
func NewPageRequest(limit int, offset int) PageRequest {
	return PageRequest{Limit: limit, Offset: offset}
}

//NewDefaultPageRequest is a conviniece constructor for the default PageRequest values limit = 25 offset = 0
func NewDefaultPageRequest() PageRequest {
	return PageRequest{Limit: 25, Offset: 0}
}

func (pageRequest PageRequest) toURLQueryValues() url.Values {
	val := url.Values{}

	if pageRequest.Offset >= 0 && pageRequest.Limit > 0 {
		val.Add("offset", strconv.Itoa(pageRequest.Offset))
		val.Add("limit", strconv.Itoa(pageRequest.Limit))
	}

	return val
}

func (request *StormpathRequest) marshalPayload() []byte {
	jsonPayload, _ := json.Marshal(request.Payload)
	return jsonPayload
}

//ToHTTPRequest coverts a StormathRequest to an http.Request, returns an error if any
//
//Headers:
//	User-Agent:   jarias/stormpath-sdk-go/ + Current SDK version
//	Accept: 	  application/json
//	Content-Type: application/json
func (request *StormpathRequest) ToHTTPRequest() (*http.Request, error) {
	var query = url.Values{}

	pageRequestQuery := request.PageRequest.toURLQueryValues()

	for k, v := range pageRequestQuery {
		query[k] = v
	}

	if request.Filter != nil {
		filterQuery := request.Filter.toURLQueryValues()

		for k, v := range filterQuery {
			query[k] = v
		}
	}

	if request.ExtraParams != nil {
		for k, v := range request.ExtraParams {
			query[k] = v
		}
	}

	req, err := http.NewRequest(request.Method, request.URL+"?"+query.Encode(), bytes.NewReader(request.marshalPayload()))

	req.Header.Set("User-Agent", "jarias/stormpath-sdk-go/"+Version)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, err
}
