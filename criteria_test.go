package stormpath_test

import (
	"net/url"
	"testing"

	. "github.com/jarias/stormpath-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestEmptyAccountCriteria(t *testing.T) {
	t.Parallel()

	assert.Empty(t, MakeAccountCriteria().ToQueryString())
}

func TestPagedAccountCriteria(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "?limit=25&offset=0", MakeAccountCriteria().Offset(0).Limit(25).ToQueryString())
}

func TestAccountCriteriaFilters(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "?givenName=test", MakeAccountCriteria().GivenNameEq("test").ToQueryString())
	assert.Equal(t, "?surname=test", MakeAccountCriteria().SurnameEq("test").ToQueryString())
	assert.Equal(t, "?email=test", MakeAccountCriteria().EmailEq("test").ToQueryString())
	assert.Equal(t, "?username=test", MakeAccountCriteria().UsernameEq("test").ToQueryString())
	assert.Equal(t, "?middleName=test", MakeAccountCriteria().MiddleNameEq("test").ToQueryString())
	assert.Equal(t, "?status=test", MakeAccountCriteria().StatusEq("test").ToQueryString())
}

func TestAccountCriteriaExpansions(t *testing.T) {
	defaultPage := url.QueryEscape("(offset:0,limit:25)")

	assert.Equal(t, "?expand=directory", MakeAccountCriteria().WithDirectory().ToQueryString())
	assert.Equal(t, "?expand=customData", MakeAccountCriteria().WithCustomData().ToQueryString())
	assert.Equal(t, "?expand=tenant", MakeAccountCriteria().WithTenant().ToQueryString())
	assert.Equal(t, "?expand=groups"+defaultPage, MakeAccountCriteria().WithGroups(DefaultPageRequest).ToQueryString())
	assert.Equal(t, "?expand=groupMemberships"+defaultPage, MakeAccountCriteria().WithGroupMemberships(DefaultPageRequest).ToQueryString())
}

func TestAccountCriteriaExpansionsPaginAndFiltering(t *testing.T) {
	c := MakeAccountCriteria()

	str := c.WithDirectory().WithTenant().UsernameEq("test").Offset(2).Limit(40).ToQueryString()

	assert.Equal(t, "?expand=directory%2Ctenant&limit=40&offset=2&username=test", str)
}
