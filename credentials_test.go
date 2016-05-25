package stormpath_test

import (
	"testing"

	. "github.com/jarias/stormpath-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestLoadCredentialsFromValidFile(t *testing.T) {
	t.Parallel()

	credentials, err := NewCredentialsFromFile("./test_files/apiKeys.properties")

	assert.NoError(t, err)
	assert.Equal(t, "APIKEY", credentials.ID)
	assert.Equal(t, "APISECRET", credentials.Secret)
}

func TestLoadCredentialsFromNoExistsFile(t *testing.T) {
	t.Parallel()

	credentials, err := NewCredentialsFromFile("./test_files/doesntexist.properties")

	assert.Error(t, err)
	assert.Equal(t, Credentials{}, credentials)
}

func TestLoadCredentialsFromEmptyFile(t *testing.T) {
	t.Parallel()

	credentials, err := NewCredentialsFromFile("./test_files/empty.properties")

	assert.NoError(t, err)
	assert.Empty(t, credentials.ID)
	assert.Empty(t, credentials.Secret)
}
