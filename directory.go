package stormpath

import "net/url"

//Directory represents a Stormpath directory object
//
//See: http://docs.stormpath.com/rest/product-guide/#directories
type Directory struct {
	resource
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Status      string      `json:"status,omitempty"`
	CustomData  *CustomData `json:"customData,omitempty"`
	Accounts    *Accounts   `json:"accounts,omitempty"`
	Groups      *Groups     `json:"groups,omitempty"`
	Tenant      *Tenant     `json:"tenant,omitempty"`
}

//Directories represnets a paged result of directories
type Directories struct {
	collectionResource
	Items []Directory `json:"items"`
}

//NewDirectory creates a new directory with the given name
func NewDirectory(name string) *Directory {
	return &Directory{Name: name}
}

//Save saves the directory in Stormpath
func (dir *Directory) Save() error {
	return client.post(dir.Href, dir, dir)
}

//Delete deletes the directory
func (dir *Directory) Delete() error {
	return client.delete(dir.Href, emptyPayload())
}

//GetGroups returns all the groups from a directory
func (dir *Directory) GetGroups(pageRequest url.Values, filter url.Values) (*Groups, error) {
	groups := &Groups{}

	err := client.get(
		buildAbsoluteURL(dir.Groups.Href, requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
		groups,
	)

	return groups, err
}

//GetAccounts returns all the accounts from the directory
func (dir *Directory) GetAccounts(pageRequest url.Values, filter url.Values) (*Accounts, error) {
	accounts := &Accounts{}

	err := client.get(
		buildAbsoluteURL(dir.Accounts.Href, requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
		accounts,
	)

	return accounts, err
}

//CreateGroup creates a new group in the directory
func (dir *Directory) CreateGroup(group *Group) error {
	return client.post(dir.Groups.Href, group, group)
}

//RegisterAccount registers a new account into the directory
//
//See: http://docs.stormpath.com/rest/product-guide/#directory-accounts
func (dir *Directory) RegisterAccount(account *Account) error {
	return client.post(dir.Accounts.Href, account, account)
}

//RegisterSocialAccount registers a new account into the application using an external provider Google, Facebook
//
//See: http://docs.stormpath.com/rest/product-guide/#accessing-accounts-with-google-authorization-codes-or-an-access-tokens
func (dir *Directory) RegisterSocialAccount(socialAccount *SocialAccount) (*Account, error) {
	account := Account{}

	err := client.post(dir.Accounts.Href, socialAccount, &account)

	return &account, err
}

//UpdateCustomData updates the directory custom data and returns that updated custom data as a map[string]interface
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (dir *Directory) UpdateCustomData(customData map[string]interface{}) (map[string]interface{}, error) {
	updatedCustomData := map[string]interface{}{}

	err := client.post(buildAbsoluteURL(dir.Href, "customData"), customData, &updatedCustomData)

	return updatedCustomData, err
}

//DeleteCustomData deletes all the directory custom data
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (dir *Directory) DeleteCustomData() error {
	return client.delete(buildAbsoluteURL(dir.Href, "customData"), emptyPayload())
}

//GetCustomData gets the directory custom data map
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (dir *Directory) GetCustomData() (map[string]interface{}, error) {
	customData := map[string]interface{}{}

	err := client.get(buildAbsoluteURL(dir.Href, "customData"), emptyPayload(), &customData)

	return customData, err
}
