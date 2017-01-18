package stormpath

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultPage = url.QueryEscape("(offset:0,limit:25)")

func TestDirectoryCriteria(t *testing.T) {
	t.Parallel()

	cases := []struct {
		expected string
		actual   DirectoryCriteria
	}{
		{"", MakeDirectoryCriteria()},
		{"?name=test", MakeDirectoryCriteria().NameEq("test")},
		{"?description=test", MakeDirectoryCriteria().DescriptionEq("test")},
		{"?status=" + Enabled, MakeDirectoryCriteria().StatusEq(Enabled)},
		{"?expand=customData", MakeDirectoryCriteria().WithCustomData()},
		{"?expand=accounts%28offset%3A0%2Climit%3A25%29", MakeDirectoryCriteria().WithAccounts(DefaultPageRequest)},
		{"?expand=groups%28offset%3A0%2Climit%3A25%29", MakeDirectoryCriteria().WithGroups(DefaultPageRequest)},
		{"?expand=tenant", MakeDirectoryCriteria().WithTenant()},
		{"?expand=provider", MakeDirectoryCriteria().WithProvider()},
		{"?expand=accountCreationPolicy", MakeDirectoryCriteria().WithAccountCreationPolicy()},
		{"?expand=passwordPolicy", MakeDirectoryCriteria().WithPasswordPolicy()},
		{"?expand=tenant%2Cprovider%2CcustomData", MakeDirectoryCriteria().WithTenant().WithProvider().WithCustomData()},
		{"?expand=tenant&name=test", MakeDirectoryCriteria().WithTenant().NameEq("test")},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual.toQueryString())
	}
}

func TestDirectoriesCriteria(t *testing.T) {
	t.Parallel()

	cases := []struct {
		expected string
		actual   DirectoryCriteria
	}{
		{"?limit=25&offset=0", MakeDirectoriesCriteria()},
		{"?name=test&limit=25&offset=0", MakeDirectoriesCriteria().NameEq("test")},
		{"?description=test&limit=25&offset=0", MakeDirectoriesCriteria().DescriptionEq("test")},
		{"?status=" + Enabled + "&limit=25&offset=0", MakeDirectoriesCriteria().StatusEq(Enabled)},
		{"?expand=customData&limit=25&offset=0", MakeDirectoriesCriteria().WithCustomData()},
		{"?expand=accounts%28offset%3A0%2Climit%3A25%29&limit=25&offset=0", MakeDirectoriesCriteria().WithAccounts(DefaultPageRequest)},
		{"?expand=groups%28offset%3A0%2Climit%3A25%29&limit=25&offset=0", MakeDirectoriesCriteria().WithGroups(DefaultPageRequest)},
		{"?expand=tenant&limit=25&offset=0", MakeDirectoriesCriteria().WithTenant()},
		{"?expand=provider&limit=25&offset=0", MakeDirectoriesCriteria().WithProvider()},
		{"?expand=accountCreationPolicy&limit=25&offset=0", MakeDirectoriesCriteria().WithAccountCreationPolicy()},
		{"?expand=passwordPolicy&limit=25&offset=0", MakeDirectoriesCriteria().WithPasswordPolicy()},
		{"?expand=tenant%2Cprovider%2CcustomData&limit=25&offset=0", MakeDirectoriesCriteria().WithTenant().WithProvider().WithCustomData()},
		{"?expand=tenant&name=test&limit=25&offset=0", MakeDirectoriesCriteria().WithTenant().NameEq("test")},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual.toQueryString())
	}
}

func TestAccountCriteria(t *testing.T) {
	t.Parallel()

	cases := []struct {
		expected string
		actual   AccountCriteria
	}{
		{"", MakeAccountCriteria()},
		{"?givenName=test", MakeAccountCriteria().GivenNameEq("test")},
		{"?surname=test", MakeAccountCriteria().SurnameEq("test")},
		{"?email=test", MakeAccountCriteria().EmailEq("test")},
		{"?username=test", MakeAccountCriteria().UsernameEq("test")},
		{"?middleName=test", MakeAccountCriteria().MiddleNameEq("test")},
		{"?status=test", MakeAccountCriteria().StatusEq("test")},
		{"?expand=tenant", MakeAccountCriteria().WithTenant()},
		{"?expand=directory", MakeAccountCriteria().WithDirectory()},
		{"?expand=customData", MakeAccountCriteria().WithCustomData()},
		{"?expand=tenant", MakeAccountCriteria().WithTenant()},
		{"?expand=groups" + defaultPage, MakeAccountCriteria().WithGroups(DefaultPageRequest)},
		{"?expand=groupMemberships" + defaultPage, MakeAccountCriteria().WithGroupMemberships(DefaultPageRequest)},
		{"?expand=tenant&givenName=test", MakeAccountCriteria().WithTenant().GivenNameEq("test")},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual.toQueryString())
	}
}

