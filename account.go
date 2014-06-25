package stormpath

type Account struct {
	Href                   *string `json:"href,omitempty"`
	Username               *string `json:"username,omitempty"`
	Email                  string  `json:"email"`
	Password               string  `json:"password"`
	FullName               *string `json:"fullName,omitempty"`
	GivenName              string  `json:"givenName"`
	MiddleName             *string `json:"middleName,omitempty"`
	Surname                string  `json:"surname"`
	Status                 *string `json:"status,omitempty"`
	CustomData             *Link   `json:"customData,omitempty"`
	Groups                 *Link   `json:"groups,omitempty"`
	GroupMemberships       *Link   `json:"groupMemberships,omitempty"`
	Directory              *Link   `json:"directory,omitempty"`
	Tenant                 *Link   `json:"tenant,omitempty"`
	EmailVerificationToken *Link   `json:"emailVerificationToken,omitempty"`
}

func NewAccount(email string, password string, givenName string, surname string) *Account {
	return &Account{Email: email, Password: password, GivenName: givenName, Surname: surname}
}
