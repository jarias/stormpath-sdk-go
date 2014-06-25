package stormpath

import (
	"net/url"
)

const (
	NAME        = "name"
	DESCRIPTION = "description"
	STATUS      = "status"
)

type Application struct {
	Href        string
	Name        string
	Description string
	Status      string
	Accounts    struct {
		Href string
	}
	Groups struct {
		Href string
	}
	Tenant struct {
		Href string
	}
	PasswordResetTokens struct {
		Href string
	}
}

type Applications struct {
	Href   string
	Offset int
	Limit  int
	Items  []Application
}

type ApplicationFilter struct {
	Name        string
	Description string
	Status      string
}

func (filter ApplicationFilter) ToUrlQueryValues() url.Values {
	values := url.Values{}

	if len(filter.Name) > 0 {
		values.Set(NAME, filter.Name)
	}
	if len(filter.Description) > 0 {
		values.Set(DESCRIPTION, filter.Description)
	}
	if len(filter.Status) > 0 {
		values.Set(STATUS, filter.Status)
	}

	return values
}
