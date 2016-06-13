package stormpath

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccountRefreshTokens(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)

	application.GetOAuthToken(account.Username, "1234567z!A89")

	tokens, err := account.GetRefreshTokens(MakeOAuthTokensCriteria())

	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.Href)
	assert.Equal(t, 0, tokens.GetOffset())
	assert.Equal(t, 25, tokens.GetLimit())
	assert.NotEmpty(t, tokens.Items)
}

func TestGetAccountAccessTokens(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)

	application.GetOAuthToken(account.Username, "1234567z!A89")

	tokens, err := account.GetAccessTokens(MakeOAuthTokensCriteria())

	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.Href)
	assert.Equal(t, 0, tokens.GetOffset())
	assert.Equal(t, 25, tokens.GetLimit())
	assert.NotEmpty(t, tokens.Items)
}

func TestRevokeAccountAccessToken(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)

	application.GetOAuthToken(account.Username, "1234567z!A89")

	tokens, _ := account.GetAccessTokens(MakeOAuthTokensCriteria())

	token := tokens.Items[0]

	err := token.Delete()

	assert.NoError(t, err)
}

func TestAccountJsonMarshaling(t *testing.T) {
	t.Parallel()
	account := NewAccount("test@test.org", "123", "test@test.org", "test", "test")

	jsonData, err := json.Marshal(account)

	assert.NoError(t, err)
	assert.Equal(t, "{\"username\":\"test@test.org\",\"email\":\"test@test.org\",\"password\":\"123\",\"givenName\":\"test\",\"surname\":\"test\",\"emailVerificationToken\":null}", string(jsonData))
}

func TestGetAccountNoExists(t *testing.T) {
	t.Parallel()

	account, err := GetAccount(GetClient().ClientConfiguration.BaseURL+"/accounts/xxxxxx", MakeAccountCriteria())

	assert.Error(t, err)
	assert.Equal(t, 404, err.(Error).Status)
	assert.Nil(t, account)
}

func TestGetAccount(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	newAccount := createTestAccount(application)

	account, err := GetAccount(newAccount.Href, MakeAccountCriteria())

	assert.NoError(t, err)
	assert.Equal(t, newAccount, account)
}

func TestVerifyInvalidEmailToken(t *testing.T) {
	t.Parallel()

	account, err := VerifyEmailToken("token")

	assert.Error(t, err)
	assert.Equal(t, 404, err.(Error).Status)
	assert.Nil(t, account)
}

func TestVerifyValidEmailToken(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()
	defer directory.Delete()

	policy, _ := directory.GetAccountCreationPolicy()
	policy.VerificationEmailStatus = Enabled
	policy.Update()

	account := newTestAccount()
	directory.RegisterAccount(account)
	account.Refresh()

	assert.Equal(t, Unverified, account.Status)

	verifyAccount, err := VerifyEmailToken(GetToken(account.EmailVerificationToken.Href))
	verifyAccount.Refresh()

	assert.NoError(t, err)
	assert.Equal(t, account.Href, verifyAccount.Href)
	assert.Equal(t, Enabled, verifyAccount.Status)
}

func TestAccountUpdate(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)

	account.GivenName = "julio"
	err := account.Update()

	assert.NoError(t, err)

	updatedAccount, err := GetAccount(account.Href, MakeAccountCriteria())

	assert.NoError(t, err)
	assert.Equal(t, "julio", updatedAccount.GivenName)
}

func TestAccountDelete(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)

	err := account.Delete()

	assert.NoError(t, err)
}

func TestAddAccountToGroup(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	group := createTestGroup(application)
	defer group.Delete()

	account := createTestAccount(application)

	_, err := account.AddToGroup(group)

	assert.NoError(t, err)

	gm, err := account.GetGroupMemberships(MakeGroupMemershipCriteria().Offset(0).Limit(25))

	assert.NoError(t, err)
	assert.Len(t, gm.Items, 1)
}

func TestRemoveAccountFromGroup(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	var groupCountBefore int
	group := createTestGroup(application)
	defer group.Delete()

	account := createTestAccount(application)

	gm, _ := account.GetGroupMemberships(MakeAccountCriteria().Offset(0).Limit(25))
	groupCountBefore = len(gm.Items)

	account.AddToGroup(group)

	err := account.RemoveFromGroup(group)
	gm, _ = account.GetGroupMemberships(MakeAccountCriteria().Offset(0).Limit(25))

	assert.NoError(t, err)
	assert.Len(t, gm.Items, groupCountBefore)
}

func TestExpandGroupMembershipsAccount(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	group := createTestGroup(application)
	defer group.Delete()

	account := createTestAccount(application)

	groupMemberships, err := account.GetGroupMemberships(MakeGroupMemershipCriteria().WithAccount().Offset(0).Limit(25))

	assert.NoError(t, err)
	for _, gm := range groupMemberships.Items {
		assert.Equal(t, account, gm.Account)
		assert.NotEqual(t, group, gm.Group)
	}
}

func TestGetAccountCustomData(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)

	customData, err := account.GetCustomData()

	assert.NoError(t, err)
	assert.NotEmpty(t, customData)
}

func TestGetNoExistsAccountCustomData(t *testing.T) {
	t.Parallel()

	account := newTestAccount()
	account.Href = GetClient().ClientConfiguration.BaseURL + "/accounts/XXXX"

	customData, err := account.GetCustomData()

	assert.Error(t, err)
	assert.Equal(t, 404, err.(Error).Status)
	assert.Nil(t, customData)
}

func TestUpdateAccountCustomData(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	account := createTestAccount(application)

	customData, err := account.UpdateCustomData(map[string]interface{}{"custom": "data"})

	assert.NoError(t, err)
	assert.Equal(t, "data", customData["custom"])
}
