package stormpath

import (
	"encoding/base64"
	"net/url"
)

type Application struct {
	Href                       string `json:"href,omitempty"`
	Name                       string `json:"name"`
	Description                string `json:"description,omitempty"`
	Status                     string `json:"status,omitempty"`
	Accounts                   *link  `json:"accounts,omitempty"`
	Groups                     *link  `json:"groups,omitempty"`
	Tenant                     *link  `json:"tenant,omitempty"`
	PasswordResetTokens        *link  `json:"passwordResetTokens,omitempty"`
	AccountStoreMappings       *link  `json:"accountStoreMappings,omitempty"`
	DefaultAccountStoreMapping *link  `json:"defaultAccountStoreMapping,omitempty"`
	DefaultGroupStoreMapping   *link  `json:"defaultGroupStoreMapping,omitempty"`
}

type Applications struct {
	list
	Items []Application `json:"items"`
}

func NewApplication(name string) *Application {
	return &Application{Name: name}
}

func (app *Application) Save() error {
	return Client.doWithResult(Client.newRequest(
		"POST",
		app.Href,
		app,
	), app)
}

func (app *Application) Delete() error {
	return Client.do(Client.newRequest(
		"DELETE",
		app.Href,
		emptyPayload(),
	))
}

func (app *Application) Purge() error {
	accountStoreMappings, err := app.GetAccountStoreMappings(NewDefaultPageRequest(), NewEmptyFilter())
	if err != nil {
		return err
	}
	for _, m := range accountStoreMappings.Items {
		Client.do(Client.newRequest(
			"DELETE",
			m.AccountStore.Href,
			emptyPayload(),
		))
	}

	return app.Delete()
}

func (app *Application) GetAccountStoreMappings(pageRequest url.Values, filter url.Values) (*AccountStoreMappings, error) {
	accountStoreMappings := &AccountStoreMappings{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		buildAbsoluteURL(app.AccountStoreMappings.Href, requestParams(pageRequest, filter, url.Values{})),
		nil,
	), accountStoreMappings)

	return accountStoreMappings, err
}

func (app *Application) GetAccounts(pageRequest url.Values, filter url.Values) (*Accounts, error) {
	accounts := &Accounts{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		buildAbsoluteURL(app.Accounts.Href, requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
	), accounts)

	return accounts, err
}

func (app *Application) RegisterAccount(account *Account) error {
	err := Client.doWithResult(Client.newRequest(
		"POST",
		app.Accounts.Href,
		account,
	), account)

	return err
}

func (app *Application) AuthenticateAccount(username string, password string) (*AccountRef, error) {
	account := &AccountRef{}

	loginAttemptPayload := make(map[string]string)

	loginAttemptPayload["type"] = "basic"
	loginAttemptPayload["value"] = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	err := Client.doWithResult(Client.newRequest(
		"POST",
		buildAbsoluteURL(app.Href, "loginAttempts"),
		loginAttemptPayload,
	), account)

	return account, err
}

func (app *Application) SendPasswordResetEmail(username string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	passwordResetPayload := make(map[string]string)
	passwordResetPayload["email"] = username

	err := Client.doWithResult(Client.newRequest(
		"POST",
		buildAbsoluteURL(app.Href, "passwordResetTokens"),
		passwordResetPayload,
	), passwordResetToken)

	return passwordResetToken, err
}

func (app *Application) ValidatePasswordResetToken(token string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		buildAbsoluteURL(app.Href, "passwordResetTokens", token),
		emptyPayload(),
	), passwordResetToken)

	return passwordResetToken, err
}

func (app *Application) ResetPassword(token string, newPassword string) (*AccountRef, error) {
	account := &AccountRef{}

	resetPasswordPayload := make(map[string]string)
	resetPasswordPayload["password"] = newPassword

	err := Client.doWithResult(Client.newRequest(
		"POST",
		buildAbsoluteURL(app.Href, "passwordResetTokens", token),
		resetPasswordPayload,
	), account)

	return account, err
}

func (app *Application) CreateApplicationGroup(group *Group) error {
	return Client.doWithResult(Client.newRequest(
		"POST",
		app.Groups.Href,
		group,
	), group)
}

func (app *Application) GetApplicationGroups(pageRequest url.Values, filter url.Values) (*Groups, error) {
	groups := &Groups{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		buildAbsoluteURL(app.Groups.Href, requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
	), groups)

	return groups, err
}