func TestGroupCriteria(t *testing.T) {
	t.Parallel()

	cases := []struct {
		expected string
		actual   GroupCriteria
	}{
		{"", MakeGroupCriteria()},
		{"?name=test", MakeGroupCriteria().NameEq("test")},
		{"?description=test", MakeGroupCriteria().DescriptionEq("test")},
		{"?status=" + Enabled, MakeGroupCriteria().StatusEq(Enabled)},
		{"?expand=customData", MakeGroupCriteria().WithCustomData()},
		{"?expand=accounts" + defaultPage, MakeGroupCriteria().WithAccounts(DefaultPageRequest)},
		{"?expand=tenant", MakeGroupCriteria().WithTenant()},
		{"?expand=directory", MakeGroupCriteria().WithDirectory()},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual.toQueryString())
	}
}

func TestApplicationCriteria(t *testing.T) {
	t.Parallel()

	cases := []struct {
		expected string
		actual   ApplicationCriteria
	}{
		{"", MakeApplicationCriteria()},
		{"?name=test", MakeApplicationCriteria().NameEq("test")},
		{"?description=test", MakeApplicationCriteria().DescriptionEq("test")},
		{"?status=" + Enabled, MakeApplicationCriteria().StatusEq(Enabled)},
		{"?expand=customData", MakeApplicationCriteria().WithCustomData()},
		{"?expand=accounts" + defaultPage, MakeApplicationCriteria().WithAccounts(DefaultPageRequest)},
		{"?expand=tenant", MakeApplicationCriteria().WithTenant()},
		{"?expand=groups" + defaultPage, MakeApplicationCriteria().WithGroups(DefaultPageRequest)},
		{"?expand=defaultAccountStoreMapping", MakeApplicationCriteria().WithDefaultAccountStoreMapping()},
		{"?expand=defaultGroupStoreMapping", MakeApplicationCriteria().WithDefaultGroupStoreMapping()},
		{"?expand=refreshTokens" + defaultPage, MakeApplicationCriteria().WithRefreshTokens(DefaultPageRequest)},
		{"?expand=accessTokens" + defaultPage, MakeApplicationCriteria().WithAccessTokens(DefaultPageRequest)},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual.toQueryString())
	}
}

func TestOrganizationCriteria(t *testing.T) {
	t.Parallel()

	cases := []struct {
		expected string
		actual   OrganizationCriteria
	}{
		{"", MakeOrganizationCriteria()},
		{"?name=test", MakeOrganizationCriteria().NameEq("test")},
		{"?description=test", MakeOrganizationCriteria().DescriptionEq("test")},
		{"?status=" + Enabled, MakeOrganizationCriteria().StatusEq(Enabled)},
		{"?expand=customData", MakeOrganizationCriteria().WithCustomData()},
		{"?expand=accounts" + defaultPage, MakeOrganizationCriteria().WithAccounts(DefaultPageRequest)},
		{"?expand=tenant", MakeOrganizationCriteria().WithTenant()},
		{"?expand=groups" + defaultPage, MakeOrganizationCriteria().WithGroups(DefaultPageRequest)},
		{"?expand=defaultAccountStoreMapping", MakeOrganizationCriteria().WithDefaultAccountStoreMapping()},
		{"?expand=defaultGroupStoreMapping", MakeOrganizationCriteria().WithDefaultGroupStoreMapping()},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual.toQueryString())
	}
}

func TestGroupMembershipCriteria(t *testing.T) {
	t.Parallel()

	cases := []struct {
		expected string
		actual   GroupMembershipCriteria
	}{
		{"", MakeGroupMemershipCriteria()},
		{"?expand=group", MakeGroupMemershipCriteria().WithGroup()},
		{"?expand=account", MakeGroupMemershipCriteria().WithAccount()},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual.toQueryString())
	}
}

func TestAPIKeyCriteria(t *testing.T) {
	t.Parallel()

	cases := []struct {
		expected string
		actual   APIKeyCriteria
	}{
		{"", MakeAPIKeyCriteria()},
		{"?expand=tenant", MakeAPIKeyCriteria().WithTenant()},
		{"?expand=account", MakeAPIKeyCriteria().WithAccount()},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual.toQueryString())
	}
}

func TestApplicationAccountStoreMappingCriteria(t *testing.T) {
	t.Parallel()

	cases := []struct {
		expected string
		actual   ApplicationAccountStoreMappingCriteria
	}{
		{"", MakeApplicationAccountStoreMappingCriteria()},
		{"?expand=application", MakeApplicationAccountStoreMappingCriteria().WithApplication()},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual.toQueryString())
	}
}

func TestOrganizationAccountStoreMappingCriteria(t *testing.T) {
	t.Parallel()

	cases := []struct {
		expected string
		actual   OrganizationAccountStoreMappingCriteria
	}{
		{"", MakeOrganizationAccountStoreMappingCriteria()},
		{"?expand=organization", MakeOrganizationAccountStoreMappingCriteria().WithOrganization()},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual.toQueryString())
	}
}
