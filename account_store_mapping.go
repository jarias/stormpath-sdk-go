package stormpath

type AccountStoreMapping struct {
	Href                  string `json:"href,omitempty"`
	ListIndex             *int   `json:"listIndex,omitempty"`
	IsDefaultAccountStore *bool  `json:"isDefaultAccountStore,omitempty"`
	IsDefaultGroupStore   *bool  `json:"isDefaultGroupStore,omitempty"`
	Application           link   `json:"application"`
	AccountStore          link   `json:"accountStore"`
}

type AccountStoreMappings struct {
	list
	Items []AccountStoreMapping `json:"items"`
}

func NewAccountStoreMapping(applicationHref string, accountStoreHref string) *AccountStoreMapping {
	return &AccountStoreMapping{Application: link{Href: applicationHref}, AccountStore: link{Href: accountStoreHref}}
}

func (mapping *AccountStoreMapping) Save() error {
	url := buildURL("accountStoreMappings")
	if mapping.Href != "" {
		url = mapping.Href
	}

	return Client.doWithResult(Client.newRequest(
		"POST",
		url,
		newPayloadReader(mapping),
	), mapping)
}
