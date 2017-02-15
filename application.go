package stormpath

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/nu7hatch/gouuid"
)

//Application is resource in Stormpath contains information about any real-world software that communicates with Stormpath via REST APIs. You control who may log in to an application by assigning (or ‘mapping’) one or more Directory, Group, or Organization resources (generically called Account Stores) to an Application resource. The Accounts in these associated Account Stores collectively form the application’s user base.
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

//Applications is the collection resource of applications.
//
//For more on Stormpath collection resources see: http://docs.stormpath.com/rest/product-guide/latest/reference.html#collection-resource
type Applications struct {
	collectionResource
	Items []Application `json:"items,omitempty"`
}

//CallbackResult is the parsed IDSite callback JWT token information and an optional account if the tocken contain one.
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

//CreateApplication creates a new application in Stormpath.
//It also passes the createDirectory param as true so the default directory would be created and mapped to the application.
func CreateApplication(app *Application) error {
	var extraParams = url.Values{}
	extraParams.Add("createDirectory", "true")

	return client.post(buildRelativeURL("applications", requestParams(extraParams)), app, app)
}

//GetApplication loads an application by href.
//It can optionally have its attributes expanded depending on the ApplicationCriteria value.
func GetApplication(href string, criteria ApplicationCriteria) (*Application, error) {
	application := &Application{}

	err := client.get(
		buildAbsoluteURL(href, criteria.toQueryString()),
		application,
	)
	if err != nil {
		return nil, err
	}

	return application, nil
}

//Refresh refreshes the application based on the latest state from Stormpath.
func (app *Application) Refresh() error {
	return client.get(app.Href, app)
}

//Update updates the application in Stormpath.
func (app *Application) Update() error {
	return client.post(app.Href, app, app)
}

//Purge deletes the application and all its account stores.
func (app *Application) Purge() error {
	accountStoreMappings, err := app.GetAccountStoreMappings(MakeApplicationAccountStoreMappingsCriteria())
	if err != nil {
		return err
	}

	for _, m := range accountStoreMappings.Items {
		client.delete(m.AccountStore.Href)
	}

	return app.Delete()
}

//GetAccountStoreMappings retrives the collection of all account store mappings associated with the Application.
//
//The collection can be filtered and/or paginated by passing the desire ApplicationAccountStoreMappingCriteria value
func (app *Application) GetAccountStoreMappings(criteria ApplicationAccountStoreMappingCriteria) (*ApplicationAccountStoreMappings, error) {
	accountStoreMappings := &ApplicationAccountStoreMappings{}

	err := client.get(
		buildAbsoluteURL(app.AccountStoreMappings.Href, criteria.toQueryString()),
		accountStoreMappings,
	)

	if err != nil {
		return nil, err
	}

	return accountStoreMappings, nil
}

//GetDefaultAccountStoreMapping retrieves the application default application account store mapping.
//
//It can optionally have its attributes expanded depending on the ApplicationAccountStoreMappingCriteria value.
func (app *Application) GetDefaultAccountStoreMapping(criteria ApplicationAccountStoreMappingCriteria) (*ApplicationAccountStoreMapping, error) {
	err := client.get(
		buildAbsoluteURL(app.DefaultAccountStoreMapping.Href, criteria.toQueryString()),
		app.DefaultAccountStoreMapping,
	)

	if err != nil {
		return nil, err
	}

	return app.DefaultAccountStoreMapping, nil
}

//RegisterAccount registers a new account into the application.
func (app *Application) RegisterAccount(account *Account) error {
	err := client.post(app.Accounts.Href, account, account)
	if err == nil {
		//Password should be cleanup so we don't keep an unhash password in memory
		account.Password = ""
	}
	return err
}

