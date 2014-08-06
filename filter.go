package stormpath

import "net/url"

//Common filter fields used in paged resutls
const (
	Name        = "name"
	Description = "description"
	Status      = "status"
)

//Account specific filter fields used in paged resutls
const (
	GivenName  = "GivenName"
	MiddleName = "MiddleName"
	Surname    = "Surname"
	Username   = "Username"
	Email      = "Email"
)

//Filter defines the interface for any filter type that can be converted to url.Values
type Filter interface {
	toURLQueryValues() url.Values
}

//DefaultFilter is the common filter among Stormpath model objects
type DefaultFilter struct {
	Name        string
	Description string
	Status      string
}

//AccountFilter is the specific filter use for Accounts
type AccountFilter struct {
	GivenName  string
	MiddleName string
	Surname    string
	Username   string
	Email      string
}

func (filter AccountFilter) toURLQueryValues() url.Values {
	values := url.Values{}

	copyFieldFilter(filter.GivenName, GivenName, values)
	copyFieldFilter(filter.MiddleName, MiddleName, values)
	copyFieldFilter(filter.Surname, Surname, values)
	copyFieldFilter(filter.Username, Username, values)
	copyFieldFilter(filter.Email, Email, values)

	return values
}

func (filter DefaultFilter) toURLQueryValues() url.Values {
	values := url.Values{}

	copyFieldFilter(filter.Name, Name, values)
	copyFieldFilter(filter.Description, Description, values)
	copyFieldFilter(filter.Status, Status, values)

	return values
}

func copyFieldFilter(fieldValue string, fieldName string, to url.Values) {
	if len(fieldValue) > 0 {
		to.Set(fieldName, fieldValue)
	}
}
