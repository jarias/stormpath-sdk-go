package stormpath

import (
	"encoding/base64"

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

func (app *Application) Purge() error {
	accountStoreMappings, err := app.GetAccountStoreMappings(NewDefaultPageRequest(), DefaultFilter{})
	if err != nil {
		return err
	}
	for _, m := range accountStoreMappings.Items {
		app.Client.Do(&StormpathRequest{
			Method: DELETE,
			URL:    m.AccountStore.Href,
		})
	}

	return app.Delete()
}

func (app *Application) GetAccountStoreMappings(pageRequest PageRequest, filter DefaultFilter) (*AccountStoreMappings, error) {
	accountStoreMappings := &AccountStoreMappings{}

	resp, err := app.Client.Do(&StormpathRequest{
		Method:      GET,
		URL:         app.AccountStoreMappings.Href,
		PageRequest: &pageRequest,
		Filter:      &filter,
	})

	if err != nil {
		return accountStoreMappings, err
	}

	err = unmarshal(resp, accountStoreMappings)
	for _, m := range accountStoreMappings.Items {
		m.Client = app.Client
	}

	return accountStoreMappings, err
}

func (app *Application) GetAccounts(pageRequest PageRequest, filter AccountFilter) (*Accounts, error) {
	accounts := &Accounts{}

	resp, err := app.Client.Do(&StormpathRequest{
		Method:      GET,
		URL:         app.Accounts.Href,
		PageRequest: &pageRequest,
		Filter:      &filter,
	})

	if err != nil {
		return accounts, err
	}

	err = unmarshal(resp, accounts)
	for _, a := range accounts.Items {
		a.Client = app.Client
	}

	return accounts, err
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

func (app *Application) AuthenticateAccount(username string, password string) (*AccountRef, error) {
	account := &AccountRef{}

	loginAttemptPayload := make(map[string]string)

	loginAttemptPayload["type"] = "basic"
	loginAttemptPayload["value"] = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	resp, err := app.Client.Do(&StormpathRequest{
		Method:  POST,
		URL:     app.Href + "/loginAttempts",
		Payload: loginAttemptPayload,
	})

	if err != nil {
		return account, err
	}

	err = unmarshal(resp, account)

	return account, err
}

func (app *Application) SendPasswordResetEmail(username string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	passwordResetPayload := make(map[string]string)
	passwordResetPayload["email"] = username

	resp, err := app.Client.Do(&StormpathRequest{
		Method:  POST,
		URL:     app.Href + "/passwordResetTokens",
		Payload: passwordResetPayload,
	})

	if err != nil {
		return passwordResetToken, err
	}

	err = unmarshal(resp, passwordResetToken)

	return passwordResetToken, err
}

func (app *Application) ValidatePasswordResetToken(token string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	resp, err := app.Client.Do(&StormpathRequest{
		Method: GET,
		URL:    app.Href + "/passwordResetTokens/" + token,
	})

	if err != nil {
		return passwordResetToken, err
	}

	err = unmarshal(resp, passwordResetToken)

	return passwordResetToken, err
}

func (app *Application) ResetPassword(token string, newPassword string) (*AccountRef, error) {
	account := &AccountRef{}

	resetPasswordPayload := make(map[string]string)
	resetPasswordPayload["password"] = newPassword

	resp, err := app.Client.Do(&StormpathRequest{
		Method:  POST,
		URL:     app.Href + "/passwordResetTokens/" + token,
		Payload: resetPasswordPayload,
	})

	if err != nil {
		return account, err
	}

	err = unmarshal(resp, account)

	return account, err
}
