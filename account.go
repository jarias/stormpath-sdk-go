package stormpath

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
	CustomData             *Link  `json:"customData,omitempty"`
	Groups                 *Link  `json:"groups,omitempty"`
	GroupMemberships       *Link  `json:"groupMemberships,omitempty"`
	Directory              *Link  `json:"directory,omitempty"`
	Tenant                 *Link  `json:"tenant,omitempty"`
	EmailVerificationToken *Link  `json:"emailVerificationToken,omitempty"`
}

type Accounts struct {
	Href   string     `json:"href"`
	Offset int        `json:"offset"`
	Limit  int        `json:"limit"`
	Items  []*Account `json:"items"`
}

type AccountRef struct {
	Account Link
}

type AccountPasswordResetToken struct {
	Href    string
	Email   string
	Account Link
}

func NewAccount(email string, password string, givenName string, surname string) *Account {
	return &Account{Email: email, Password: password, GivenName: givenName, Surname: surname}
}

func (account *Account) Save() error {
	return Client.DoWithResult(&StormpathRequest{
		Method:  POST,
		URL:     account.Href,
		Payload: account,
	}, account)
}

func (account *Account) Delete() error {
	return Client.Do(&StormpathRequest{
		Method: DELETE,
		URL:    account.Href,
	})
}

func (account *Account) AddToGroup(group *Group) (*GroupMembership, error) {
	groupMembership := NewGroupMembership(account.Href, group.Href)

	err := Client.DoWithResult(&StormpathRequest{
		Method:  POST,
		URL:     GroupMembershipBaseUrl,
		Payload: groupMembership,
	}, groupMembership)

	return groupMembership, err
}

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

func (account *Account) GetGroupMemberships(pageRequest PageRequest) (*GroupMemberships, error) {
	groupMemberships := &GroupMemberships{}

	err := Client.DoWithResult(&StormpathRequest{
		Method:      GET,
		URL:         account.GroupMemberships.Href,
		PageRequest: pageRequest,
	}, groupMemberships)

	return groupMemberships, err
}

func (account *Account) GetCustomData() (map[string]string, error) {
	customData := make(map[string]string)

	err := Client.DoWithResult(&StormpathRequest{
		Method: GET,
		URL:    account.CustomData.Href,
	}, customData)

	return customData, err
}

func (account *Account) SetCustomData(data map[string]string) error {
	return Client.DoWithResult(&StormpathRequest{
		Method:  POST,
		URL:     account.CustomData.Href,
		Payload: data,
	}, data)
}
