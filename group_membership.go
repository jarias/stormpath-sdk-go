package stormpath

const (
	GroupMembershipBaseUrl = "https://api.stormpath.com/v1/groupMemberships"
)

type GroupMembership struct {
	Href    string `json:"href,omitempty"`
	Account Link   `json:"account"`
	Group   Link   `json:"group"`
}

func NewGroupMembership(accountHref string, groupHref string) *GroupMembership {
	return &GroupMembership{Account: Link{accountHref}, Group: Link{groupHref}}
}

func (groupmembership *GroupMembership) Delete() error {
	return Client.Do(&StormpathRequest{
		Method: DELETE,
		URL:    groupmembership.Href,
	})
}

func (groupmembership *GroupMembership) GetAccount() (*Account, error) {
	account := &Account{}

	err := Client.DoWithResult(&StormpathRequest{
		Method: GET,
		URL:    groupmembership.Account.Href,
	}, account)

	return account, err
}

func (groupmembership *GroupMembership) GetGroup() (*Group, error) {
	group := &Group{}

	err := Client.DoWithResult(&StormpathRequest{
		Method: GET,
		URL:    groupmembership.Group.Href,
	}, group)

	return group, err
}
