package stormpath

import "net/url"

//ApplicationAccountStoreMappingCriteria is the criteria type for the ApplicationAccountStoreMapping resource
type ApplicationAccountStoreMappingCriteria struct {
	baseCriteria
}

//OrganizationAccountStoreMappingCriteria is the criteria type for OrganizationAccountStoreMapping
type OrganizationAccountStoreMappingCriteria struct {
	baseCriteria
}

//MakeApplicationAccountStoreMappingCriteria creates a default ApplicationAccountStoreMappingCriteria for a single ApplicationAccountStoreMapping resource
func MakeApplicationAccountStoreMappingCriteria() ApplicationAccountStoreMappingCriteria {
	return ApplicationAccountStoreMappingCriteria{baseCriteria{filter: url.Values{}}}
}

//MakeApplicationAccountStoreMappingsCriteria creates a default ApplicationAccountStoreMappingCriteria for a ApplicationAccountStoreMappings collection resource
func MakeApplicationAccountStoreMappingsCriteria() ApplicationAccountStoreMappingCriteria {
	return ApplicationAccountStoreMappingCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//MakeOrganizationAccountStoreMappingCriteria creates a default OrganizationAccountStoreMappingCriteria for a single OrganizationAccountStoreMapping resource
func MakeOrganizationAccountStoreMappingCriteria() OrganizationAccountStoreMappingCriteria {
	return OrganizationAccountStoreMappingCriteria{baseCriteria{filter: url.Values{}}}
}

//MakeOrganizationAccountStoreMappingsCriteria creates a default OrganizationAccountStoreMappingCriteria for a OrganizationAccountStoreMappings collection resource
func MakeOrganizationAccountStoreMappingsCriteria() OrganizationAccountStoreMappingCriteria {
	return OrganizationAccountStoreMappingCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//Pagination
func (c ApplicationAccountStoreMappingCriteria) Limit(limit int) ApplicationAccountStoreMappingCriteria {
	c.limit = limit
	return c
}
func (c ApplicationAccountStoreMappingCriteria) Offset(offset int) ApplicationAccountStoreMappingCriteria {
	c.offset = offset
	return c
}

func (c OrganizationAccountStoreMappingCriteria) Limit(limit int) OrganizationAccountStoreMappingCriteria {
	c.limit = limit
	return c
}

func (c OrganizationAccountStoreMappingCriteria) Offset(offset int) OrganizationAccountStoreMappingCriteria {
	c.offset = offset
	return c
}

//Expansion related functions

//WithApplication adds the application expansion to the given ApplicationAccountStoreMappingCriteria
func (c ApplicationAccountStoreMappingCriteria) WithApplication() ApplicationAccountStoreMappingCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "application")
	return c
}

//WithOrganization adds the organization expansion to the given OrganizationAccountStoreMappingCriteria
func (c OrganizationAccountStoreMappingCriteria) WithOrganization() OrganizationAccountStoreMappingCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "organization")
	return c
}

//TODO: Right not this function is disable cause the type of accountStore is resource and thus can't be
//expanded, probably need to change it interface and create a custom serializer depending on the href value
//directory, application or group
//func (c AccountStoreMappingCriteria) WithAccountStore() AccountStoreMappingCriteria {
//	c.expandedAttributes = append(c.expandedAttributes, "accountStore")
//	return c
//}
