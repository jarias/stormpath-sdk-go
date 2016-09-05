package stormpath

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/nu7hatch/gouuid"
)

//Application represents a Stormpath application object
//
//See: http://docs.stormpath.com/rest/product-guide/#applications
type Application struct {
	accountStoreResource
	Name                       string                           `json:"name,omitempty"`
	Description                string                           `json:"description,omitempty"`
	Status                     string                           `json:"status,omitempty"`
	Groups                     *Groups                          `json:"groups,omitempty"`
	Tenant                     *Tenant                          `json:"tenant,omitempty"`
	PasswordResetTokens        *resource                        `json:"passwordResetTokens,omitempty"`
	AccountStoreMappings       *ApplicationAccountStoreMappings `json:"accountStoreMappings,omitempty"`
	DefaultAccountStoreMapping *ApplicationAccountStoreMapping  `json:"defaultAccountStoreMapping,omitempty"`
	DefaultGroupStoreMapping   *ApplicationAccountStoreMapping  `json:"defaultGroupStoreMapping,omitempty"`
	OAuthPolicy                *OAuthPolicy                     `json:"oAuthPolicy,omitempty"`
	APIKeys                    *APIKeys                         `json:"apiKeys,omitempty"`
}

//Applications represents a paged result or applications
type Applications struct {
	collectionResource
	Items []Application `json:"items,omitempty"`
}

//CallbackResult holds the ID Site callback parsed JWT token information + the acccount if one was given
type CallbackResult struct {
	Account *Account
	State   string
	IsNew   bool
	Status  string
}

//IDSiteOptions represents the posible options to generate an new IDSite URL.
type IDSiteOptions struct {
	Logout      bool
	Path        string
	CallbackURL string
	State       string
}

//NewApplication creates a new application
func NewApplication(name string) *Application {
	return &Application{Name: name}
}

