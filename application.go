package stormpath

import (
	"net/url"
)

const (
	Name        = "name"
	Description = "description"
	Status      = "status"
)

type Application struct {
	Href                *string          `json:"href,omitempty"`
	Name                string           `json:"name"`
	Description         *string          `json:"description,omitempty"`
	Status              *string          `json:"status,omitempty"`
	Accounts            *Link            `json:"accounts,omitempty"`
	Groups              *Link            `json:"groups,omitempty"`
	Tenant              *Link            `json:"tenant,omitempty"`
	PasswordResetTokens *Link            `json:"passwordResetTokens,omitempty"`
	Client              *StormpathClient `json:"-"`
}

type Applications struct {
	Href   string        `json:"href"`
	Offset int           `json:"offset"`
	Limit  int           `json:"limit"`
	Items  []Application `json:"items"`
}

type ApplicationFilter struct {
	Name        string
	Description string
	Status      string
}

func NewApplication(name string) *Application {
	return &Application{Name: name}
}

func (filter ApplicationFilter) ToUrlQueryValues() url.Values {
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
