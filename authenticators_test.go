package stormpath

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOAuthStormpathTokenAuthenticator(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)

	claims := GrantTypeStormpathTokenClaims{}
	claims.IssuedAt = time.Now().Unix()
	claims.Issuer = application.Href
	claims.Subject = account.Href
	claims.ExpiresAt = time.Now().Add(1 * time.Minute).Unix()
	claims.Status = "AUTHENTICATED"
	claims.Audience = client.ClientConfiguration.APIKeyID

	jwtString := JWT(
		claims,
		map[string]interface{}{
			"kid": client.ClientConfiguration.APIKeyID,
		},
	)

	authenticator := NewOAuthStormpathTokenAuthenticator(application)

	authResult, err := authenticator.Authenticate(jwtString)

	assert.NoError(t, err)
	assert.Equal(t, account.Href, authResult.GetAccount().Href)
}

func TestOAuthStormpathTokenAuthenticatorInvalidToken(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	authenticator := NewOAuthStormpathTokenAuthenticator(application)

	authResult, err := authenticator.Authenticate("I'm not a JWT token really")

	assert.Error(t, err)
	assert.EqualError(t, err, "Token is invalid")
	assert.Nil(t, authResult)
}

func TestOAuthClientCredentialsAuthenticator(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)
	apiKey, _ := account.CreateAPIKey()

	authenticator := NewOAuthClientCredentialsAuthenticator(application)

	authResult, err := authenticator.Authenticate(apiKey.ID, apiKey.Secret, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, authResult.AccessToken)
	assert.Empty(t, authResult.RefreshToken)
	assert.Equal(t, account.Href, authResult.GetAccount().Href)
}

func TestOAuthClientCredentialsAuthenticatorInvalidCredentials(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	authenticator := NewOAuthClientCredentialsAuthenticator(application)

	authResult, err := authenticator.Authenticate("foo", "bar", "")

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid_client")
	assert.Nil(t, authResult)
}

func TestOAuthClientCredentialsAuthenticatorScopeFactory(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	authenticator := NewOAuthClientCredentialsAuthenticator(application)
	authenticator.ScopeFactory = ScopeFactoryFunc(func(scope string) bool {
		return false
	})

	authResult, err := authenticator.Authenticate("foo", "bar", "bar")

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid_scope")
	assert.Nil(t, authResult)
}

func TestOAuthClientCredentialsAuthenticatorScopeFactoryValidScope(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)
	apiKey, _ := account.CreateAPIKey()

	authenticator := NewOAuthClientCredentialsAuthenticator(application)
	authenticator.ScopeFactory = ScopeFactoryFunc(func(scope string) bool {
		return true
	})

	authResult, err := authenticator.Authenticate(apiKey.ID, apiKey.Secret, "bar")

	assert.NoError(t, err)
	assert.NotEmpty(t, authResult.AccessToken)
	assert.Empty(t, authResult.RefreshToken)
	assert.Equal(t, account.Href, authResult.GetAccount().Href)
}

func TestBasicAuthenticator(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)
	apiKey, _ := account.CreateAPIKey()

	authenticator := NewBasicAuthenticator(application)

	authResult, err := authenticator.Authenticate(apiKey.ID, apiKey.Secret)

	assert.NoError(t, err)
	assert.Equal(t, account.Href, authResult.Account.Href)
}

func TestBasicAuthenticatorInvalidCredentials(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	authenticator := NewBasicAuthenticator(application)

	authResult, err := authenticator.Authenticate("foo", "bar")

	assert.Error(t, err)
	assert.Nil(t, authResult)
}

func TestBasicAuthenticatorDisabledCredentials(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)
	apiKey, _ := account.CreateAPIKey()
	apiKey.Status = Disabled
	apiKey.Update()

	authenticator := NewBasicAuthenticator(application)

	authResult, err := authenticator.Authenticate(apiKey.ID, apiKey.Secret)

	assert.Error(t, err)
	assert.EqualError(t, err, "API Key disabled")
	assert.Nil(t, authResult)
}

func TestBasicAuthenticatorDisabledAccount(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)
	apiKey, _ := account.CreateAPIKey()
	account.Status = Disabled
	account.Update()

	authenticator := NewBasicAuthenticator(application)

	authResult, err := authenticator.Authenticate(apiKey.ID, apiKey.Secret)

	assert.Error(t, err)
	assert.EqualError(t, err, "Account is disable")
	assert.Nil(t, authResult)
}

func TestBasicAuthenticatorInvalidSecret(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)
	apiKey, _ := account.CreateAPIKey()

	authenticator := NewBasicAuthenticator(application)

	authResult, err := authenticator.Authenticate(apiKey.ID, "foo")

	assert.Error(t, err)
	assert.EqualError(t, err, "Invalid API Key Secret")
	assert.Nil(t, authResult)
}

func TestOAuthPasswordAuthenticator(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)

	authenticator := NewOAuthPasswordAuthenticator(application)

	authResult, err := authenticator.Authenticate(account.Username, "1234567z!A89")

	assert.NoError(t, err)
	assert.NotEmpty(t, authResult.AccessToken)
	assert.NotEmpty(t, authResult.RefreshToken)
	assert.Equal(t, account.Href, authResult.GetAccount().Href)
}

func TestOAuthPasswordAuthenticatorInvalidCredentials(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)

	authenticator := NewOAuthPasswordAuthenticator(application)

	authResult, err := authenticator.Authenticate(account.Username, "foo")

	assert.Error(t, err)
	assert.EqualError(t, err, "Invalid username or password.")
	assert.Nil(t, authResult)
}

func TestOAuthRefreshTokenAuthenticator(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)
	oauthResponse, _ := application.GetOAuthToken(account.Username, "1234567z!A89")

	authenticator := NewOAuthRefreshTokenAuthenticator(application)

	authResult, err := authenticator.Authenticate(oauthResponse.RefreshToken)

	assert.NoError(t, err)
	assert.NotEmpty(t, authResult.AccessToken)
	assert.NotEmpty(t, authResult.RefreshToken)
	assert.Equal(t, account.Href, authResult.GetAccount().Href)
}

func TestOAuthRefreshTokenAuthenticatorInvalidToken(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	authenticator := NewOAuthRefreshTokenAuthenticator(application)

	authResult, err := authenticator.Authenticate("foo")

	assert.Error(t, err)
	assert.EqualError(t, err, "Token is invalid")
	assert.Nil(t, authResult)
}
