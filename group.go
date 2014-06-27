package stormpath

type Group struct {
	Href        string `json:"href,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	CustomData  *Link  `json:"customData,omitempty"`
	Accounts    *Link  `json:"accounts,omitempty"`
	Tenant      *Link  `json:"tenant,omitempty"`
	Directory   *Link  `json:"directory,omitempty"`
}

func NewGroup(name string) *Group {
	return &Group{Name: name}
}

func (group *Group) Save() error {
	return Client.DoWithResult(&StormpathRequest{
		Method:  POST,
		URL:     group.Href,
		Payload: group,
	}, group)
}

func (group *Group) Delete() error {
	_, err := Client.Do(&StormpathRequest{
		Method: DELETE,
		URL:    group.Href,
	})
	return err
}

func (group *Group) GetAccounts(pageRequest PageRequest, filter Filter) (*Accounts, error) {
	accounts := &Accounts{}

	err := Client.DoWithResult(&StormpathRequest{
		Method:      GET,
		URL:         group.Accounts.Href,
		PageRequest: &pageRequest,
		Filter:      filter,
	}, accounts)

	return accounts, err
}

func (group *Group) GetGroupMemberships(pageRequest PageRequest, filter Filter) (*GroupMemberships, error) {
	groupMemberships := &GroupMemberships{}

	err := Client.DoWithResult(&StormpathRequest{
		Method:      GET,
		URL:         group.Href + "/accountMemberships",
		PageRequest: &pageRequest,
		Filter:      filter,
	}, groupMemberships)

	return groupMemberships, err
}
