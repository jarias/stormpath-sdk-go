package stormpath

import (
	"fmt"
	"net/url"
)

//AccountCriteria represents the Criteria type for accounts
type AccountCriteria struct {
	baseCriteria
}

//MakeAccountCriteria creates a new AccountCriteria for a single Account resource
func MakeAccountCriteria() AccountCriteria {
	return AccountCriteria{baseCriteria{filter: url.Values{}}}
}

//MakeAccountsCriteria creates a new AccountCriteria for a AccountList resource
func MakeAccountsCriteria() AccountCriteria {
	return AccountCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//Pagination

func (c AccountCriteria) Limit(limit int) AccountCriteria {
	c.limit = limit
	return c
}

func (c AccountCriteria) Offset(offset int) AccountCriteria {
	c.offset = offset
	return c
}

//Filter related functions

//Possible filters:
//* givenName
//* surname
//* email
//* username
//* middleName
//* status
//* customData

//GivenNameEq adds the givenName equal filter to the given AccountCriteria
func (c AccountCriteria) GivenNameEq(givenName string) AccountCriteria {
	c.filter.Add("givenName", givenName)
	return c
}

//SurnameEq adds the surname equals filter to the given AccountCriteria
func (c AccountCriteria) SurnameEq(surname string) AccountCriteria {
	c.filter.Add("surname", surname)
	return c
}

//EmailEq adds the email equals filter to the given AccountCriteria
func (c AccountCriteria) EmailEq(email string) AccountCriteria {
	c.filter.Add("email", email)
	return c
}

//UsernameEq adds the username equals fitler to the given AccountCriteria
func (c AccountCriteria) UsernameEq(username string) AccountCriteria {
	c.filter.Add("username", username)
	return c
}

//MiddleNameEq adds the middleName equals filter to the given AccountCriteria
func (c AccountCriteria) MiddleNameEq(middleName string) AccountCriteria {
	c.filter.Add("middleName", middleName)
	return c
}

//StatusEq adds the status equals filter to the given AccountCriteria
func (c AccountCriteria) StatusEq(status string) AccountCriteria {
	c.filter.Add("status", status)
	return c
}

func (c AccountCriteria) CustomDataEq(k string, v string) AccountCriteria {
	c.filter.Add(fmt.Sprintf("customData.%s", k), v)
	return c
}

//Expansion related functions

//WithDirectory adds the directory expansion to the given AccountCriteria
func (c AccountCriteria) WithDirectory() AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "directory")
	return c
}

//WithCustomData adds the customData expansion to the given AccountCriteria
func (c AccountCriteria) WithCustomData() AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "customData")
	return c
}

//WithTenant adds the tenant expansion to the given AccountCriteria
func (c AccountCriteria) WithTenant() AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "tenant")
	return c
}

//WithGroups adds the groups expansion to the given AccountCriteria
func (c AccountCriteria) WithGroups(pageRequest PageRequest) AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("groups"))
	return c
}

//WithGroupMemberships adds the groupMembership expansion to the given AccountCriteria
func (c AccountCriteria) WithGroupMemberships(pageRequest PageRequest) AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("groupMemberships"))
	return c
}

//WithProviderData adds the providerData expansion to the given AccountCriteria
func (c AccountCriteria) WithProviderData() AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "providerData")
	return c
}

//WithAPIKeys adds the apiKeys expansion to the given AccountCriteria
func (c AccountCriteria) WithAPIKeys() AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "apiKeys")
	return c
}

//WithApplications adds the applications expansion to the given AccountCriteria
func (c AccountCriteria) WithApplications() AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "applications")
	return c
}
