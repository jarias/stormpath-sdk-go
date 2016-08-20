package stormpath

import "net/url"

type ApplicationAccountStoreMappingCriteria struct {
	baseCriteria
}

type OrganizationAccountStoreMappingCriteria struct {
	baseCriteria
}

func MakeApplicationAccountStoreMappingCriteria() ApplicationAccountStoreMappingCriteria {
	return ApplicationAccountStoreMappingCriteria{baseCriteria{filter: url.Values{}}}
}

func MakeApplicationAccountStoreMappingsCriteria() ApplicationAccountStoreMappingCriteria {
	return ApplicationAccountStoreMappingCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

func MakeOrganizationAccountStoreMappingCriteria() OrganizationAccountStoreMappingCriteria {
	return OrganizationAccountStoreMappingCriteria{baseCriteria{filter: url.Values{}}}
}

func MakeOrganizationAccountStoreMappingsCriteria() OrganizationAccountStoreMappingCriteria {
	return OrganizationAccountStoreMappingCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//Expansion related functions

func (c ApplicationAccountStoreMappingCriteria) WithApplication() ApplicationAccountStoreMappingCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "application")
	return c
}

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