//GetApplication loads an application by href and criteria
func GetApplication(href string, criteria Criteria) (*Application, error) {
	application := &Application{}

	err := client.get(
		buildAbsoluteURL(href, criteria.ToQueryString()),
		application,
	)

	return application, err
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (app *Application) Refresh() error {
	return client.get(app.Href, app)
}

//Update updates the given resource, by doing a POST to the resource Href
func (app *Application) Update() error {
	return client.post(app.Href, app, app)
}

//Purge deletes all the account stores before deleting the application
//
//See: http://docs.stormpath.com/rest/product-guide/#delete-an-application
func (app *Application) Purge() error {
	accountStoreMappings, err := app.GetAccountStoreMappings(MakeApplicationAccountStoreMappingCriteria().Offset(0).Limit(25))
	if err != nil {
		return err
	}
	for _, m := range accountStoreMappings.Items {
		client.delete(m.AccountStore.Href)
	}

	return app.Delete()
}

//GetAccountStoreMappings returns all the applications account store mappings
//
//See: http://docs.stormpath.com/rest/product-guide/#application-account-store-mappings
func (app *Application) GetAccountStoreMappings(criteria Criteria) (*ApplicationAccountStoreMappings, error) {
	accountStoreMappings := &ApplicationAccountStoreMappings{}

	err := client.get(
		buildAbsoluteURL(app.AccountStoreMappings.Href, criteria.ToQueryString()),
		accountStoreMappings,
	)

	if err != nil {
		return nil, err
	}

	return accountStoreMappings, nil
}

//GetDefaultAccountStoreMapping retrieves the default ApplicationAccountStoreMapping for the given Application
func (app *Application) GetDefaultAccountStoreMapping(criteria Criteria) (*ApplicationAccountStoreMapping, error) {
	err := client.get(
		buildAbsoluteURL(app.DefaultAccountStoreMapping.Href, criteria.ToQueryString()),
		app.DefaultAccountStoreMapping,
	)

	if err != nil {
		return nil, err
	}

	return app.DefaultAccountStoreMapping, nil
}

//RegisterAccount registers a new account into the application
//
//See: http://docs.stormpath.com/rest/product-guide/#application-accounts
func (app *Application) RegisterAccount(account *Account) error {
	err := client.post(app.Accounts.Href, account, account)
	if err == nil {
		//Password should be cleanup so we don't keep an unhash password in memory
		account.Password = ""
	}
	return err
}

//RegisterSocialAccount registers a new account into the application using an external provider Google, Facebook
//
//See: http://docs.stormpath.com/rest/product-guide/#accessing-accounts-with-google-authorization-codes-or-an-access-tokens
func (app *Application) RegisterSocialAccount(socialAccount *SocialAccount) (*Account, error) {
	account := &Account{}

	err := client.post(app.Accounts.Href, socialAccount, account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

//AuthenticateAccount authenticates an account against the application
//
//See: http://docs.stormpath.com/rest/product-guide/#authenticate-an-account
func (app *Application) AuthenticateAccount(username string, password string, accountStoreHref string) (*Account, error) {
	accountRef := &accountRef{Account: &Account{}}

	loginAttemptPayload := make(map[string]interface{})
	loginAttemptPayload["type"] = "basic"
	loginAttemptPayload["value"] = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	if accountStoreHref != "" {
		loginAttemptPayload["accountStore"] = map[string]string{
			"href": accountStoreHref,
		}
	}

	err := client.post(buildAbsoluteURL(app.Href, "loginAttempts"), loginAttemptPayload, accountRef)

	if err != nil {
		return nil, err
	}

	account := accountRef.Account
	err = account.Refresh()
	if err != nil {
		return nil, err
	}

	return account, nil
}

//ResendVerificationEmail resends the verification email to the given email address
//
//See: https://docs.stormpath.com/rest/product-guide/latest/accnt_mgmt.html#how-to-verify-an-account-s-email
func (app *Application) ResendVerificationEmail(email string) error {
	resendVerificationEmailPayload := map[string]string{
		"login": email,
	}
	return client.post(buildAbsoluteURL(app.Href, "verificationEmails"), resendVerificationEmailPayload, nil)
}

//SendPasswordResetEmail sends a password reset email to the given user
//
//See: http://docs.stormpath.com/rest/product-guide/#reset-an-accounts-password
func (app *Application) SendPasswordResetEmail(email string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	passwordResetPayload := make(map[string]string)
	passwordResetPayload["email"] = email

	err := client.post(buildAbsoluteURL(app.Href, "passwordResetTokens"), passwordResetPayload, passwordResetToken)

	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

//ValidatePasswordResetToken validates a password reset token
//
//See: http://docs.stormpath.com/rest/product-guide/#reset-an-accounts-password
func (app *Application) ValidatePasswordResetToken(token string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	err := client.get(buildAbsoluteURL(app.Href, "passwordResetTokens", token), passwordResetToken)

	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

//ResetPassword resets a user password based on the reset token
//
//See: http://docs.stormpath.com/rest/product-guide/#reset-an-accounts-password
func (app *Application) ResetPassword(token string, newPassword string) (*Account, error) {
	accountRef := &accountRef{}
	account := &Account{}

	resetPasswordPayload := make(map[string]string)
	resetPasswordPayload["password"] = newPassword

	err := client.post(buildAbsoluteURL(app.Href, "passwordResetTokens", token), resetPasswordPayload, accountRef)

	if err != nil {
		return nil, err
	}
	account.Href = accountRef.Account.Href

	return account, nil
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
func (app *Application) GetGroups(criteria Criteria) (*Groups, error) {
	groups := &Groups{}

	err := client.get(
		buildAbsoluteURL(app.Groups.Href, criteria.ToQueryString()),
		groups,
	)

	if err != nil {
		return nil, err
	}

	return groups, nil
}

//CreateIDSiteURL creates the IDSite URL for the application
func (app *Application) CreateIDSiteURL(options IDSiteOptions) (string, error) {
	nonce, _ := uuid.NewV4()

	if options.Path == "" {
		options.Path = "/"
	}

	claims := SSOTokenClaims{}
	claims.Id = nonce.String()
	claims.IssuedAt = time.Now().Unix()
	claims.Issuer = client.ClientConfiguration.APIKeyID
	claims.Subject = app.Href
	claims.State = options.State
	claims.Path = options.Path
	claims.CallbackURI = options.CallbackURL

	jwtString := JWT(claims, map[string]interface{}{})

	p, _ := url.Parse(app.Href)
	ssoURL := p.Scheme + "://" + p.Host + "/sso"

	if options.Logout {
		ssoURL = ssoURL + "/logout" + "?jwtRequest=" + jwtString
	} else {
		ssoURL = ssoURL + "?jwtRequest=" + jwtString
	}

	return ssoURL, nil
}

//HandleCallback handles the URL from an ID Site callback or SAML callback it parses the JWT token
//validates it and return an CallbackResult with the token info + the Account if the sub was given
func (app *Application) HandleCallback(URL string) (*CallbackResult, error) {
	result := &CallbackResult{}

	cbURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	jwtResponse := cbURL.Query().Get("jwtResponse")

	claims := &IDSiteAssertionTokenClaims{}

	ParseJWT(jwtResponse, claims)

	if claims.Audience != client.ClientConfiguration.APIKeyID {
		return nil, errors.New("ID Site invalid aud")
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, errors.New("ID Site JWT has expired")
	}

	if claims.Subject != "" {
		account, err := GetAccount(claims.Subject, MakeAccountCriteria())
		if err != nil {
			return nil, err
		}
		result.Account = account
	}

	result.State = claims.State
	result.Status = claims.Status

	return result, nil
}

//GetOAuthToken creates a OAuth2 token response for a given user credentials
func (app *Application) GetOAuthToken(username string, password string) (*OAuthResponse, error) {
	values := url.Values{
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
	}

	return app.getOAuthTokenCommon(values)
}

//GetOAuthTokenStormpathGrantType creates an OAuth2 token response for a given Stormpath token
func (app *Application) GetOAuthTokenStormpathGrantType(token string) (*OAuthResponse, error) {
	values := url.Values{
		"grant_type": {"stormpath_token"},
		"token":      {token},
	}

	return app.getOAuthTokenCommon(values)
}

func (app *Application) GetOAuthTokenClientCredentialsGrantType(apiKeyID, apiKeySecret string) (*OAuthResponse, error) {
	values := url.Values{
		"grant_type":   {"client_credentials"},
		"apiKeyId":     {apiKeyID},
		"apiKeySecret": {apiKeySecret},
	}

	return app.getOAuthTokenCommon(values)
}

//GetOAuthTokenSocialGrantType creates a OAuth2 token response for a given social provider token
func (app *Application) GetOAuthTokenSocialGrantType(providerID string, token string) (*OAuthResponse, error) {
	values := url.Values{
		"grant_type":  {"stormpath_social"},
		"providerId":  {providerID},
		"accessToken": {token},
	}

	return app.getOAuthTokenCommon(values)
}

func (app *Application) getOAuthTokenCommon(values url.Values) (*OAuthResponse, error) {
	response := &OAuthResponse{}

	err := client.postURLEncodedForm(
		buildAbsoluteURL(app.Href, "oauth/token"),
		values.Encode(),
		response,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

//RefreshOAuthToken refreshes an OAuth2 token using the provided refresh_token and returns a new OAuth reponse
func (app *Application) RefreshOAuthToken(refreshToken string) (*OAuthResponse, error) {
	response := &OAuthResponse{}

	values := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}
	body := &bytes.Buffer{}
	canonicalizeQueryString(body, values)

	err := client.postURLEncodedForm(
		buildAbsoluteURL(app.Href, "oauth/token"),
		body.String(),
		response,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

//ValidateToken against the application
func (app *Application) ValidateToken(token string) (*OAuthToken, error) {
	response := &OAuthToken{}

	err := client.get(
		buildAbsoluteURL(app.Href, "authTokens", token),
		response,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

//GetAPIKey retrives an APIKey from the given Application by its ID and optional criteria
func (app *Application) GetAPIKey(apiKeyID string, criteria APIKeyCriteria) (*APIKey, error) {
	apiKeys := &APIKeys{}

	err := client.get(buildAbsoluteURL(app.APIKeys.Href, criteria.idEq(apiKeyID).ToQueryString()), apiKeys)
	if err != nil {
		return nil, err
	}

	if len(apiKeys.Items) == 0 {
		return nil, fmt.Errorf("API Key not found")
	}

	return &apiKeys.Items[0], nil
}
