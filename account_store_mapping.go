package stormpath

//AccountStoreMapping represents an Stormpath account store mapping
//
//See: http://docs.stormpath.com/rest/product-guide/#account-store-mappings
type AccountStoreMapping struct {
	Href                  string `json:"href,omitempty"`
	ListIndex             *int   `json:"listIndex,omitempty"`
	IsDefaultAccountStore *bool  `json:"isDefaultAccountStore,omitempty"`
	IsDefaultGroupStore   *bool  `json:"isDefaultGroupStore,omitempty"`
	Application           link   `json:"application"`
	AccountStore          link   `json:"accountStore"`
}

//AccountStoreMappings represents a pages result of account store mappings
//
//See: http://docs.stormpath.com/rest/product-guide/#list-account-store-mappings
type AccountStoreMappings struct {
	list
	Items []AccountStoreMapping `json:"items"`
}

//NewAccountStoreMapping creates a new account store mappings
func NewAccountStoreMapping(applicationHref string, accountStoreHref string) *AccountStoreMapping {
	return &AccountStoreMapping{Application: link{Href: applicationHref}, AccountStore: link{Href: accountStoreHref}}
}

//Save saves the given account store mapping
func (mapping *AccountStoreMapping) Save() error {
	url := buildRelativeURL("accountStoreMappings")
	if mapping.Href != "" {
		url = mapping.Href
	}

	return client.post(url, mapping, mapping)
}
