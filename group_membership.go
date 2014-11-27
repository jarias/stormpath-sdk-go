package stormpath

type GroupMembership struct {
	Href    string `json:"href,omitempty"`
	Account link   `json:"account"`
	Group   link   `json:"group"`
}

type GroupMemberships struct {
	list
	Items []GroupMembership `json:"items"`
}

func NewGroupMembership(accountHref string, groupHref string) *GroupMembership {
	return &GroupMembership{Account: link{accountHref}, Group: link{groupHref}}
}

func (groupmembership *GroupMembership) Delete() error {
	return client.delete(groupmembership.Href, emptyPayload())
}

func (groupmembership *GroupMembership) GetAccount() (*Account, error) {
	account := &Account{}

	err := client.get(groupmembership.Account.Href, emptyPayload(), account)

	return account, err
}

func (groupmembership *GroupMembership) GetGroup() (*Group, error) {
	group := &Group{}

	err := client.get(groupmembership.Group.Href, emptyPayload(), group)

	return group, err
}
