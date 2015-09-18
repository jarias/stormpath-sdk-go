package stormpath

import (
	"time"
)

type GroupMembership struct {
	resource
	Account    Account   `json:"account"`
	Group      Group     `json:"group"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	ModifiedAt time.Time `json:"modifiedAt,omitempty"`
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
		emptyPayload(),
		account,
	)

	return account, err
}

func (groupmembership *GroupMembership) GetGroup(criteria Criteria) (*Group, error) {
	group := &Group{}

	err := client.get(
		buildAbsoluteURL(groupmembership.Group.Href, criteria.ToQueryString()),
		emptyPayload(),
		group,
	)

	return group, err
}
