package stormpath

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetOAuthTokenValidAccount(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	oauthResponse, err := application.GetOAuthToken(account.Username, "1234567z!A89")

	assert.NoError(t, err)
	assert.NotNil(t, oauthResponse)
	assert.NotEmpty(t, oauthResponse.AccessToken)
	assert.NotEmpty(t, oauthResponse.RefreshToken)
	assert.Equal(t, 3600, oauthResponse.ExpiresIn)
}

func TestRefreshOAuthTokenValidAccount(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	oauthResponse, err := application.GetOAuthToken(account.Username, "1234567z!A89")

	assert.NoError(t, err)
	assert.NotNil(t, oauthResponse)
	assert.NotEmpty(t, oauthResponse.AccessToken)
	assert.NotEmpty(t, oauthResponse.RefreshToken)
	assert.Equal(t, 3600, oauthResponse.ExpiresIn)

	refreshOauthResponse, err := application.RefreshOAuthToken(oauthResponse.RefreshToken)

	assert.NoError(t, err)
	assert.NotNil(t, refreshOauthResponse)
	assert.NotEmpty(t, refreshOauthResponse.AccessToken)
	assert.NotEmpty(t, refreshOauthResponse.RefreshToken)
	assert.Equal(t, 3600, refreshOauthResponse.ExpiresIn)

	assert.NotEqual(t, oauthResponse.AccessToken, refreshOauthResponse.AccessToken)
}

func TestValidateOAuthAccessToken(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	response, err := application.GetOAuthToken(account.Username, "1234567z!A89")
	token, err := application.ValidateToken(response.AccessToken)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.JWT)
}

func TestValidateOAuthInvalidAccessToken(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	_, err := application.ValidateToken("anInvalidToken")

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(Error).Status)
}

func TestApplicationJsonMarshaling(t *testing.T) {
	t.Parallel()

	application := Application{Name: "name"}

	jsonData, _ := json.Marshal(application)

	assert.Equal(t, "{\"name\":\"name\"}", string(jsonData))
}

func TestUpdateApplication(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	application.Name = "new-name" + randomName()
	err := application.Update()

	assert.NoError(t, err)

	updatedApplication, err := GetApplication(application.Href, MakeApplicationCriteria())

	assert.NoError(t, err)
	assert.Equal(t, application.Name, updatedApplication.Name)
}

func TestApplicationRegisterAccount(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := newTestAccount()
	err := application.RegisterAccount(account)

	assert.NoError(t, err)
	assert.NotEmpty(t, account.Href)
}

func TestAuthenticateAccount(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	authenticatedAccount, err := application.AuthenticateAccount(account.Email, "1234567z!A89", "")

	assert.NoError(t, err)
	assert.Equal(t, account.Href, authenticatedAccount.Href)
	assert.Equal(t, account.GivenName, authenticatedAccount.GivenName)
	assert.Equal(t, account.Surname, authenticatedAccount.Surname)
	assert.Equal(t, account.Email, authenticatedAccount.Email)
}

func TestApplicationCreateInvalidGroup(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	err := application.CreateGroup(&Group{})

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(Error).Status)
}

func TestApplicationCreateGroup(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	group := newTestGroup()
	defer group.Delete()

	err := application.CreateGroup(group)

	assert.NoError(t, err)
	assert.NotEmpty(t, group.Href)
	assert.Equal(t, Enabled, group.Status)
}

func TestGetApplicationGroups(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	group := createTestGroup(application, t)
	defer group.Delete()

	groups, err := application.GetGroups(MakeGroupCriteria())

	assert.NoError(t, err)
	assert.NotEmpty(t, groups.Href)
	assert.Equal(t, 0, groups.GetOffset())
	assert.Equal(t, 25, groups.GetLimit())
	assert.NotEmpty(t, groups.Items)
}

