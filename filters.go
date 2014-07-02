package stormpath

import "net/url"

const (
	Name        = "name"
	Description = "description"
	Status      = "status"

	GivenName  = "GivenName"
	MiddleName = "MiddleName"
	Surname    = "Surname"
	Username   = "Username"
	Email      = "Email"
)

type Filter interface {
	toURLQueryValues() url.Values
}

type DefaultFilter struct {
	Name        string
	Description string
	Status      string
}

type AccountFilter struct {
	GivenName  string
	MiddleName string
	Surname    string
	Username   string
	Email      string
}

func (filter AccountFilter) toURLQueryValues() url.Values {
	values := url.Values{}

	if len(filter.GivenName) > 0 {
		values.Set(GivenName, filter.GivenName)
	}
	if len(filter.MiddleName) > 0 {
		values.Set(MiddleName, filter.MiddleName)
	}
	if len(filter.Surname) > 0 {
		values.Set(Surname, filter.Surname)
	}
	if len(filter.Username) > 0 {
		values.Set(Username, filter.Username)
	}
	if len(filter.Email) > 0 {
		values.Set(Email, filter.Email)
	}

	return values
}

func (filter DefaultFilter) toURLQueryValues() url.Values {
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
