package stormpath

import "strings"

//ApplicationAccountStoreMapping represents an Stormpath account store mapping
//
//See: https://docs.stormpath.com/rest/product-guide/latest/reference.html#account-store-mapping
type ApplicationAccountStoreMapping struct {
	resource
	ListIndex             *int         `json:"collectionResourceIndex,omitempty"`
	IsDefaultAccountStore bool         `json:"isDefaultAccountStore"`
	IsDefaultGroupStore   bool         `json:"isDefaultGroupStore"`
	Application           *Application `json:"application,omitempty"`
	AccountStore          *resource    `json:"accountStore,omitempty"`
}

//OrganizationAccountStoreMapping represents an Stormpath account store mapping for an Organization resource
//
//See: https://docs.stormpath.com/rest/product-guide/latest/reference.html?organization-account-store-mapping-operations
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

//OrganizationAccountStoreMappings represents a collection of OrganizationAccountStoreMapping
type OrganizationAccountStoreMappings struct {
	collectionResource
	Items []OrganizationAccountStoreMapping `json:"items,omitempty"`
}

//NewApplicationAccountStoreMapping creates a new account store mapping for the Application resource
func NewApplicationAccountStoreMapping(applicationHref string, accountStoreHref string) *ApplicationAccountStoreMapping {
	app := Application{}
	app.Href = applicationHref
	return &ApplicationAccountStoreMapping{
		Application:  &app,
		AccountStore: &resource{Href: accountStoreHref},
	}
}

//NewOrganizationAccountStoreMapping creates a new account mapping for the Organization resource
func NewOrganizationAccountStoreMapping(organizationHref string, accountStoreHref string) *OrganizationAccountStoreMapping {
	org := Organization{}
	org.Href = organizationHref
	return &OrganizationAccountStoreMapping{
		Organization: &org,
		AccountStore: &resource{Href: accountStoreHref},
	}
}

//Save saves the given ApplicationAccountStoreMapping
func (mapping *ApplicationAccountStoreMapping) Save() error {
	url := buildRelativeURL("accountStoreMappings")
	if mapping.Href != "" {
		url = mapping.Href
	}

	return client.post(url, mapping, mapping)
}

//Save saves the given OrganizationAccountStoreMapping
func (mapping *OrganizationAccountStoreMapping) Save() error {
	url := buildRelativeURL("organizationAccountStoreMappings")
	if mapping.Href != "" {
		url = mapping.Href
	}

	return client.post(url, mapping, mapping)
}

//IsAccountStoreDirectory checks if a given ApplicationAccountStoreMapping maps an Application to a Directory
func (mapping *ApplicationAccountStoreMapping) IsAccountStoreDirectory() bool {
	return strings.Contains(mapping.AccountStore.Href, "/directories/")
}

//IsAccountStoreGroup checks if a given ApplicationAccountStoreMapping maps an Application to a Group
func (mapping *ApplicationAccountStoreMapping) IsAccountStoreGroup() bool {
	return strings.Contains(mapping.AccountStore.Href, "/groups/")
}

//IsAccountStoreOrganization checks if a given ApplicationAccountStoreMapping maps an Application to an Organization
func (mapping *ApplicationAccountStoreMapping) IsAccountStoreOrganization() bool {
	return strings.Contains(mapping.AccountStore.Href, "/organizations/")
}

//IsAccountStoreDirectory checks if a given OrganizationAccountStoreMapping maps an Application to a Directory
func (mapping *OrganizationAccountStoreMapping) IsAccountStoreDirectory() bool {
	return strings.Contains(mapping.AccountStore.Href, "/directories/")
}

//IsAccountStoreGroup checks if a given OrganizationAccountStoreMapping maps an Application to a Directory
func (mapping *OrganizationAccountStoreMapping) IsAccountStoreGroup() bool {
	return strings.Contains(mapping.AccountStore.Href, "/groups/")
}

//IsAccountStoreOrganization checks if a given ApplicationAccountStoreMapping maps an Application to an Organization
func (mapping *OrganizationAccountStoreMapping) IsAccountStoreOrganization() bool {
	return strings.Contains(mapping.AccountStore.Href, "/organizations/")
}
