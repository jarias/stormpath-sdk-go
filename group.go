package stormpath

import "net/url"

//Group represents a Stormpath Group
//
//See: http://docs.stormpath.com/rest/product-guide/#groups
type Group struct {
	resource
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Status      string      `json:"status,omitempty"`
	CustomData  *CustomData `json:"customData,omitempty"`
	Accounts    *Accounts   `json:"accounts,omitempty"`
	Tenant      *Tenant     `json:"tenant,omitempty"`
	Directory   *Directory  `json:"directory,omitempty"`
}

//Groups represent a paged result of groups
type Groups struct {
	collectionResource
	Items []Groups `json:"items"`
}

//NewGroup creates a new Group with the given name
func NewGroup(name string) *Group {
	return &Group{Name: name}
}

func MakeGroup(href string) *Group {
	return &Group{resource: resource{Href: href}}
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

//UpdateCustomData updates the group custom data and returns that updated custom data as a map[string]interface
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (group *Group) UpdateCustomData(customData map[string]interface{}) (map[string]interface{}, error) {
	customData = cleanCustomData(customData)

	err := client.post(buildAbsoluteURL(group.Href, "customData"), customData, &customData)

	return customData, err
}

//DeleteCustomData deletes all the group custom data
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (group *Group) DeleteCustomData() error {
	return client.delete(buildAbsoluteURL(group.Href, "customData"), emptyPayload())
}

//GetCustomData gets the group custom data map
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (group *Group) GetCustomData() (map[string]interface{}, error) {
	customData := map[string]interface{}{}

	err := client.get(buildAbsoluteURL(group.Href, "customData"), emptyPayload(), &customData)

	return customData, err
}
