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
func NewPageRequest(limit int, offset int) url.Values {
	params := url.Values{}
	params.Add("offset", strconv.Itoa(offset))
	params.Add("limit", strconv.Itoa(limit))
	return params
}

//NewDefaultPageRequest is a conviniece constructor for the default PageRequest values limit = 25 offset = 0
func NewDefaultPageRequest() url.Values {
	return NewPageRequest(25, 0)
}