func TestSendPasswordResetEmail(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	token, err := application.SendPasswordResetEmail(account.Email, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestSendPasswordResetEmailWithAccStore(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	//account store: Organization, Directory or Group
	token, err := application.SendPasswordResetEmail(account.Email, account.Directory.Href)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestResetPassword(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	token, _ := application.SendPasswordResetEmail(account.Email, "")

	re := regexp.MustCompile("[^\\/]+$")

	a, err := application.ResetPassword(re.FindString(token.Href), "8787987!kJKJdfW")

	assert.NoError(t, err)
	assert.Equal(t, account.Href, a.Href)
}

func TestResetPasswordWithAccStore(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	//account store: Organization, Directory or Group
	token, _ := application.SendPasswordResetEmail(account.Email, account.Directory.Href)

	re := regexp.MustCompile("[^\\/]+$")

	a, err := application.ResetPassword(re.FindString(token.Href), "8787987!kJKJdfW")

	assert.NoError(t, err)
	assert.Equal(t, account.Href, a.Href)
}

func TestValidatePasswordResetToken(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	token, _ := application.SendPasswordResetEmail(account.Email, "")

	re := regexp.MustCompile("[^\\/]+$")

	validatedToken, err := application.ValidatePasswordResetToken(re.FindString(token.Href))

	assert.NoError(t, err)
	assert.Equal(t, token.Href, validatedToken.Href)
}

func TestValidatePasswordResetTokenWithAccStore(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	//account store: Organization, Directory or Group
	token, _ := application.SendPasswordResetEmail(account.Email, account.Directory.Href)

	re := regexp.MustCompile("[^\\/]+$")

	validatedToken, err := application.ValidatePasswordResetToken(re.FindString(token.Href))

	assert.NoError(t, err)
	assert.Equal(t, token.Href, validatedToken.Href)
}

func TestValidateInvalidPasswordResetToken(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	_, err := application.ValidatePasswordResetToken("invalid token")

	assert.Error(t, err)
}

func TestCreateIDSiteURL(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	idSiteOptions := IDSiteOptions{
		CallbackURL: "http://localhost:8080",
	}

	idSiteURL, err := application.CreateIDSiteURL(idSiteOptions)

	u, _ := url.Parse(idSiteURL)

	assert.NoError(t, err)
	assert.Equal(t, "/sso", u.Path)
	assert.NotEmpty(t, u.Query)

	//Check Token
	jwtString := u.Query().Get("jwtRequest")

	claims := &SSOTokenClaims{}

	token := ParseJWT(jwtString, claims)

	assert.True(t, token.Valid)

	assert.Equal(t, "http://localhost:8080", claims.CallbackURI)
	assert.Equal(t, "", claims.State)
	assert.Equal(t, "/", claims.Path)
	assert.Equal(t, client.ClientConfiguration.APIKeyID, claims.Issuer)
	assert.Equal(t, application.Href, claims.Subject)
	assert.NotEmpty(t, claims.Id)
	assert.NotEmpty(t, claims.IssuedAt)
}

func TestCreateIDSiteLogoutURL(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	idSiteOptions := IDSiteOptions{
		CallbackURL: "http://localhost:8080",
		Logout:      true,
	}

	idSiteURL, err := application.CreateIDSiteURL(idSiteOptions)

	u, _ := url.Parse(idSiteURL)

	assert.NoError(t, err)
	assert.Equal(t, "/sso/logout", u.Path)
	assert.NotEmpty(t, u.Query)

	//Check Token
	jwtString := u.Query().Get("jwtRequest")

	claims := &SSOTokenClaims{}

	token := ParseJWT(jwtString, claims)

	assert.True(t, token.Valid)

	assert.Equal(t, "http://localhost:8080", claims.CallbackURI)
	assert.Equal(t, "", claims.State)
	assert.Equal(t, "/", claims.Path)
	assert.Equal(t, client.ClientConfiguration.APIKeyID, claims.Issuer)
	assert.Equal(t, application.Href, claims.Subject)
	assert.NotEmpty(t, claims.Id)
	assert.NotEmpty(t, claims.IssuedAt)
}

func TestGetApplicationDefaultAccountStoreMapping(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	defaultMapping, err := application.GetDefaultAccountStoreMapping(MakeApplicationAccountStoreMappingCriteria())

	assert.NoError(t, err)
	assert.Equal(t, application.Href, defaultMapping.Application.Href)
}

func TestGetOAuthTokenStormpathGrantType(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

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

	oauthResponse, err := application.GetOAuthTokenStormpathGrantType(jwtString)

	assert.NoError(t, err)
	assert.NotNil(t, oauthResponse)
}

func TestApplicationGetAPIKey(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	apiKey, err := account.CreateAPIKey()
	assert.NoError(t, err)
	assert.NotNil(t, apiKey)

	accountAPIKey, err := application.GetAPIKey(apiKey.ID, MakeAPIKeyCriteria())

	assert.NoError(t, err)
	assert.Equal(t, apiKey, accountAPIKey)
}
