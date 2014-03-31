package stormpath

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
