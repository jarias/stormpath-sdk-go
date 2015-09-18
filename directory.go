package stormpath

import (
	"time"
)

//Directory represents a Stormpath directory object
//
//See: http://docs.stormpath.com/rest/product-guide/#directories
type Directory struct {
	accountStoreResource
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	Groups      *Groups   `json:"groups,omitempty"`
	Tenant      *Tenant   `json:"tenant,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	ModifiedAt  time.Time `json:"modifiedAt,omitempty"`
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

func GetDirectory(href string, criteria Criteria) (*Directory, error) {
	directory := &Directory{}

	err := client.get(
		buildAbsoluteURL(href, criteria.ToQueryString()),
		emptyPayload(),
		directory,
	)

	return directory, err
}

//GetGroups returns all the groups from a directory
func (dir *Directory) GetGroups(criteria Criteria) (*Groups, error) {
	groups := &Groups{}

	err := client.get(
		buildAbsoluteURL(dir.Groups.Href, criteria.ToQueryString()),
		emptyPayload(),
		groups,
	)

	return groups, err
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
