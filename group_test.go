package stormpath_test

import (
	"encoding/json"
	"testing"

	. "github.com/jarias/stormpath-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestGroupJsonMarshaling(t *testing.T) {
	group := NewGroup("name")

	jsonData, _ := json.Marshal(group)

	assert.Equal(t, "{\"name\":\"name\"}", string(jsonData))
}
