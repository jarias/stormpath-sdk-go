package stormpath

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganizationJsonMarshaling(t *testing.T) {
	t.Parallel()

	organization := NewOrganization("name", "nameKey")

	jsonData, _ := json.Marshal(organization)

	assert.Equal(t, "{\"name\":\"name\",\"nameKey\":\"nameKey\"}", string(jsonData))
}

func TestUpdateOrganization(t *testing.T) {
	t.Parallel()

	org := createTestOrganization(t)
	defer org.Delete()

	org.Name = "new-name" + randomName()
	err := org.Update()

	assert.NoError(t, err)

	updatedOrg, err := GetOrganization(org.Href, MakeOrganizationCriteria())

	assert.NoError(t, err)
	assert.Equal(t, org.Name, updatedOrg.Name)
}

func TestRefreshOrganization(t *testing.T) {
	t.Parallel()

	org := createTestOrganization(t)
	defer org.Delete()

	newName := "new-name" + randomName()
	org.Name = newName
	err := org.Refresh()

	assert.NoError(t, err)
	assert.NotEqual(t, newName, org.Name)
}

func TestGetOrganizationDefaultAccountStoreMapping(t *testing.T) {
	t.Parallel()

	org := createTestOrganization(t)
	defer org.Delete()

	directory := createTestDirectory(t)
	defer directory.Delete()

	mapping := NewOrganizationAccountStoreMapping(org.Href, directory.Href)
	mapping.IsDefaultAccountStore = true
	mapping.Save()

	org.Refresh()

	defaultMapping, err := org.GetDefaultAccountStoreMapping(MakeOrganizationAccountStoreMappingCriteria())

	assert.NoError(t, err)
	assert.Equal(t, org.Href, defaultMapping.Organization.Href)
	assert.Equal(t, directory.Href, defaultMapping.AccountStore.Href)
}

func TestOrganizationRegisterAccount(t *testing.T) {
	t.Parallel()

	org := createTestOrganization(t)
	defer org.Delete()

	directory := createTestDirectory(t)
	defer directory.Delete()

	mapping := NewOrganizationAccountStoreMapping(org.Href, directory.Href)
	mapping.IsDefaultAccountStore = true
	mapping.Save()

	account := newTestAccount()
	err := org.RegisterAccount(account)

	assert.NoError(t, err)
	assert.NotEmpty(t, account.Href)
}
