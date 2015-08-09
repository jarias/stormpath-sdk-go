package stormpath

import (
	"net/url"

	"github.com/asaskevich/govalidator"
)

//Group represents a Stormpath Group
//
//See: http://docs.stormpath.com/rest/product-guide/#groups
type Group struct {
	resource
	Name        string      `json:"name,omitempty" valid:"required,length(1|255)"`
	Description string      `json:"description,omitempty" valid:"length(0|1000)"`
	Status      string      `json:"status,omitempty"`
	CustomData  *CustomData `json:"customData,omitempty"`
	Accounts    *Accounts   `json:"accounts,omitempty"`
	Tenant      *Tenant     `json:"tenant,omitempty"`
	Directory   *Directory  `json:"directory,omitempty"`
}

//Groups represent a paged result of groups
type Groups struct {
	collectionResource
	Items []Group `json:"items"`
}

//NewGroup creates a new Group with the given name
func NewGroup(name string) *Group {
	return &Group{Name: name}
}

//Validate validates a group, returns true if valid and false + error if not
func (group *Group) Validate() (bool, error) {
	return govalidator.ValidateStruct(group)
}

//Refresh refreshes the group resource by doing a GET to the group href endpoint
func (group *Group) Refresh() error {
	return client.get(group.Href, emptyPayload(), group)
}

func (group *Group) Save() error {
	ok, err := group.Validate()
	if !ok && err != nil {
		return err
	}
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
