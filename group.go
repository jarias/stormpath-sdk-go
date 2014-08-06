package stormpath

import "net/url"

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

type Groups struct {
	list
	Items []Groups `json:"items"`
}

func NewGroup(name string) *Group {
	return &Group{Name: name}
}

func (group *Group) Save() error {
	return Client.doWithResult(Client.newRequest(
		"POST",
		group.Href,
		newPayloadReader(group),
	), group)
}

func (group *Group) Delete() error {
	return Client.do(Client.newRequest(
		"DELETE",
		group.Href,
		nil,
	))
}

func (group *Group) GetAccounts(pageRequest PageRequest, filter Filter) (*Accounts, error) {
	accounts := &Accounts{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		group.Accounts.Href+requestParams(&pageRequest, filter, url.Values{}),
		nil,
	), accounts)

	return accounts, err
}

func (group *Group) GetGroupMemberships(pageRequest PageRequest, filter Filter) (*GroupMemberships, error) {
	groupMemberships := &GroupMemberships{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		group.Href+"/accountMemberships"+requestParams(&pageRequest, filter, url.Values{}),
		nil,
	), groupMemberships)

	return groupMemberships, err
}
