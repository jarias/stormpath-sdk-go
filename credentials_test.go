package stormpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadCredentialsFromFile(t *testing.T) {
	id, secret, err := loadCredentialsFromFile("./test_files/apiKeys.properties")

	assert.NoError(t, err)
	assert.Equal(t, "APIKEY", id)
	assert.Equal(t, "APISECRET", secret)
}

func TestLoadCredentialsFromEmptyFile(t *testing.T) {
	id, secret, err := loadCredentialsFromFile("./test_files/empty.properties")

	assert.NoError(t, err)
	assert.Equal(t, "", id)
	assert.Equal(t, "", secret)
}

func TestLoadCredentialsFromNoExistingFile(t *testing.T) {
	id, secret, err := loadCredentialsFromFile("./test_files/bla.properties")

	assert.Error(t, err)
	assert.Equal(t, "", id)
	assert.Equal(t, "", secret)
}
