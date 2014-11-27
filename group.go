package stormpath

import "net/url"

//Group represents a Stormpath Group
//
//See: http://docs.stormpath.com/rest/product-guide/#groups
type Group struct {
	Href        string `json:"href,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	CustomData  *link  `json:"customData,omitempty"`
	Accounts    *link  `json:"accounts,omitempty"`
	Tenant      *link  `json:"tenant,omitempty"`
	Directory   *link  `json:"directory,omitempty"`
}

//Groups represent a paged result of groups
type Groups struct {
	list
	Items []Groups `json:"items"`
}

//NewGroup creates a new Group with the given name
func NewGroup(name string) *Group {
	return &Group{Name: name}
}

func (group *Group) Save() error {
	return client.post(group.Href, group, group)
}

func (group *Group) Delete() error {
	return client.delete(group.Href, emptyPayload())
}

func (group *Group) GetAccounts(pageRequest url.Values, filter url.Values) (*Accounts, error) {
	accounts := &Accounts{}

	err := client.get(
		buildAbsoluteURL(group.Accounts.Href, requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
		accounts,
	)

	return accounts, err
}

func (group *Group) GetGroupMemberships(pageRequest url.Values, filter url.Values) (*GroupMemberships, error) {
	groupMemberships := &GroupMemberships{}

	err := client.get(
		buildAbsoluteURL(group.Href, "accountMemberships", requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
		groupMemberships,
	)

	return groupMemberships, err
}
