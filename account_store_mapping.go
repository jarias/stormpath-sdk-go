package stormpath

import "strings"

//ApplicationAccountStoreMapping represents an Stormpath account store mapping
//
//See: http://docs.stormpath.com/rest/product-guide/#account-store-mappings
type ApplicationAccountStoreMapping struct {
	resource
	ListIndex             *int         `json:"collectionResourceIndex,omitempty"`
	IsDefaultAccountStore bool         `json:"isDefaultAccountStore"`
	IsDefaultGroupStore   bool         `json:"isDefaultGroupStore"`
	Application           *Application `json:"application,omitempty"`
	AccountStore          *resource    `json:"accountStore,omitempty"`
}

type OrganizationAccountStoreMapping struct {
	resource
	ListIndex             *int          `json:"collectionResourceIndex,omitempty"`
	IsDefaultAccountStore bool          `json:"isDefaultAccountStore"`
	IsDefaultGroupStore   bool          `json:"isDefaultGroupStore"`
	Organization          *Organization `json:"organization,omitempty"`
	AccountStore          *resource     `json:"accountStore,omitempty"`
}

//ApplicationAccountStoreMappings represents a pages result of account store mappings
//
//See: http://docs.stormpath.com/rest/product-guide/#collectionResource-account-store-mappings
type ApplicationAccountStoreMappings struct {
	collectionResource
	Items []ApplicationAccountStoreMapping `json:"items,omitempty"`
}

type OrganizationAccountStoreMappings struct {
	collectionResource
	Items []OrganizationAccountStoreMapping `json:"items,omitempty"`
}

//NewAccountStoreMapping creates a new account store mappings
func NewApplicationAccountStoreMapping(applicationHref string, accountStoreHref string) *ApplicationAccountStoreMapping {
	app := Application{}
	app.Href = applicationHref
	return &ApplicationAccountStoreMapping{
		Application:  &app,
		AccountStore: &resource{Href: accountStoreHref},
	}
}

func NewOrganizationAccountStoreMapping(organizationHref string, accountStoreHref string) *OrganizationAccountStoreMapping {
	org := Organization{}
	org.Href = organizationHref
	return &OrganizationAccountStoreMapping{
		Organization: &org,
		AccountStore: &resource{Href: accountStoreHref},
	}
}

//Save saves the given account store mapping
func (mapping *ApplicationAccountStoreMapping) Save() error {
	url := buildRelativeURL("accountStoreMappings")
	if mapping.Href != "" {
		url = mapping.Href
	}

	return client.post(url, mapping, mapping)
}

func (mapping *OrganizationAccountStoreMapping) Save() error {
	url := buildRelativeURL("organizationAccountStoreMappings")
	if mapping.Href != "" {
		url = mapping.Href
	}

	return client.post(url, mapping, mapping)
}

func (mapping *ApplicationAccountStoreMapping) IsAccountStoreDirectory() bool {
	return strings.Contains(mapping.AccountStore.Href, "/directories/")
}

func (mapping *ApplicationAccountStoreMapping) IsAccountStoreGroup() bool {
	return strings.Contains(mapping.AccountStore.Href, "/groups/")
}

func (mapping *ApplicationAccountStoreMapping) IsAccountStoreOrganization() bool {
	return strings.Contains(mapping.AccountStore.Href, "/organizations/")
}

func (mapping *OrganizationAccountStoreMapping) IsAccountStoreDirectory() bool {
	return strings.Contains(mapping.AccountStore.Href, "/directories/")
}

func (mapping *OrganizationAccountStoreMapping) IsAccountStoreGroup() bool {
	return strings.Contains(mapping.AccountStore.Href, "/groups/")
}

func (mapping *OrganizationAccountStoreMapping) IsAccountStoreOrganization() bool {
	return strings.Contains(mapping.AccountStore.Href, "/organizations/")
}
