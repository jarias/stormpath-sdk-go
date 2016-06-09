package stormpath_test

import (
	"encoding/json"
	"net/http"
	"testing"

	. "github.com/jarias/stormpath-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestAccountStoreMappingJsonMarshaling(t *testing.T) {
	t.Parallel()

	accountStoreMapping := NewAccountStoreMapping("http://appurl", "http://storeUrl")

	jsonData, _ := json.Marshal(accountStoreMapping)

	assert.Equal(t, "{\"application\":{\"href\":\"http://appurl\"},\"accountStore\":{\"href\":\"http://storeUrl\"}}", string(jsonData))
}

func TestSaveAccountStoreMapping(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	directory := createTestDirectory()
	defer directory.Delete()

	asm := NewAccountStoreMapping(application.Href, directory.Href)
	err := asm.Save()

	assert.NoError(t, err)
	assert.NotEmpty(t, asm.Href)
}

func TestSaveAccountStoreMappingApplicationNoExists(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()
	defer directory.Delete()

	asm := NewAccountStoreMapping(GetClient().ClientConfiguration.BaseURL+"applications/XXX", directory.Href)
	err := asm.Save()

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(Error).Status)
	assert.Equal(t, 2014, err.(Error).Code)
}

func TestSaveAccountStoreMappingDirectoryNoExists(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	asm := NewAccountStoreMapping(application.Href, GetClient().ClientConfiguration.BaseURL+"directories/XXX")
	err := asm.Save()

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(Error).Status)
	assert.Equal(t, 2014, err.(Error).Code)
}
