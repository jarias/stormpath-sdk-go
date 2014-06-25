package stormpath

import "net/url"

const (
	Name        = "name"
	Description = "description"
	Status      = "status"
)

type DefaultFilter struct {
	Name        string
	Description string
	Status      string
}

func (filter DefaultFilter) ToUrlQueryValues() url.Values {
	values := url.Values{}

	if len(filter.Name) > 0 {
		values.Set(Name, filter.Name)
	}
	if len(filter.Description) > 0 {
		values.Set(Description, filter.Description)
	}
	if len(filter.Status) > 0 {
		values.Set(Status, filter.Status)
	}

	return values
}
