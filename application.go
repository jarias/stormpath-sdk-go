package stormpath

import (
	"encoding/base64"
	"net/url"
	"time"

	"github.com/nu7hatch/gouuid"

	"github.com/dgrijalva/jwt-go"
)

//Application represents a Stormpath application object
//
//See: http://docs.stormpath.com/rest/product-guide/#applications
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

//Applications represents a paged result or applications
type Applications struct {
	list
	Items []Application `json:"items"`
}

//NewApplication creates a new application
func NewApplication(name string) *Application {
	return &Application{Name: name}
}

//Save saves the given application
//
//See: http://docs.stormpath.com/rest/product-guide/#create-an-application-aka-register-an-application-with-stormpath
func (app *Application) Save() error {
	return client.post(app.Href, app, app)
}

//Delete deletes the given applicaiton
//
//See: http://docs.stormpath.com/rest/product-guide/#delete-an-application
func (app *Application) Delete() error {
	return client.delete(app.Href, emptyPayload())
}

//Purge deletes all the account stores before deleting the application
//
//See: http://docs.stormpath.com/rest/product-guide/#delete-an-application
func (app *Application) Purge() error {
	accountStoreMappings, err := app.GetAccountStoreMappings(NewDefaultPageRequest(), NewEmptyFilter())
	if err != nil {
		return err
	}
	for _, m := range accountStoreMappings.Items {
		client.delete(m.AccountStore.Href, emptyPayload())
	}

	return app.Delete()
}

//GetAccountStoreMappings returns all the applications account store mappings
//
//See: http://docs.stormpath.com/rest/product-guide/#application-account-store-mappings
func (app *Application) GetAccountStoreMappings(pageRequest url.Values, filter url.Values) (*AccountStoreMappings, error) {
	accountStoreMappings := &AccountStoreMappings{}

	err := client.get(
		buildAbsoluteURL(app.AccountStoreMappings.Href, requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
		accountStoreMappings,
	)

	return accountStoreMappings, err
}

//GetAccounts returns all the accounts of the application
//
//See: http://docs.stormpath.com/rest/product-guide/#application-accounts
func (app *Application) GetAccounts(pageRequest url.Values, filter url.Values) (*Accounts, error) {
	accounts := &Accounts{}

	err := client.get(
		buildAbsoluteURL(app.Accounts.Href, requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
		accounts,
	)

	return accounts, err
}

//RegisterAccount registers a new account into the application
//
//See: http://docs.stormpath.com/rest/product-guide/#create-an-application-aka-register-an-application-with-stormpath
func (app *Application) RegisterAccount(account *Account) error {
	return client.post(app.Accounts.Href, account, account)
}

//AuthenticateAccount authenticates an account against the application
//
//See: http://docs.stormpath.com/rest/product-guide/#authenticate-an-account
func (app *Application) AuthenticateAccount(username string, password string) (*AccountRef, error) {
	account := &AccountRef{}

	loginAttemptPayload := make(map[string]string)

	loginAttemptPayload["type"] = "basic"
	loginAttemptPayload["value"] = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	err := client.post(buildAbsoluteURL(app.Href, "loginAttempts"), loginAttemptPayload, account)

	return account, err
}

//SendPasswordResetEmail sends a password reset email to the given user
//
//See: http://docs.stormpath.com/rest/product-guide/#reset-an-accounts-password
func (app *Application) SendPasswordResetEmail(username string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	passwordResetPayload := make(map[string]string)
	passwordResetPayload["email"] = username

	err := client.post(buildAbsoluteURL(app.Href, "passwordResetTokens"), passwordResetPayload, passwordResetToken)

	return passwordResetToken, err
}

//ValidatePasswordResetToken validates a password reset token
//
//See: http://docs.stormpath.com/rest/product-guide/#reset-an-accounts-password
func (app *Application) ValidatePasswordResetToken(token string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	err := client.get(buildAbsoluteURL(app.Href, "passwordResetTokens", token), emptyPayload(), passwordResetToken)

	return passwordResetToken, err
}

//ResetPassword resets a user password based on the reset token
//
//See: http://docs.stormpath.com/rest/product-guide/#reset-an-accounts-password
func (app *Application) ResetPassword(token string, newPassword string) (*AccountRef, error) {
	account := &AccountRef{}

	resetPasswordPayload := make(map[string]string)
	resetPasswordPayload["password"] = newPassword

	err := client.post(buildAbsoluteURL(app.Href, "passwordResetTokens", token), resetPasswordPayload, account)

	return account, err
}

//CreateGroup creates a new group in the application
//
//See: http://docs.stormpath.com/rest/product-guide/#application-groups
func (app *Application) CreateGroup(group *Group) error {
	return client.post(app.Groups.Href, group, group)
}

//GetGroups returns all the application groups
//
//See: http://docs.stormpath.com/rest/product-guide/#application-groups
func (app *Application) GetGroups(pageRequest url.Values, filter url.Values) (*Groups, error) {
	groups := &Groups{}

	err := client.get(
		buildAbsoluteURL(app.Groups.Href, requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
		groups,
	)

	return groups, err
}

//CreateIDSiteURL creates the IDSite URL for the application
func (app *Application) CreateIDSiteURL(options map[string]string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	nonce, _ := uuid.NewV4()

	if options["path"] == "" {
		options["path"] = "/"
	}

	token.Claims["jti"] = nonce.String()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["iss"] = client.Credentials.ID
	token.Claims["sub"] = app.Href
	token.Claims["state"] = options["state"]
	token.Claims["path"] = options["path"]
	token.Claims["cb_uri"] = options["callbackURI"]

	tokenString, err := token.SignedString([]byte(client.Credentials.Secret))
	if err != nil {
		return "", err
	}

	p, _ := url.Parse(app.Href)
	ssoURL := p.Scheme + "://" + p.Host + "/sso" + "?jwtRequest=" + tokenString

	return ssoURL, nil
}
