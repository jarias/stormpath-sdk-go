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

	k, err := GetAPIKey(apiKey.Href, MakeAPIKeyCriteria())

	assert.NoError(t, err)
	assert.Equal(t, apiKey, k)
}

func TestDeleteAPIKey(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)

	apiKey, _ := account.CreateAPIKey()

	err := apiKey.Delete()

	assert.NoError(t, err)

	k, err := GetAPIKey(apiKey.Href, MakeAPIKeyCriteria())

	assert.Error(t, err)
	assert.Nil(t, k)
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

	updatedAPIKey, _ := GetAPIKey(apiKey.Href, MakeAPIKeyCriteria())
	assert.Equal(t, Disabled, updatedAPIKey.Status)
}
