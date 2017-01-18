package stormpath

import "net/url"

type GroupMembershipCriteria struct {
	baseCriteria
}

func MakeGroupMemershipCriteria() GroupMembershipCriteria {
	return GroupMembershipCriteria{baseCriteria{filter: url.Values{}}}
}

func MakeGroupMemershipsCriteria() GroupMembershipCriteria {
	return GroupMembershipCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//Pagination

func (c GroupMembershipCriteria) Limit(limit int) GroupMembershipCriteria {
	c.limit = limit
	return c
}

func (c GroupMembershipCriteria) Offset(offset int) GroupMembershipCriteria {
	c.offset = offset
	return c
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
