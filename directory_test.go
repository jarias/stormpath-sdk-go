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
