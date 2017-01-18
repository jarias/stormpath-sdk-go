package stormpath

const (
	Facebook = "facebook"
	Google   = "google"
	GitHub   = "github"
	LinkedIn = "linkedin"
)

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
	Provider              *Provider              `json:"provider,omitempty"`
	AccountCreationPolicy *AccountCreationPolicy `json:"accountCreationPolicy,omitempty"`
	PasswordPolicy        *PasswordPolicy        `json:"passwordPolicy,omitempty"`
}

//Directories represnets a paged result of directories
type Directories struct {
	collectionResource
	Items []Directory `json:"items,omitempty"`
}

//Provider represents the directory provider (cloud, google, github, facebook or linkedin)
type Provider struct {
	resource
	OAuthProvider
	ProviderID string `json:"providerId,omitempty"`
}

//OAuthProvider represents a generic OAuth2 provider for all the social type directories
type OAuthProvider struct {
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	RedirectURI  string `json:"redirectUri,omitempty"`
}

//NewDirectory creates a new directory with the given name
func NewDirectory(name string) *Directory {
	return &Directory{Name: name}
}

func newSocialDirectory(name string, clientID string, clientSecret string, redirectURI string, provider string) *Directory {
	directory := NewDirectory(name)
	directory.Provider = &Provider{
		ProviderID: provider,
		OAuthProvider: OAuthProvider{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURI:  redirectURI,
		},
	}

	return directory
}

//NewFacebookDirectory creates a new directory with a Facebook backed provider
func NewFacebookDirectory(name string, clientID string, clientSecret string) *Directory {
	return newSocialDirectory(name, clientID, clientSecret, "", Facebook)
}

//NewGithubDirectory creates a new directory with a GitHub backed provider
func NewGithubDirectory(name string, clientID string, clientSecret string) *Directory {
	return newSocialDirectory(name, clientID, clientSecret, "", GitHub)
}

//NewGoogleDirectory creates a new directory with a Google backed provider
func NewGoogleDirectory(name string, clientID string, clientSecret string, redirectURI string) *Directory {
	return newSocialDirectory(name, clientID, clientSecret, redirectURI, Google)
}

//NewLinkedInDirectory creates a new directory with a LinkedIn backend provider
func NewLinkedInDirectory(name string, clientID string, clientSecret string, redirectURI string) *Directory {
	return newSocialDirectory(name, clientID, clientSecret, redirectURI, LinkedIn)
}

//CreateDirectory creates a new directory for the given tenant
//
//See: http://docs.stormpath.com/rest/product-guide/#tenant-directories
func CreateDirectory(dir *Directory) error {
	return client.post(buildRelativeURL("directories"), dir, dir)
}

//GetDirectory loads a directory by href and criteria
func GetDirectory(href string, criteria DirectoryCriteria) (*Directory, error) {
	directory := &Directory{}

	err := client.get(
		buildAbsoluteURL(href, criteria.toQueryString()),
		directory,
	)

	if err != nil {
		return nil, err
	}

	return directory, nil
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (dir *Directory) Refresh() error {
	return client.get(dir.Href, dir)
}

//Update updates the given resource, by doing a POST to the resource Href
func (dir *Directory) Update() error {
	return client.post(dir.Href, dir, dir)
}

//GetAccountCreationPolicy loads the directory account creation policy
func (dir *Directory) GetAccountCreationPolicy() (*AccountCreationPolicy, error) {
	err := client.get(buildAbsoluteURL(dir.AccountCreationPolicy.Href), dir.AccountCreationPolicy)

	if err != nil {
		return nil, err
	}

	return dir.AccountCreationPolicy, nil
}

//GetGroups returns all the groups from a directory
func (dir *Directory) GetGroups(criteria GroupCriteria) (*Groups, error) {
	err := client.get(
		buildAbsoluteURL(dir.Groups.Href, criteria.toQueryString()),
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
