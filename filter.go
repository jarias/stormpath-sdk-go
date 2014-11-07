package stormpath

import "net/url"

func NewDefaultFilter(name string, description string, status string) url.Values {
	filter := url.Values{}

	if len(name) > 0 {
		filter.Add("name", name)
	}
	if len(description) > 0 {
		filter.Add("description", description)
	}
	if len(description) > 0 {
		filter.Add("status", status)
	}

	return filter
}

func NewAccountFilter(givenName string, middleName string, surname string, username string, email string) url.Values {
	filter := url.Values{}

	if len(givenName) > 0 {
		filter.Add("givenName", givenName)
	}
	if len(middleName) > 0 {
		filter.Add("middleName", middleName)
	}
	if len(surname) > 0 {
		filter.Add("surname", surname)
	}
	if len(username) > 0 {
		filter.Add("username", username)
	}
	if len(email) > 0 {
		filter.Add("email", email)
	}

	return filter
}

func NewEmptyFilter() url.Values {
	return url.Values{}
}
