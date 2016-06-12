package stormpath

import "net/url"

const (
	Name        = "name"
	Description = "description"
	Status      = "status"
)

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
		buildExpandParam(c.expandedAttributes),
		c.filter,
		NewPageRequest(c.limit, c.offset),
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
