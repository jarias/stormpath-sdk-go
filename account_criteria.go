package stormpath

import (
	"fmt"
	"net/url"
)

type AccountCriteria struct {
	baseCriteria
}

func MakeAccountCriteria() AccountCriteria {
	return AccountCriteria{baseCriteria{filter: url.Values{}}}
}

func MakeAccountsCriteria() AccountCriteria {
	return AccountCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
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

func (c AccountCriteria) GivenNameEq(givenName string) AccountCriteria {
	c.filter.Add("givenName", givenName)
	return c
}

func (c AccountCriteria) SurnameEq(surname string) AccountCriteria {
	c.filter.Add("surname", surname)
	return c
}

func (c AccountCriteria) EmailEq(email string) AccountCriteria {
	c.filter.Add("email", email)
	return c
}

func (c AccountCriteria) UsernameEq(username string) AccountCriteria {
	c.filter.Add("username", username)
	return c
}

func (c AccountCriteria) MiddleNameEq(middleName string) AccountCriteria {
	c.filter.Add("middleName", middleName)
	return c
}

func (c AccountCriteria) StatusEq(status string) AccountCriteria {
	c.filter.Add("status", status)
	return c
}

func (c AccountCriteria) CustomDataEq(k string, v string) AccountCriteria {
	c.filter.Add(fmt.Sprintf("customData.%s", k), v)
	return c
}

//Expansion related functions

func (c AccountCriteria) WithDirectory() AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "directory")
	return c
}

func (c AccountCriteria) WithCustomData() AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "customData")
	return c
}

func (c AccountCriteria) WithTenant() AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "tenant")
	return c
}

func (c AccountCriteria) WithGroups(pageRequest PageRequest) AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("groups"))
	return c
}

func (c AccountCriteria) WithGroupMemberships(pageRequest PageRequest) AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("groupMemberships"))
	return c
}

func (c AccountCriteria) WithProviderData() AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "providerData")
	return c
}

func (c AccountCriteria) WithAPIKeys() AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "apiKeys")
	return c
}

func (c AccountCriteria) WithApplications() AccountCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "applications")
	return c
}
