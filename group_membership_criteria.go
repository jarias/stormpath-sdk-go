package stormpath

import "net/url"

type GroupMembershipCriteria struct {
	baseCriteria
}

func MakeGroupMemershipCriteria() GroupMembershipCriteria {
	return GroupMembershipCriteria{baseCriteria{filter: url.Values{}}}
}

func MakeGroupMemershipsCriteria() GroupCriteria {
	return GroupCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//Expansion related functions

func (c GroupMembershipCriteria) WithGroup() GroupMembershipCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "group")
	return c
}

func (c GroupMembershipCriteria) WithAccount() GroupMembershipCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "account")
	return c
}
