package stormpath

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAPIKey(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	apiKey, _ := account.CreateAPIKey()

	k := &APIKey{}
	err := GetAPIKeys(apiKey.Href, MakeAPIKeyCriteria(), k)

	assert.NoError(t, err)
	assert.Equal(t, apiKey, k)
}

func TestGetAPIKeys(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	apiKey1, _ := account.CreateAPIKey()
	apiKey2, _ := account.CreateAPIKey()

	keys := &APIKeys{}
	err := GetAPIKeys(account.APIKeys.Href, MakeAPIKeyCriteria(), keys)

	assert.NoError(t, err)
	assert.Equal(t, apiKey1, &keys.Items[0])
	assert.Equal(t, apiKey2, &keys.Items[1])
}

func TestDeleteAPIKey(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	apiKey, _ := account.CreateAPIKey()

	err := apiKey.Delete()

	assert.NoError(t, err)
	k := &APIKey{}

	err = GetAPIKeys(apiKey.Href, MakeAPIKeyCriteria(), k)

	assert.Error(t, err)
	assert.Equal(t, &APIKey{}, k)
	assert.Equal(t, http.StatusNotFound, err.(Error).Status)
}

func TestUpdateAPIKey(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	apiKey, _ := account.CreateAPIKey()

	apiKey.Status = Disabled
	err := apiKey.Update()

	assert.NoError(t, err)

	updatedAPIKey := &APIKey{}
	GetAPIKeys(apiKey.Href, MakeAPIKeyCriteria(), updatedAPIKey)
	assert.Equal(t, Disabled, updatedAPIKey.Status)
}
