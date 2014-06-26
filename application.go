package stormpath

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
)

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
	resp, err := app.Client.Do(&StormpathRequest{
		Method:      POST,
		URL:         ApplicationBaseUrl,
		Payload:     app,
		ExtraParams: extraParams,
	})

	if err != nil {
		return err
	}

	err = unmarshal(resp, app)

	return err
}

func (app *Application) Delete() error {
	_, err := app.Client.Do(&StormpathRequest{
		Method: DELETE,
		URL:    app.Href,
	})

	return err
}

func (app *Application) RegisterAccount(account *Account) error {
	resp, err := app.Client.Do(&StormpathRequest{
		Method:  POST,
		URL:     app.Accounts.Href,
		Payload: account,
	})

	if err != nil {
		return err
	}

	err = unmarshal(resp, account)

	return err
}

func (app *Application) Authorize(username string, password string) (*Account, error) {
	account := &Account{}

	login := make(map[string]string)

	login["type"] = "basic"
	login["value"] = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	jsonLogin, _ := json.Marshal(login)

	resp, err := app.Client.Do(&StormpathRequest{
		Method:  POST,
		URL:     app.Href + "/loginAttempts",
		Payload: jsonLogin,
	})

	if err != nil {
		return account, err
	}

	err = unmarshal(resp, account)

	return account, err
}
