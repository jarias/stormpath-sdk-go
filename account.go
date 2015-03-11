package stormpath

import "net/url"

//Account represents an Stormpath account object
//
//See: http://docs.stormpath.com/rest/product-guide/#accounts
type Account struct {
	Href                   string `json:"href,omitempty"`
	Username               string `json:"username,omitempty"`
	Email                  string `json:"email"`
	Password               string `json:"password"`
	FullName               string `json:"fullName,omitempty"`
	GivenName              string `json:"givenName"`
	MiddleName             string `json:"middleName,omitempty"`
	Surname                string `json:"surname"`
	Status                 string `json:"status,omitempty"`
	CustomData             *link  `json:"customData,omitempty"`
	Groups                 *link  `json:"groups,omitempty"`
	GroupMemberships       *link  `json:"groupMemberships,omitempty"`
	Directory              *link  `json:"directory,omitempty"`
	Tenant                 *link  `json:"tenant,omitempty"`
	EmailVerificationToken *link  `json:"emailVerificationToken,omitempty"`
}

//Accounts represents a paged result of Account objects
//
//See: http://docs.stormpath.com/rest/product-guide/#accounts-list
type Accounts struct {
	list
	Items []Account `json:"items"`
}

//AccountRef represent a link to an account, this type of resource is return when expand is not especify,
//use only in account authentication
//
//See: http://docs.stormpath.com/rest/product-guide/#application-accounts (Log In (Authenticate) an Account)
type AccountRef struct {
	Account link
}

//AccountPasswordResetToken represents an password reset token for a given account
//
//See: http://docs.stormpath.com/rest/product-guide/#application-accounts (Reset An Account’s Password)
type AccountPasswordResetToken struct {
	Href    string
	Email   string
	Account link
}

type SocialAccount struct {
	Data ProviderData `json:"providerData"`
}

type ProviderData struct {
	ProviderId string `json:"providerId"`
	AccessToken string `json:"accessToken"`
}

//NewAccount is a conviniece constructor for an Account, it accepts all the required fields according to
//the Stormpath API, it returns a pointer to an Account
func NewAccount(email string, password string, givenName string, surname string) *Account {
	return &Account{Email: email, Password: password, GivenName: givenName, Surname: surname}
}

//GetAccount returns the Account from an AccountRef
func (accountRef *AccountRef) GetAccount() (*Account, error) {
	account := &Account{}

	err := client.get(accountRef.Account.Href, emptyPayload(), account)

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
		buildAbsoluteURL(account.GroupMemberships.Href, requestParams(pageRequest, NewEmptyFilter(), url.Values{})),
		emptyPayload(),
		groupMemberships,
	)

	return groupMemberships, err
}

//GetCustomData returns the given account custom data as a map
func (account *Account) GetCustomData() (map[string]interface{}, error) {
	customData := make(map[string]interface{})

	err := client.get(account.CustomData.Href, emptyPayload(), &customData)

	return customData, err
}

//UpdateCustomData sets or updates the given account custom data
func (account *Account) UpdateCustomData(data map[string]interface{}) error {
	return client.post(account.CustomData.Href, data, &data)
}
