package stormpath

import "net/url"

//Account represents an Stormpath account object
//
//See: http://docs.stormpath.com/rest/product-guide/#accounts
type Account struct {
	resource
	Username               string            `json:"username,omitempty"`
	Email                  string            `json:"email"`
	Password               string            `json:"password"`
	FullName               string            `json:"fullName,omitempty"`
	GivenName              string            `json:"givenName"`
	MiddleName             string            `json:"middleName,omitempty"`
	Surname                string            `json:"surname"`
	Status                 string            `json:"status,omitempty"`
	CustomData             *CustomData       `json:"customData,omitempty"`
	Groups                 *Groups           `json:"groups,omitempty"`
	GroupMemberships       *GroupMemberships `json:"groupMemberships,omitempty"`
	Directory              *Directory        `json:"directory,omitempty"`
	Tenant                 *Tenant           `json:"tenant,omitempty"`
	EmailVerificationToken *resource         `json:"emailVerificationToken,omitempty"`
}

//Accounts represents a paged result of Account objects
//
//See: http://docs.stormpath.com/rest/product-guide/#accounts-collectionResource
type Accounts struct {
	collectionResource
	Items []Account `json:"items"`
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
	Account resource `json:"account"`
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
}

//NewAccount is a conviniece constructor for an Account, it accepts all the required fields according to
//the Stormpath API, it returns a pointer to an Account
func NewAccount(email string, password string, givenName string, surname string) *Account {
	return &Account{Email: email, Password: password, GivenName: givenName, Surname: surname}
}

//MakeAccount creates an account resource from an href
func MakeAccount(href string) *Account {
	return &Account{resource: resource{Href: href}}
}

//VerifyEmailToken verifies an email verification token associated with an account
//
//See: http://docs.stormpath.com/rest/product-guide/#account-verify-email
func VerifyEmailToken(token string) (*Account, error) {
	account := &Account{}
	err := client.post(buildAbsoluteURL(BaseURL, "accounts/emailVerificationTokens", token), emptyPayload(), account)

	return account, err
}

//Load refreses the account from the account href
func (account *Account) Load() (*Account, error) {
	err := client.get(account.Href, emptyPayload(), account)

	return account, err
}

//Save updates the given account, by doing a POST to the account Href, if the account is a new account
//it should be created via Application.RegisterAccount
func (account *Account) Save() error {
	return client.post(account.Href, account, account)
}

//Delete deletes the given account, it wont modify the calling account
func (account *Account) Delete() error {
	return client.delete(account.Href, emptyPayload())
}

//AddToGroup adds the given account to a given group and returns the respective GroupMembership
func (account *Account) AddToGroup(group *Group) (*GroupMembership, error) {
	groupMembership := NewGroupMembership(account.Href, group.Href)

	err := client.post(buildRelativeURL("groupMemberships"), groupMembership, groupMembership)

	return groupMembership, err
}

//RemoveFromGroup removes the given account from the given group by searching the account groupmemberships,
//and deleting the corresponding one
func (account *Account) RemoveFromGroup(group *Group) error {
	groupMemberships, err := account.GetGroupMemberships(NewDefaultPageRequest())

	if err != nil {
		return err
	}

	for i := 1; len(groupMemberships.Items) > 0; i++ {
		for _, gm := range groupMemberships.Items {
			if gm.Group.Href == group.Href {
				return gm.Delete()
			}
		}
		groupMemberships, err = account.GetGroupMemberships(NewPageRequest(25, i*25))
		if err != nil {
			return err
		}
	}

	return nil
}

//GetGroupMemberships returns a paged result of the group memeberships of the given account
func (account *Account) GetGroupMemberships(pageRequest url.Values) (*GroupMemberships, error) {
	groupMemberships := &GroupMemberships{}

	err := client.get(
		buildAbsoluteURL(
			account.GroupMemberships.Href,
			requestParams(pageRequest, NewEmptyFilter(), url.Values{}),
		),
		emptyPayload(),
		groupMemberships,
	)

	return groupMemberships, err
}

//GetCustomData returns the given account custom data as a map
func (account *Account) GetCustomData() (map[string]interface{}, error) {
	customData := make(map[string]interface{})

	err := client.get(buildAbsoluteURL(account.Href, "customData"), emptyPayload(), &customData)

	return customData, err
}

//UpdateCustomData sets or updates the given account custom data
func (account *Account) UpdateCustomData(data map[string]interface{}) error {
	// delete illegal keys from data
	// http://docs.stormpath.com/rest/product-guide/#custom-data
	keys := []string{
		"href", "createdAt", "modifiedAt", "meta",
		"spMeta", "spmeta", "ionmeta", "ionMeta",
	}

	for i := range keys {
		delete(data, keys[i])
	}

	return client.post(buildAbsoluteURL(account.Href, "customData"), data, &data)
}
