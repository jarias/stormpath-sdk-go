package stormpath

import "net/url"

const (
	ApplicationBaseUrl = "https://api.stormpath.com/v1/applications"
)

type Application struct {
	Href                       string           `json:"href,omitempty"`
	Name                       string           `json:"name"`
	Description                string           `json:"description,omitempty"`
	Status                     string           `json:"status,omitempty"`
	Accounts                   *Link            `json:"accounts,omitempty"`
	Groups                     *Link            `json:"groups,omitempty"`
	Tenant                     *Link            `json:"tenant,omitempty"`
	PasswordResetTokens        *Link            `json:"passwordResetTokens,omitempty"`
	AccountStoreMappings       *Link            `json:"accountStoreMappings,omitempty"`
	DefaultAccountStoreMapping *Link            `json:"defaultAccountStoreMapping,omitempty"`
	DefaultGroupStoreMapping   *Link            `json:"defaultGroupStoreMapping,omitempty"`
	Client                     *StormpathClient `json:"-"`
}

type Applications struct {
	Href   string         `json:"href"`
	Offset int            `json:"offset"`
	Limit  int            `json:"limit"`
	Items  []*Application `json:"items"`
}

func NewApplication(name string, client *StormpathClient) *Application {
	return &Application{Name: name, Client: client}
}

func (app *Application) Save() error {
	var extraParams = url.Values{}
	extraParams.Add("createDirectory", "true")
	resp, err := app.Client.Do(NewStormpathPostRequest(ApplicationBaseUrl, app, extraParams))

	if err != nil {
		return err
	}

	err = Unmarshal(resp, app)

	return err
}

func (app *Application) Delete() error {
	_, err := app.Client.Do(NewStormpathDeleteRequest(app.Href))

	return err
}

func (app *Application) RegisterAccount(account *Account) error {
	resp, err := app.Client.Do(NewStormpathPostRequest(app.Accounts.Href, account, url.Values{}))

	if err != nil {
		return err
	}

	err = Unmarshal(resp, account)

	return err
}
