package stormpath

import "net/url"

//Directory represents a Stormpath directory object
//
//See: http://docs.stormpath.com/rest/product-guide/#directories
type Directory struct {
	Href        string `json:"href,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Accounts    *link  `json:"accounts,omitempty"`
	Groups      *link  `json:"groups,omitempty"`
	Tenant      *link  `json:"tenant,omitempty"`
}

//Directories represnets a paged result of directories
type Directories struct {
	list
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
