package stormpath

type GroupMembership struct {
	resource
	Account *Account `json:"account"`
	Group   *Group   `json:"group"`
}

type GroupMemberships struct {
	collectionResource
	Items []GroupMembership `json:"items,omitempty"`
}

func NewGroupMembership(accountHref string, groupHref string) *GroupMembership {
	account := Account{}
	account.Href = accountHref
	group := Group{}
	group.Href = groupHref
	return &GroupMembership{
		Account: &account,
		Group:   &group,
	}
}

func (groupmembership *GroupMembership) GetAccount(criteria AccountCriteria) (*Account, error) {
	err := client.get(
		buildAbsoluteURL(groupmembership.Account.Href, criteria.toQueryString()),
		groupmembership.Account,
	)

	if err != nil {
		return nil, err
	}

	return groupmembership.Account, nil
}

func (groupmembership *GroupMembership) GetGroup(criteria GroupCriteria) (*Group, error) {
	err := client.get(
		buildAbsoluteURL(groupmembership.Group.Href, criteria.toQueryString()),
		groupmembership.Group,
	)

	if err != nil {
		return nil, err
	}

	return groupmembership.Group, nil
}
