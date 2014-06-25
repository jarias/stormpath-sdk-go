package stormpath

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

func NewApplication(name string) *Application {
	return &Application{Name: name}
}
