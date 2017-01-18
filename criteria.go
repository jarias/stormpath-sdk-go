package stormpath

import "net/url"

const (
	Name        = "name"
	Description = "description"
	Status      = "status"
)

type baseCriteria struct {
	offset             int
	limit              int
	filter             url.Values
	expandedAttributes []string
}

func (c baseCriteria) toQueryString() string {
	return requestParams(
		buildExpandParam(c.expandedAttributes),
		c.filter,
		NewPageRequest(c.limit, c.offset),
	)
}
