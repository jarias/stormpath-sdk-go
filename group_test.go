package stormpath

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupJsonMarshaling(t *testing.T) {
	t.Parallel()

	group := NewGroup("name")

	jsonData, _ := json.Marshal(group)

	assert.Equal(t, "{\"name\":\"name\"}", string(jsonData))
}

func TestGetGroup(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	group := createTestGroup(application, t)
	defer group.Delete()

	existingGroup, err := GetGroup(group.Href, MakeGroupCriteria())

	assert.NoError(t, err)
	assert.Equal(t, group, existingGroup)
}

func TestGetGroupNotFound(t *testing.T) {
	t.Parallel()

	noExistingGroup, err := GetGroup(GetClient().ClientConfiguration.BaseURL+"/groups/XXXX", MakeGroupCriteria())

	assert.Error(t, err)
	assert.Nil(t, noExistingGroup)
	assert.Equal(t, http.StatusNotFound, err.(Error).Status)
}

func TestGroupRefresh(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	group := createTestGroup(application, t)
	defer group.Delete()

	g := &Group{}
	g.Href = group.Href

	err := g.Refresh()

	assert.NoError(t, err)
	assert.Equal(t, group, g)
}

func TestGroupRefreshNotFound(t *testing.T) {
	t.Parallel()

	noExistingGroup := &Group{}
	noExistingGroup.Href = GetClient().ClientConfiguration.BaseURL + "/groups/XXXX"
	err := noExistingGroup.Refresh()

	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, err.(Error).Status)
}

func TestUpdateGroup(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	group := createTestGroup(application, t)
	defer group.Delete()

	group.Name = "newName" + randomName()
	group.Update()

	updatedGroup, err := GetGroup(group.Href, MakeGroupCriteria())

	assert.NoError(t, err)
	assert.Equal(t, group.Name, updatedGroup.Name)
}

func TestGetGroupAccountMemberships(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	group := createTestGroup(application, t)
	defer group.Delete()

	gm, err := group.GetGroupAccountMemberships(MakeGroupMemershipsCriteria())

	assert.NoError(t, err)
	assert.Empty(t, gm.Items)
	assert.Equal(t, 25, gm.GetLimit())
	assert.Equal(t, 0, gm.GetOffset())
}
