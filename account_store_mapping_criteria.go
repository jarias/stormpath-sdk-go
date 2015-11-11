package stormpath

import "net/url"

type AccountStoreMappingCriteria struct {
	baseCriteria
}

func MakeAccountStoreMappingCriteria() AccountStoreMappingCriteria {
	return AccountStoreMappingCriteria{baseCriteria{filter: url.Values{}}}
}

func MakeAccountStoreMappingsCriteria() AccountStoreMappingCriteria {
	return AccountStoreMappingCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//Expansion related functions

func (c AccountStoreMappingCriteria) WithApplication() AccountStoreMappingCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "application")
	return c
}

//TODO: Right not this function is disable cause the type of accountStore is resource and thus can't be
//expanded, probably need to change it interface and create a custom serializer depending on the href value
//directory, application or group
//func (c AccountStoreMappingCriteria) WithAccountStore() AccountStoreMappingCriteria {
//	c.expandedAttributes = append(c.expandedAttributes, "accountStore")
//	return c
//}
