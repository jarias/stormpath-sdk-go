package stormpath

//Directory represents a Stormpath directory object
//
//See: http://docs.stormpath.com/rest/product-guide/#directories
type Directory struct {
	accountStoreResource
	Name                  string                 `json:"name,omitempty"`
	Description           string                 `json:"description,omitempty"`
	Status                string                 `json:"status,omitempty"`
	Groups                *Groups                `json:"groups,omitempty"`
	Tenant                *Tenant                `json:"tenant,omitempty"`
	AccountCreationPolicy *AccountCreationPolicy `json:"accountCreationPolicy,omitempty"`
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

//GetDirectory loads a directory by href and criteria
func GetDirectory(href string, criteria Criteria) (*Directory, error) {
	directory := &Directory{}

	err := client.get(
		buildAbsoluteURL(href, criteria.ToQueryString()),
		emptyPayload(),
		directory,
	)

	if err != nil {
		return nil, err
	}

	return directory, nil
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (dir *Directory) Refresh() error {
	return client.get(dir.Href, emptyPayload(), dir)
}

//Update updates the given resource, by doing a POST to the resource Href
func (dir *Directory) Update() error {
	return client.post(dir.Href, dir, dir)
}

//GetAccountCreationPolicy loads the directory account creation policy
func (dir *Directory) GetAccountCreationPolicy() (*AccountCreationPolicy, error) {
	err := client.get(buildAbsoluteURL(dir.AccountCreationPolicy.Href), emptyPayload(), dir.AccountCreationPolicy)

	if err != nil {
		return nil, err
	}

	return dir.AccountCreationPolicy, nil
}

//GetGroups returns all the groups from a directory
func (dir *Directory) GetGroups(criteria Criteria) (*Groups, error) {
	err := client.get(
		buildAbsoluteURL(dir.Groups.Href, criteria.ToQueryString()),
		emptyPayload(),
		dir.Groups,
	)

	if err != nil {
		return nil, err
	}

	return dir.Groups, nil
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
	account := &Account{}

	err := client.post(dir.Accounts.Href, socialAccount, account)

	if err != nil {
		return nil, err
	}

	return account, nil
}
