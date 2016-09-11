package stormpath

//Account represents an Stormpath account object
//
//See: http://docs.stormpath.com/rest/product-guide/#accounts
type Account struct {
	customDataAwareResource
	Username               string            `json:"username,omitempty"`
	Email                  string            `json:"email,omitempty"`
	Password               string            `json:"password,omitempty"`
	FullName               string            `json:"fullName,omitempty"`
	GivenName              string            `json:"givenName,omitempty"`
	MiddleName             string            `json:"middleName,omitempty"`
	Surname                string            `json:"surname,omitempty"`
	Status                 string            `json:"status,omitempty"`
	Groups                 *Groups           `json:"groups,omitempty"`
	GroupMemberships       *GroupMemberships `json:"groupMemberships,omitempty"`
	Directory              *Directory        `json:"directory,omitempty"`
	Tenant                 *Tenant           `json:"tenant,omitempty"`
	EmailVerificationToken *resource         `json:"emailVerificationToken"`
	AccessTokens           *OAuthTokens      `json:"accessTokens,omitempty"`
	RefreshTokens          *OAuthTokens      `json:"refreshTokens,omitempty"`
	ProviderData           *ProviderData     `json:"providerData,omitempty"`
	APIKeys                *APIKeys          `json:"apiKeys,omitempty"`
	Applications           *Applications     `json:"applications,omitempty"`
}

//Accounts represents a paged result of Account objects
//
//See: http://docs.stormpath.com/rest/product-guide/#accounts-collectionResource
type Accounts struct {
	collectionResource
	Items []Account `json:"items,omitempty"`
}

//AccountPasswordResetToken represents an password reset token for a given account
//
//See: http://docs.stormpath.com/rest/product-guide/#application-accounts (Reset An Account’s Password)
type AccountPasswordResetToken struct {
	Href    string
	Email   string
	Account Account
}

type accountRef struct {
	Account *Account `json:"account"`
}

//SocialAccount represents the JSON payload use to create an account for a social backend directory
//(Google, Facebook, Github, etc)
type SocialAccount struct {
	Data ProviderData `json:"providerData"`
}

//ProviderData represents the especific information needed by the social provider (Google, Github, Faceboo, etc)
type ProviderData struct {
	ProviderID  string `json:"providerId"`
	AccessToken string `json:"accessToken,omitempty"`
	Code        string `json:"code,omitempty"`
}

//NewAccount returns a pointer to an Account with the minimum data required
func NewAccount(username, password, email, givenName, surname string) *Account {
	return &Account{Username: username, Password: password, Email: email, GivenName: givenName, Surname: surname}
}

//GetAccount fetches an account by href and criteria
func GetAccount(href string, criteria AccountCriteria) (*Account, error) {
	account := &Account{}

	err := client.get(
		buildAbsoluteURL(href, criteria.toQueryString()),
		account,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (account *Account) Refresh() error {
	return client.get(account.Href, account)
}

//Update updates the given resource, by doing a POST to the resource Href
func (account *Account) Update() error {
	return client.post(account.Href, account, account)
}

//AddToGroup adds the given account to a given group and returns the respective GroupMembership
func (account *Account) AddToGroup(group *Group) (*GroupMembership, error) {
	groupMembership := NewGroupMembership(account.Href, group.Href)

	err := client.post(buildRelativeURL("groupMemberships"), groupMembership, groupMembership)

	if err != nil {
		return nil, err
	}

	return groupMembership, nil
}

//RemoveFromGroup removes the given account from the given group by searching the account groupmemberships,
//and deleting the corresponding one
func (account *Account) RemoveFromGroup(group *Group) error {
	groupMemberships, err := account.GetGroupMemberships(
		MakeGroupMemershipCriteria().Offset(0).Limit(25),
	)

	if err != nil {
		return err
	}

	for i := 1; len(groupMemberships.Items) > 0; i++ {
		for _, gm := range groupMemberships.Items {
			if gm.Group.Href == group.Href {
				return gm.Delete()
			}
		}
		groupMemberships, err = account.GetGroupMemberships(
			MakeGroupMemershipCriteria().Offset(i * 25).Limit(25),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

//GetGroupMemberships returns a paged result of the group memeberships of the given account
func (account *Account) GetGroupMemberships(criteria GroupMembershipCriteria) (*GroupMemberships, error) {
	groupMemberships := &GroupMemberships{}

	err := client.get(
		buildAbsoluteURL(
			account.GroupMemberships.Href,
			criteria.toQueryString(),
		),
		groupMemberships,
	)

	if err != nil {
		return nil, err
	}

	return groupMemberships, nil
}

//VerifyEmailToken verifies an email verification token associated with an account
//
//See: http://docs.stormpath.com/rest/product-guide/#account-verify-email
func VerifyEmailToken(token string) (*Account, error) {
	account := &Account{}
	err := client.post(buildRelativeURL("accounts/emailVerificationTokens", token), emptyPayload(), account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

//GetRefreshTokens returns the account's refreshToken collection
func (account *Account) GetRefreshTokens(criteria OAuthTokenCriteria) (*OAuthTokens, error) {
	refreshTokens := &OAuthTokens{}

	err := client.get(
		buildAbsoluteURL(account.RefreshTokens.Href, criteria.toQueryString()),
		refreshTokens,
	)

	if err != nil {
		return nil, err
	}

	return refreshTokens, nil
}

//GetAccessTokens returns the acounts's accessToken collection
func (account *Account) GetAccessTokens(criteria OAuthTokenCriteria) (*OAuthTokens, error) {
	accessTokens := &OAuthTokens{}

	err := client.get(
		buildAbsoluteURL(account.AccessTokens.Href, criteria.toQueryString()),
		accessTokens,
	)

	if err != nil {
		return nil, err
	}

	return accessTokens, nil
}

//CreateAPIKey creates a new API key pair for the given account, it returns a pointer to the APIKey pair.
func (account *Account) CreateAPIKey() (*APIKey, error) {
	apiKey := &APIKey{}

	err := client.post(account.APIKeys.Href, emptyPayload(), apiKey)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}
