package stormpath

import (
	"net/url"

	"github.com/asaskevich/govalidator"
)

//Group represents a Stormpath Group
//
//See: http://docs.stormpath.com/rest/product-guide/#groups
type Group struct {
	accountStoreResource
	Name        string     `json:"name,omitempty" valid:"required,length(1|255)"`
	Description string     `json:"description,omitempty" valid:"length(0|1000)"`
	Status      string     `json:"status,omitempty"`
	Tenant      *Tenant    `json:"tenant,omitempty"`
	Directory   *Directory `json:"directory,omitempty"`
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

func (group *Group) GetGroupMemberships(pageRequest url.Values, filter url.Values) (*GroupMemberships, error) {
	groupMemberships := &GroupMemberships{}

	err := client.get(
		buildAbsoluteURL(group.Href, "accountMemberships", requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
		groupMemberships,
	)

	return groupMemberships, err
}
