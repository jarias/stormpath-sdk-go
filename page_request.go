package stormpath

import (
	"fmt"
	"net/url"
	"strconv"
)

var DefaultPageRequest = PageRequest{25, 0}

//PageRequest contains the limit and offset values for any paginated Stormpath request
type PageRequest struct {
	Limit  int
	Offset int
}

func (r PageRequest) toExpansion(attribute string) string {
	return fmt.Sprintf("%s(offset:%d,limit:%d)", attribute, r.Offset, r.Limit)
}

//NewPageRequest is a conviniece constructor for a PageRequest
func NewPageRequest(limit int, offset int) url.Values {
	params := url.Values{}
	//limit == 0 invalid so we sanitize it for the user
	if limit != 0 {
		params.Add("offset", strconv.Itoa(offset))
		params.Add("limit", strconv.Itoa(limit))
	}
	return params
}
