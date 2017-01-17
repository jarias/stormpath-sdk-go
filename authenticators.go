package stormpath

import (
	"fmt"
	"net/http"
	"time"
)

type AuthResult interface {
	GetAccount() *Account
}

//AuthenticationResult base authentication result for all authenticators
type AuthenticationResult struct {
	Account *Account
}

type OAuthAccessTokenResult OAuthResponse

type OAuthClientCredentialsAuthenticationResult OAuthResponse

type StormpathAssertionAuthenticationResult CallbackResult

//Authenticator is the base authenticator type
//
//See https://github.com/stormpath/stormpath-sdk-spec/blob/master/specifications/authenticators.md
type Authenticator struct {
	Application *Application
}

//BasicAuthenticator will authenticate the API Key and Secret of a Stormpath Account object. Authentication should succeed only if the following are true:
//
// * The provided API Key and Secret exist for an account that is reachable by the application.
// * The API Key is not disabled.
// * The Account is not disabled.
type BasicAuthenticator Authenticator

//OAuthRequestAuthenticator should authenticate OAuth2 requests. It will eventually support authenticating all 4 OAuth2 grant types.
//
//Specifically, right now, this class will authenticate OAuth2 access tokens, as well as handle API key for access token exchanges using the OAuth2 client credentials grant type.
type OAuthRequestAuthenticator struct {
	Authenticator
	ScopeFactory ScopeFactoryFunc
	TTL          time.Duration
}

//OAuthBearerAuthenticator should authenticate OAuth2 bearer tokens only. The token is an access token JWT that has been created by Stormpath. The token may have been created by the client_credential or password_grant flow. This can be determined by looking at the kid property in the header of the JWT. Password grant JWTs will have a kid, but client credential JWTs will not.
type OAuthBearerAuthenticator Authenticator

//OAuthClientCredentialsAuthenticator this authenticator accepts an Account's API Key and Secret, and gives back an access token in response. The authenticator should follow the same authentication rules as the BasicAuthenticator. The end-user (account) can request scope, if the scope factory determines that this scope is permitted, then the scope should be added to the access token.
//
//This authenticator is responsible for creating the access token. The Stormpath REST API does not yet provide the client_credential grant on the appplication's /oauth/token endpoint.
type OAuthClientCredentialsAuthenticator struct {
	Authenticator
	ScopeFactory ScopeFactoryFunc
	TTL          time.Duration
}

type ScopeFactoryFunc func(string) bool

//OAuthPasswordAuthenticator this authenticator accepts an account's username and password, and returns an access token response that is obtained by posting the username and password to the application's /oauth/token endpoint with the grant_type=password parameter.
type OAuthPasswordAuthenticator Authenticator

//OAuthRefreshTokenAuthenticator this authenticator accepts a previously-issued refresh token and post's it to the application's /oauth/token endpoint with the grant_type=refresh_token parameter. The response is a new access token response.
type OAuthRefreshTokenAuthenticator Authenticator

//OAuthStormpathTokenAuthenticator this authenticator takes a Stormpath Token JWT and posts it to the application's /oauth/token endpoint, as grant_type=stormpath_token. The result is an OAuthAccessTokenResult.
type OAuthStormpathTokenAuthenticator Authenticator

//StormpathAssertionAuthenticator this authenticator will verify the a JWT from an ID Site or SAML callback. It should verify that:
//
// * The token is not expired
// * The signature can be verified
// * The claims body does not contain an err property.
type StormpathAssertionAuthenticator Authenticator

//NewBasicAuthenticator returns a BasicAuthenticator for the given application
func NewBasicAuthenticator(application *Application) BasicAuthenticator {
	return BasicAuthenticator{application}
}

//Authenticate authenticates the given account APIKey and APISecret
func (a BasicAuthenticator) Authenticate(accountAPIKey, accountAPISecret string) (*AuthenticationResult, error) {
	apiKey, err := a.Application.GetAPIKey(accountAPIKey, MakeAPIKeysCriteria().WithAccount())
	if err != nil {
		return nil, err
	}

	if apiKey.Secret != accountAPISecret {
		return nil, fmt.Errorf("Invalid API Key Secret")
	}

	if apiKey.Status == Disabled {
		return nil, fmt.Errorf("API Key disabled")
	}

	if apiKey.Account.Status == Disabled {
		return nil, fmt.Errorf("Account is disable")
	}

	return &AuthenticationResult{Account: apiKey.Account}, nil
}

func NewOAuthRequestAuthenticator(application *Application) OAuthRequestAuthenticator {
	authenticator := OAuthRequestAuthenticator{}
	authenticator.Application = application
	return authenticator
}

func (a OAuthRequestAuthenticator) Authenticate(r *http.Request) (*OAuthAccessTokenResult, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	grantType := r.Form.Get("grant_type")

	switch grantType {
	case "password":
		authResult, err := NewOAuthPasswordAuthenticator(a.Application).Authenticate(r.Form.Get("username"), r.Form.Get("password"))
		if err != nil {
			return nil, err
		}
		return authResult, nil
	case "client_credentials":
		accountAPIKeyID, accountAPIKeyScret, ok := r.BasicAuth()
		if !ok {
			return nil, fmt.Errorf("invalid_client")
		}
		oauthResponse, err := NewOAuthClientCredentialsAuthenticator(a.Application).Authenticate(accountAPIKeyID, accountAPIKeyScret, r.Form.Get("scope"))
		if err != nil {
			return nil, err
		}
		result := OAuthAccessTokenResult(*oauthResponse)
		return &result, nil
	case "refresh_token":
		authResult, err := NewOAuthRefreshTokenAuthenticator(a.Application).Authenticate(r.Form.Get("refresh_token"))
		if err != nil {
			return nil, err
		}
		return authResult, nil
	case "stormpath_social":
		oauthResponse, err := a.Application.GetOAuthTokenSocialGrantType(r.Form.Get("providerId"), r.Form.Get("accessToken"))
		if err != nil {
			return nil, err
		}
		result := OAuthAccessTokenResult(*oauthResponse)
		return &result, nil
	}

	return nil, fmt.Errorf("unsupported_grant_type")
}

