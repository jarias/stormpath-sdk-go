package stormpath

const (
	GroupMembershipBaseUrl = "https://api.stormpath.com/v1/groupMemberships"
)

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
	return Client.do(Client.newRequest(
		"DELETE",
		groupmembership.Href,
		nil,
	))
}

func (groupmembership *GroupMembership) GetAccount() (*Account, error) {
	account := &Account{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		groupmembership.Account.Href,
		nil,
	), account)

	return account, err
}

func (groupmembership *GroupMembership) GetGroup() (*Group, error) {
	group := &Group{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		groupmembership.Group.Href,
		nil,
	), group)

	return group, err
}
