package stormpath

import "net/url"

type Criteria interface {
	ToQueryString() string
	Offset(offset int) Criteria
	Limit(limit int) Criteria
}

type baseCriteria struct {
	offset             int
	limit              int
	filter             url.Values
	expandedAttributes []string
}

func (c baseCriteria) ToQueryString() string {
	return requestParams(
		NewPageRequest(c.limit, c.offset),
		c.filter,
		buildExpandParam(c.expandedAttributes),
	)
}

func (c baseCriteria) Offset(offset int) Criteria {
	c.offset = offset
	return c
}

func (c baseCriteria) Limit(limit int) Criteria {
	c.limit = limit
	return c
}
