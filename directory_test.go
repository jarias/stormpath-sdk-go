package stormpath_test

import (
	"encoding/json"
	"testing"

	. "github.com/jarias/stormpath-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestDirectoryJsonMarshaling(t *testing.T) {
	t.Parallel()

	directory := NewDirectory("name")

	jsonData, _ := json.Marshal(directory)

	assert.Equal(t, "{\"name\":\"name\"}", string(jsonData))
}

func TestDeleteDirectory(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()

	err := directory.Delete()

	assert.NoError(t, err)
}

func TestGetAccountCreationPolicy(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()
	defer directory.Delete()

	policy, err := directory.GetAccountCreationPolicy()

	assert.NoError(t, err)
	assert.Equal(t, directory.AccountCreationPolicy, policy)
	assert.Equal(t, Disabled, policy.VerificationEmailStatus)
	assert.Equal(t, Disabled, policy.VerificationSuccessEmailStatus)
	assert.Equal(t, Disabled, policy.WelcomeEmailStatus)
}

func TestGetDirectoryEmptyGroupsCollection(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()
	defer directory.Delete()

	groups, err := directory.GetGroups(MakeGroupsCriteria())

	assert.NoError(t, err)
	assert.NotEmpty(t, groups.Href)
	assert.Equal(t, 0, groups.Offset)
	assert.Equal(t, 25, groups.Limit)
	assert.Empty(t, groups.Items)
}

func TestGetDirectoryEmptyAccountsCollection(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()
	defer directory.Delete()

	accounts, err := directory.GetAccounts(MakeAccountsCriteria())

	assert.NoError(t, err)
	assert.NotEmpty(t, accounts.Href)
	assert.Equal(t, 0, accounts.Offset)
	assert.Equal(t, 25, accounts.Limit)
	assert.Empty(t, accounts.Items)
}

func TestDirectoryCreateGroup(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()
	defer directory.Delete()

	group := newTestGroup()
	err := directory.CreateGroup(group)

	assert.NoError(t, err)
	assert.NotEmpty(t, group.Href)
}

func TestDirectoryRegisterAccount(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()
	defer directory.Delete()

	account := newTestAccount()
	err := directory.RegisterAccount(account)

	assert.NoError(t, err)
	assert.NotEmpty(t, account.Href)
}

func TestCreateGoogleDirectory(t *testing.T) {
	t.Parallel()

	directory := NewGoogleDirectory("google-"+randomName(), "ClientID", "ClientSercret", "http://localhost:8080")
	defer directory.Delete()

	err := tenant.CreateDirectory(directory)

	assert.NoError(t, err)
	assert.NotEmpty(t, directory.Href)

	d, err := GetDirectory(directory.Href, MakeDirectoryCriteria().WithProvider())

	assert.NoError(t, err)
	assert.Equal(t, Google, d.Provider.ProviderID)
	assert.Equal(t, directory.Provider.ClientID, d.Provider.ClientID)
	assert.Equal(t, directory.Provider.ClientSecret, d.Provider.ClientSecret)
	assert.Equal(t, directory.Provider.RedirectURI, d.Provider.RedirectURI)
}

func TestCreateLinkedInDirectory(t *testing.T) {
	t.Parallel()

	directory := NewLinkedInDirectory("linkedin-"+randomName(), "ClientID", "ClientSercret", "http://localhost:8080")
	defer directory.Delete()

	err := tenant.CreateDirectory(directory)

	assert.NoError(t, err)
	assert.NotEmpty(t, directory.Href)

	d, err := GetDirectory(directory.Href, MakeDirectoryCriteria().WithProvider())

	assert.NoError(t, err)
	assert.Equal(t, LinkedIn, d.Provider.ProviderID)
	assert.Equal(t, directory.Provider.ClientID, d.Provider.ClientID)
	assert.Equal(t, directory.Provider.ClientSecret, d.Provider.ClientSecret)
	assert.Equal(t, directory.Provider.RedirectURI, d.Provider.RedirectURI)
}

func TestCreateFacebookDirectory(t *testing.T) {
	t.Parallel()

	directory := NewFacebookDirectory("facebook-"+randomName(), "ClientID", "ClientSercret")
	defer directory.Delete()

	err := tenant.CreateDirectory(directory)

	assert.NoError(t, err)
	assert.NotEmpty(t, directory.Href)

	d, err := GetDirectory(directory.Href, MakeDirectoryCriteria().WithProvider())

	assert.NoError(t, err)
	assert.Equal(t, Facebook, d.Provider.ProviderID)
	assert.Equal(t, directory.Provider.ClientID, d.Provider.ClientID)
	assert.Equal(t, directory.Provider.ClientSecret, d.Provider.ClientSecret)
	assert.Empty(t, d.Provider.RedirectURI)
}

func TestCreateGithubDirectory(t *testing.T) {
	t.Parallel()

	directory := NewGithubDirectory("github-"+randomName(), "ClientID", "ClientSercret")
	defer directory.Delete()

	err := tenant.CreateDirectory(directory)

	assert.NoError(t, err)
	assert.NotEmpty(t, directory.Href)

	d, err := GetDirectory(directory.Href, MakeDirectoryCriteria().WithProvider())

	assert.NoError(t, err)
	assert.Equal(t, GitHub, d.Provider.ProviderID)
	assert.Equal(t, directory.Provider.ClientID, d.Provider.ClientID)
	assert.Equal(t, directory.Provider.ClientSecret, d.Provider.ClientSecret)
	assert.Empty(t, d.Provider.RedirectURI)
}
