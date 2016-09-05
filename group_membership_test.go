package stormpath

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGroupMembershipAccount(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	group := createTestGroup(application, t)
	defer group.Delete()

	account := createTestAccount(application, t)

	groupMembership, _ := account.AddToGroup(group)

	a, err := groupMembership.GetAccount(MakeAccountCriteria())

	assert.NoError(t, err)
	assert.Equal(t, account.Href, a.Href)
	assert.Equal(t, groupMembership.Account, a)
}

func TestGetGroupMembershipAccountNotFound(t *testing.T) {
	t.Parallel()

	gm := NewGroupMembership(client.ClientConfiguration.BaseURL+"/accounts/XXX", client.ClientConfiguration.BaseURL+"/groups/XXX")

	account, err := gm.GetAccount(MakeAccountCriteria())

	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Equal(t, http.StatusNotFound, err.(Error).Status)
}

func TestGetGroupMembershipGroup(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	group := createTestGroup(application, t)
	defer group.Delete()

	account := createTestAccount(application, t)

	groupMembership, _ := account.AddToGroup(group)

	g, err := groupMembership.GetGroup(MakeGroupCriteria())

	assert.NoError(t, err)
	assert.Equal(t, group.Href, g.Href)
	assert.Equal(t, groupMembership.Group, g)
}

func TestGetGroupMembershipGroupNotFound(t *testing.T) {
	t.Parallel()

	gm := NewGroupMembership(client.ClientConfiguration.BaseURL+"/accounts/XXX", client.ClientConfiguration.BaseURL+"/groups/XXX")

	group, err := gm.GetGroup(MakeAccountCriteria())

	assert.Error(t, err)
	assert.Nil(t, group)
	assert.Equal(t, http.StatusNotFound, err.(Error).Status)
}
