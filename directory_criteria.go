package stormpath

import (
	"net/url"
)

type DirectoryCriteria struct {
	baseCriteria
}

func MakeDirectoryCriteria() DirectoryCriteria {
	return DirectoryCriteria{baseCriteria{filter: url.Values{}}}
}

func MakeDirectoriesCriteria() DirectoryCriteria {
	return DirectoryCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//Pagination

func (c DirectoryCriteria) Limit(limit int) DirectoryCriteria {
	c.limit = limit
	return c
}

func (c DirectoryCriteria) Offset(offset int) DirectoryCriteria {
	c.offset = offset
	return c
}

//Filter related functions

//Possible filters:
//* name
//* description
//* status

func (c DirectoryCriteria) NameEq(name string) DirectoryCriteria {
	c.filter.Add(Name, name)
	return c
}

func (c DirectoryCriteria) DescriptionEq(description string) DirectoryCriteria {
	c.filter.Add(Description, description)
	return c
}

func (c DirectoryCriteria) StatusEq(status string) DirectoryCriteria {
	c.filter.Add(Status, status)
	return c
}

//Expansion related functions

func (c DirectoryCriteria) WithCustomData() DirectoryCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "customData")
	return c
}

func (c DirectoryCriteria) WithAccounts(pageRequest PageRequest) DirectoryCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("accounts"))
	return c
}

func (c DirectoryCriteria) WithGroups(pageRequest PageRequest) DirectoryCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("groups"))
	return c
}

func (c DirectoryCriteria) WithTenant() DirectoryCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "tenant")
	return c
}

func (c DirectoryCriteria) WithProvider() DirectoryCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "provider")
	return c
}

func (c DirectoryCriteria) WithAccountCreationPolicy() DirectoryCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "accountCreationPolicy")
	return c
}

func (c DirectoryCriteria) WithPasswordPolicy() DirectoryCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "passwordPolicy")
	return c
}