//RegisterSocialAccount registers a new account into the application using an external social provider Google, Facebook, GitHub or LinkedIn.
func (app *Application) RegisterSocialAccount(socialAccount *SocialAccount) (*Account, error) {
	account := &Account{}

	err := client.post(app.Accounts.Href, socialAccount, account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

//AuthenticateAccount authenticates an account against the application, using its username and password.
//It can also include an optional account store HREF, if the accountStoreHref is a zero value string, then it won't be used.
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
	//Refresh the account since we only get the account href from the loginAttempts endpoit.
	err = account.Refresh()
	if err != nil {
		return nil, err
	}

	return account, nil
}

//ResendVerificationEmail triggers a resend of the account verification email in Stormpath for a given email address.
//
//For more info on the Stormpath verification workflow see: http://docs.stormpath.com/rest/product-guide/latest/accnt_mgmt.html#how-to-verify-an-account-s-email
func (app *Application) ResendVerificationEmail(email string) error {
	resendVerificationEmailPayload := map[string]string{
		"login": email,
	}
	return client.post(buildAbsoluteURL(app.Href, "verificationEmails"), resendVerificationEmailPayload, nil)
}

//SendPasswordResetEmail triggers a send of the password reset email in Stormpath for a given email address.
//
//For more info on the Stormpath password reset workflow see: http://docs.stormpath.com/rest/product-guide/latest/accnt_mgmt.html#password-reset-flow
func (app *Application) SendPasswordResetEmail(email, accountStoreHref string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	passwordResetPayload := make(map[string]interface{})
	passwordResetPayload["email"] = email
	if accountStoreHref != "" {
		passwordResetPayload["accountStore"] = map[string]string{
			"href": accountStoreHref,
		}
	}

	err := client.post(buildAbsoluteURL(app.Href, "passwordResetTokens"), passwordResetPayload, passwordResetToken)

	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

//ValidatePasswordResetToken validates the given password reset token against Stormpath.
func (app *Application) ValidatePasswordResetToken(token string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	err := client.get(buildAbsoluteURL(app.Href, "passwordResetTokens", token), passwordResetToken)

	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

//ResetPassword resets a user password based on the reset password token.
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

//CreateGroup creates a new application group.
//Creating a group for an application automatically creates the proper account store mapping between the group and the application.
func (app *Application) CreateGroup(group *Group) error {
	return client.post(app.Groups.Href, group, group)
}

//GetGroups retrives the collection of all groups associated with the Application.
//
//The collection can be filtered and/or paginated by passing the desire GroupCriteria value.
func (app *Application) GetGroups(criteria GroupCriteria) (*Groups, error) {
	groups := &Groups{}

	err := client.get(
		buildAbsoluteURL(app.Groups.Href, criteria.toQueryString()),
		groups,
	)

	if err != nil {
		return nil, err
	}

	return groups, nil
}

//CreateIDSiteURL generates the IDSite URL for the application. This URL is used to initiate an IDSite workflow.
//You can pass an IDSiteOptions values to customize the IDSite workflow.
//
//For more information on Stormpath's IDSite feature see: http://docs.stormpath.com/rest/product-guide/latest/idsite.html
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

//HandleCallback handles the URL from an ID Site or SAML callback it parses the JWT token
//validates it and returns a CallbackResult, if the JWT was valid.
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

//GetOAuthToken creates a OAuth2 token response for an account, using the password grant type.
func (app *Application) GetOAuthToken(username string, password string) (*OAuthResponse, error) {
	values := url.Values{
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
	}

	return app.getOAuthTokenCommon(values)
}

//GetOAuthTokenStormpathGrantType creates an OAuth2 token response, for a given Stormpath JWT, using the stormpath_token grant type.
//This grant type is use together with IDSite.
//
//For more information on the stormpath_token grant type see: http://docs.stormpath.com/rest/product-guide/latest/idsite.html#exchanging-the-id-site-jwt-for-an-oauth-token
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

//GetOAuthTokenSocialGrantType creates a OAuth2 token response, for an account using a socical provider via the stormpath_social grant type.
//This grant type supports exchanging a social provider accessToken or code for a Stormpath access/refresh tokens.
//You can either pass an accessToken or code but not both, the one that is not used should be pass as a zero value string.
//If both values are empty or both values are not empty an error is return.
//
//For more information on the stormpath_social grant type see: http://docs.stormpath.com/rest/product-guide/latest/auth_n.html#social
func (app *Application) GetOAuthTokenSocialGrantType(providerID, accessToken, code string) (*OAuthResponse, error) {
	values := url.Values{
		"grant_type": {"stormpath_social"},
		"providerId": {providerID},
	}

	if accessToken == "" && code == "" {
		return nil, fmt.Errorf("You must either pass a valid accessToken or code.")
	}

	if accessToken != "" && code != "" {
		return nil, fmt.Errorf("You must either pass a valid accessToken or code but not both.")
	}

	if accessToken != "" {
		values.Add("accessToken", accessToken)
	}

	if code != "" {
		values.Add("code", code)
	}

	return app.getOAuthTokenCommon(values)
}

//RefreshOAuthToken creates an OAuth2 response using the refresh_token grant type.
func (app *Application) RefreshOAuthToken(refreshToken string) (*OAuthResponse, error) {
	values := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
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

//ValidateToken validates either an OAuth2 access or refresh token against Stormpath.
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

//GetAPIKey retrives an APIKey from the application by its ID.
//
//It can optionally have its attributes expanded depending on the APIKeyCriteria value.
func (app *Application) GetAPIKey(apiKeyID string, criteria APIKeyCriteria) (*APIKey, error) {
	apiKeys := &APIKeys{}

	err := client.get(buildAbsoluteURL(app.APIKeys.Href, criteria.IDEq(apiKeyID).toQueryString()), apiKeys)
	if err != nil {
		return nil, err
	}

	if len(apiKeys.Items) == 0 {
		return nil, fmt.Errorf("API Key not found")
	}

	return &apiKeys.Items[0], nil
}
