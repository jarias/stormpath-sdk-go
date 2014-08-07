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
		group,
	), group)
}

func (group *Group) Delete() error {
	return Client.do(Client.newRequest(
		"DELETE",
		group.Href,
		emptyPayload(),
	))
}

func (group *Group) GetAccounts(pageRequest url.Values, filter url.Values) (*Accounts, error) {
	accounts := &Accounts{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		buildAbsoluteURL(group.Accounts.Href, requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
	), accounts)

	return accounts, err
}

func (group *Group) GetGroupMemberships(pageRequest url.Values, filter url.Values) (*GroupMemberships, error) {
	groupMemberships := &GroupMemberships{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		buildAbsoluteURL(group.Href, "accountMemberships", requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
	), groupMemberships)

	return groupMemberships, err
}
