package stormpath

import "net/url"

type OrganizationCriteria struct {
	baseCriteria
}

func MakeOrganizationCriteria() OrganizationCriteria {
	return OrganizationCriteria{baseCriteria{filter: url.Values{}}}
}

func MakeOrganizationsCriteria() OrganizationCriteria {
	return OrganizationCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//Pagination

func (c OrganizationCriteria) Limit(limit int) OrganizationCriteria {
	c.limit = limit
	return c
}

func (c OrganizationCriteria) Offset(offset int) OrganizationCriteria {
	c.offset = offset
	return c
}

//Filter related functions

//Possible filters:
//* name
//* description
//* status

func (c OrganizationCriteria) NameEq(name string) OrganizationCriteria {
	c.filter.Add("name", name)
	return c
}

func (c OrganizationCriteria) DescriptionEq(description string) OrganizationCriteria {
	c.filter.Add("description", description)
	return c
}

func (c OrganizationCriteria) StatusEq(status string) OrganizationCriteria {
	c.filter.Add("status", status)
	return c
}

func (c OrganizationCriteria) NameKeyEq(status string) OrganizationCriteria {
	c.filter.Add("nameKey", status)
	return c
}

//Expansion related functions

func (c OrganizationCriteria) WithCustomData() OrganizationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "customData")
	return c
}

func (c OrganizationCriteria) WithAccounts(pageRequest PageRequest) OrganizationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("accounts"))
	return c
}

func (c OrganizationCriteria) WithGroups(pageRequest PageRequest) OrganizationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("groups"))
	return c
}

func (c OrganizationCriteria) WithTenant() OrganizationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "tenant")
	return c
}

func (c OrganizationCriteria) WithAccountStoreMappings(pageRequest PageRequest) OrganizationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("accountStoreMappings"))
	return c
}

func (c OrganizationCriteria) WithDefaultAccountStoreMapping() OrganizationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "defaultAccountStoreMapping")
	return c
}

func (c OrganizationCriteria) WithDefaultGroupStoreMapping() OrganizationCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "defaultGroupStoreMapping")
	return c
}
