package stormpath

import "net/url"

//ApplicationCriteria rerpresents the criteria object for an application or an applications collection.
type ApplicationCriteria struct {
	baseCriteria
}

//MakeApplicationCriteria an empty ApplicationCriteria for an application.
func MakeApplicationCriteria() ApplicationCriteria {
	return ApplicationCriteria{baseCriteria{filter: url.Values{}}}
}

//MakeApplicationsCriteria an empty ApplicationCriteria for an applications collection.
func MakeApplicationsCriteria() ApplicationCriteria {
	return ApplicationCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//Pagination

func (c ApplicationCriteria) Limit(limit int) ApplicationCriteria {
	c.limit = limit
	return c
}

func (c ApplicationCriteria) Offset(offset int) ApplicationCriteria {
	c.offset = offset
	return c
}

//Filter related functions

//Possible filters:
//* name
//* description
//* status

//NameEq adds the name filter to the given ApplicationCriteria
func (c ApplicationCriteria) NameEq(name string) ApplicationCriteria {
	c.filter.Add("name", name)
	return c
}

//DescriptionEq adds the description filter to the given ApplicationCriteria
func (c ApplicationCriteria) DescriptionEq(description string) ApplicationCriteria {
	c.filter.Add("description", description)
	return c
}

//StatusEq adds the status filter to the given ApplicationCriteria
func (c ApplicationCriteria) StatusEq(status string) ApplicationCriteria {
	c.filter.Add("status", status)
	return c
}

//Expansion related functions

//WithCustomData adds the customData expansion to the given ApplicationCriteria
func (c ApplicationCriteria) WithCustomData() ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "customData")
	return c
}

//WithAccounts adds the accounts expansion to the given ApplicationCriteria
func (c ApplicationCriteria) WithAccounts(pageRequest PageRequest) ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("accounts"))
	return c
}

//WithGroups adds the groups expansion to the given ApplicationCriteria
func (c ApplicationCriteria) WithGroups(pageRequest PageRequest) ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("groups"))
	return c
}

//WithTenant adds the tenant expansion to the given ApplicationCriteria
func (c ApplicationCriteria) WithTenant() ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "tenant")
	return c
}

//WithAccountStoreMappings adds the accountStoreMapping expansion to the given ApplicationCriteria
func (c ApplicationCriteria) WithAccountStoreMappings(pageRequest PageRequest) ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("accountStoreMappings"))
	return c
}

//WithDefaultAccountStoreMapping adds the defaultGroupStoreMapping expansion to the given ApplicationCriteria
func (c ApplicationCriteria) WithDefaultAccountStoreMapping() ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "defaultAccountStoreMapping")
	return c
}

//WithDefaultGroupStoreMapping adds the defaultGroupStoreMapping expansion to the given ApplicationCriteria
func (c ApplicationCriteria) WithDefaultGroupStoreMapping() ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "defaultGroupStoreMapping")
	return c
}

//WithRefreshTokens adds the refreshTokens expansion to the given ApplicationCriteria
func (c ApplicationCriteria) WithRefreshTokens(pageRequest PageRequest) ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("refreshTokens"))
	return c
}

//WithAccessTokens adds the accessTokens expansion to the given ApplicationCriteria
func (c ApplicationCriteria) WithAccessTokens(pageRequest PageRequest) ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("accessTokens"))
	return c
}
