package stormpath

type GroupMembership struct {
	resource
	Account Account `json:"account"`
	Group   Group   `json:"group"`
}

type GroupMemberships struct {
	collectionResource
	Items []GroupMembership `json:"items"`
}

func NewGroupMembership(accountHref string, groupHref string) *GroupMembership {
	account := Account{}
	account.Href = accountHref
	group := Group{}
	group.Href = groupHref
	return &GroupMembership{
		Account: account,
		Group:   group,
	}
}

func (groupmembership *GroupMembership) GetAccount(criteria Criteria) (*Account, error) {
	account := &Account{}

	err := client.get(
		buildAbsoluteURL(groupmembership.Account.Href, criteria.ToQueryString()),
		account,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (groupmembership *GroupMembership) GetGroup(criteria Criteria) (*Group, error) {
	group := &Group{}

	err := client.get(
		buildAbsoluteURL(groupmembership.Group.Href, criteria.ToQueryString()),
		group,
	)

	if err != nil {
		return nil, err
	}

	return group, nil
}
