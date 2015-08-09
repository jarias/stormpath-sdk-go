package stormpath

import (
	"net/url"

	"github.com/asaskevich/govalidator"
)

//Directory represents a Stormpath directory object
//
//See: http://docs.stormpath.com/rest/product-guide/#directories
type Directory struct {
	accountStoreResource
	Name        string  `json:"name,omitempty" valid:"required,length(1|255)"`
	Description string  `json:"description,omitempty" valid:"length(0|1000)"`
	Status      string  `json:"status,omitempty"`
	Groups      *Groups `json:"groups,omitempty"`
	Tenant      *Tenant `json:"tenant,omitempty"`
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

//Validate validates a directory, returns true if valid and false + error if not
func (dir *Directory) Validate() (bool, error) {
	return govalidator.ValidateStruct(dir)
}

//Refresh refreshes the directory resource by doing a GET to the directory href endpoint
func (dir *Directory) Refresh() error {
	return client.get(dir.Href, emptyPayload(), dir)
}

//Save saves the directory in Stormpath
func (dir *Directory) Save() error {
	ok, err := dir.Validate()
	if !ok && err != nil {
		return err
	}
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
