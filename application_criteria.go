package stormpath

import "net/url"

type ApplicationCriteria struct {
	baseCriteria
}

func MakeApplicationCriteria() ApplicationCriteria {
	return ApplicationCriteria{baseCriteria{filter: url.Values{}}}
}

func MakeApplicationsCriteria() ApplicationCriteria {
	return ApplicationCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//Filter related functions

//Possible filters:
//* name
//* description
//* status

func (c ApplicationCriteria) NameEq(name string) ApplicationCriteria {
	c.filter.Add("name", name)
	return c
}

func (c ApplicationCriteria) DescriptionEq(description string) ApplicationCriteria {
	c.filter.Add("description", description)
	return c
}

func (c ApplicationCriteria) StatusEq(status string) ApplicationCriteria {
	c.filter.Add("statu", status)
	return c
}

//Expansion related functions

func (c ApplicationCriteria) WithCustomData() ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "customData")
	return c
}

func (c ApplicationCriteria) WithAccounts(pageRequest PageRequest) ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("accounts"))
	return c
}

func (c ApplicationCriteria) WithGroups(pageRequest PageRequest) ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("groups"))
	return c
}

func (c ApplicationCriteria) WithTenant() ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "tenant")
	return c
}

func (c ApplicationCriteria) WithAccountStoreMappings(pageRequest PageRequest) ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("accountStoreMappings"))
	return c
}

func (c ApplicationCriteria) WithDefaultAccountStoreMapping() ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "defaultAccountStoreMapping")
	return c
}

func (c ApplicationCriteria) WithDefaultGroupStoreMapping() ApplicationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "defaultGroupStoreMapping")
	return c
}