func NewOAuthClientCredentialsAuthenticator(application *Application) OAuthClientCredentialsAuthenticator {
	authenticator := OAuthClientCredentialsAuthenticator{}
	authenticator.Application = application
	authenticator.TTL = 3600 * time.Second
	return authenticator
}

func (a OAuthClientCredentialsAuthenticator) Authenticate(accountAPIKeyID, accountAPIKeySecret, scope string) (*OAuthClientCredentialsAuthenticationResult, error) {
	if a.ScopeFactory != nil {
		if !a.ScopeFactory(scope) {
			return nil, fmt.Errorf("invalid_scope")
		}
	}

	oAuthResponse, err := a.Application.GetOAuthTokenClientCredentialsGrantType(accountAPIKeyID, accountAPIKeySecret)
	if err != nil {
		return nil, fmt.Errorf("invalid_client")
	}

	oauthResult := OAuthClientCredentialsAuthenticationResult(*oAuthResponse)

	return &oauthResult, nil
}

func NewOAuthPasswordAuthenticator(application *Application) OAuthPasswordAuthenticator {
	return OAuthPasswordAuthenticator{application}
}

func (a OAuthPasswordAuthenticator) Authenticate(username, password string) (*OAuthAccessTokenResult, error) {
	oauthResponse, err := a.Application.GetOAuthToken(username, password)
	if err != nil {
		return nil, err
	}

	authResult := OAuthAccessTokenResult(*oauthResponse)
	return &authResult, nil
}

func NewOAuthRefreshTokenAuthenticator(application *Application) OAuthRefreshTokenAuthenticator {
	return OAuthRefreshTokenAuthenticator{application}
}

func (a OAuthRefreshTokenAuthenticator) Authenticate(refreshToken string) (*OAuthAccessTokenResult, error) {
	oauthResponse, err := a.Application.RefreshOAuthToken(refreshToken)
	if err != nil {
		return nil, err
	}

	authResult := OAuthAccessTokenResult(*oauthResponse)
	return &authResult, nil
}

func NewOAuthStormpathTokenAuthenticator(application *Application) OAuthStormpathTokenAuthenticator {
	return OAuthStormpathTokenAuthenticator{application}
}

func (a OAuthStormpathTokenAuthenticator) Authenticate(stormpathJWT string) (*OAuthAccessTokenResult, error) {
	oauthResponse, err := a.Application.GetOAuthTokenStormpathGrantType(stormpathJWT)
	if err != nil {
		return nil, err
	}

	authResult := OAuthAccessTokenResult(*oauthResponse)
	return &authResult, nil
}

func NewStormpathAssertionAuthenticator(application *Application) StormpathAssertionAuthenticator {
	return StormpathAssertionAuthenticator{application}
}

func (a StormpathAssertionAuthenticator) Authenticate(stormpathJWT string) (*StormpathAssertionAuthenticationResult, error) {
	callbackResponse, err := a.Application.HandleCallback("http://fake?jwtResponse=" + stormpathJWT)
	if err != nil {
		return nil, err
	}

	authResult := StormpathAssertionAuthenticationResult(*callbackResponse)
	return &authResult, nil
}

func NewOAuthBearerAuthenticator(application *Application) OAuthBearerAuthenticator {
	return OAuthBearerAuthenticator{application}
}

func (a OAuthBearerAuthenticator) Authenticate(accessTokenJWT string) (*AuthenticationResult, error) {
	oauthToken, err := a.Application.ValidateToken(accessTokenJWT)
	if err != nil {
		return nil, err
	}

	if oauthToken.ExpandedJWT.Header.STT != "access" {
		//This is a refresh token
		return nil, fmt.Errorf("can't use refresh token as access token")
	}

	return &AuthenticationResult{oauthToken.Account}, nil
}

func (ar *AuthenticationResult) GetAccount() *Account {
	account, err := GetAccount(ar.Account.Href, MakeAccountCriteria().WithProviderData().WithDirectory())
	if err != nil {
		return nil
	}
	return account
}

func (ar *OAuthAccessTokenResult) GetAccount() *Account {
	claims := &AccessTokenClaims{}

	ParseJWT(ar.AccessToken, claims)

	account, err := GetAccount(claims.Subject, MakeAccountCriteria().WithProviderData().WithDirectory())
	if err != nil {
		return nil
	}

	return account
}

func (ar *OAuthClientCredentialsAuthenticationResult) GetAccount() *Account {
	claims := &AccessTokenClaims{}

	ParseJWT(ar.AccessToken, claims)

	account, err := GetAccount(claims.Subject, MakeAccountCriteria().WithProviderData().WithDirectory())
	if err != nil {
		return nil
	}

	return account
}

func (ar *StormpathAssertionAuthenticationResult) GetAccount() *Account {
	account, err := GetAccount(ar.Account.Href, MakeAccountCriteria().WithProviderData().WithDirectory())
	if err != nil {
		return nil
	}
	return account
}
