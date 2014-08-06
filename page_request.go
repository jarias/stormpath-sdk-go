package stormpath

import (
	"net/url"
	"strconv"
)

//PageRequest contains the limit and offset values for any paginated Stormpath request
type PageRequest struct {
	Limit  int
	Offset int
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
