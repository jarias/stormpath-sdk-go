package stormpath

import (
	"encoding/base64"

	"net/url"
)

const (
	ApplicationBaseUrl = "https://api.stormpath.com/v1/applications"
)

type Application struct {
	Href                       string `json:"href,omitempty"`
	Name                       string `json:"name"`
	Description                string `json:"description,omitempty"`
	Status                     string `json:"status,omitempty"`
	Accounts                   *Link  `json:"accounts,omitempty"`
	Groups                     *Link  `json:"groups,omitempty"`
	Tenant                     *Link  `json:"tenant,omitempty"`
	PasswordResetTokens        *Link  `json:"passwordResetTokens,omitempty"`
	AccountStoreMappings       *Link  `json:"accountStoreMappings,omitempty"`
	DefaultAccountStoreMapping *Link  `json:"defaultAccountStoreMapping,omitempty"`
	DefaultGroupStoreMapping   *Link  `json:"defaultGroupStoreMapping,omitempty"`
}

func NewApplication(name string) *Application {
	return &Application{Name: name}
}

func (app *Application) Save() error {
	var extraParams = url.Values{}
	extraParams.Add("createDirectory", "true")

	return Client.DoWithResult(&StormpathRequest{
		Method:      POST,
		URL:         ApplicationBaseUrl,
		Payload:     app,
		ExtraParams: extraParams,
	}, app)
}

func (app *Application) Delete() error {
	_, err := Client.Do(&StormpathRequest{
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
		Client.Do(&StormpathRequest{
			Method: DELETE,
			URL:    m.AccountStore.Href,
		})
	}

	return app.Delete()
}

func (app *Application) GetAccountStoreMappings(pageRequest PageRequest, filter DefaultFilter) (*AccountStoreMappings, error) {
	accountStoreMappings := &AccountStoreMappings{}

	err := Client.DoWithResult(&StormpathRequest{
		Method:      GET,
		URL:         app.AccountStoreMappings.Href,
		PageRequest: &pageRequest,
		Filter:      &filter,
	}, accountStoreMappings)

	return accountStoreMappings, err
}

func (app *Application) GetAccounts(pageRequest PageRequest, filter AccountFilter) (*Accounts, error) {
	accounts := &Accounts{}

	err := Client.DoWithResult(&StormpathRequest{
		Method:      GET,
		URL:         app.Accounts.Href,
		PageRequest: &pageRequest,
		Filter:      &filter,
	}, accounts)

	return accounts, err
}

func (app *Application) RegisterAccount(account *Account) error {
	err := Client.DoWithResult(&StormpathRequest{
		Method:  POST,
		URL:     app.Accounts.Href,
		Payload: account,
	}, account)

	return err
}

func (app *Application) AuthenticateAccount(username string, password string) (*AccountRef, error) {
	account := &AccountRef{}

	loginAttemptPayload := make(map[string]string)

	loginAttemptPayload["type"] = "basic"
	loginAttemptPayload["value"] = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	err := Client.DoWithResult(&StormpathRequest{
		Method:  POST,
		URL:     app.Href + "/loginAttempts",
		Payload: loginAttemptPayload,
	}, account)

	return account, err
}

func (app *Application) SendPasswordResetEmail(username string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	passwordResetPayload := make(map[string]string)
	passwordResetPayload["email"] = username

	err := Client.DoWithResult(&StormpathRequest{
		Method:  POST,
		URL:     app.Href + "/passwordResetTokens",
		Payload: passwordResetPayload,
	}, passwordResetToken)

	return passwordResetToken, err
}

func (app *Application) ValidatePasswordResetToken(token string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	err := Client.DoWithResult(&StormpathRequest{
		Method: GET,
		URL:    app.Href + "/passwordResetTokens/" + token,
	}, passwordResetToken)

	return passwordResetToken, err
}

func (app *Application) ResetPassword(token string, newPassword string) (*AccountRef, error) {
	account := &AccountRef{}

	resetPasswordPayload := make(map[string]string)
	resetPasswordPayload["password"] = newPassword

	err := Client.DoWithResult(&StormpathRequest{
		Method:  POST,
		URL:     app.Href + "/passwordResetTokens/" + token,
		Payload: resetPasswordPayload,
	}, account)

	return account, err
}

func (app *Application) CreateApplicationGroup(group *Group) error {
	return Client.DoWithResult(&StormpathRequest{
		Method:  POST,
		URL:     app.Groups.Href,
		Payload: group,
	}, group)
}

func (app *Application) GetApplicationGroups(pageRequest PageRequest, filters Filter) (*Groups, error) {
	groups := &Groups{}

	err := Client.DoWithResult(&StormpathRequest{
		Method: GET,
		URL:    app.Groups.Href,
	}, groups)

	return groups, err
}
