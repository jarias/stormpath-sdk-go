package stormpath

//Group represents a Stormpath Group
//
//See: http://docs.stormpath.com/rest/product-guide/#groups
type Group struct {
	accountStoreResource
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status,omitempty"`
	Tenant      *Tenant    `json:"tenant,omitempty"`
	Directory   *Directory `json:"directory,omitempty"`
	CreatedAt   Date       `json:"createdAt,omitempty"`
	ModifiedAt  Date       `json:"modifiedAt,omitempty"`
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

func GetGroup(href string, criteria Criteria) (*Group, error) {
	group := &Group{}

	err := client.get(
		buildAbsoluteURL(href, criteria.ToQueryString()),
		emptyPayload(),
		group,
	)

	return group, err
}

func (group *Group) GetGroupMemberships(criteria Criteria) (*GroupMemberships, error) {
	groupMemberships := &GroupMemberships{}

	err := client.get(
		buildAbsoluteURL(group.Href, "accountMemberships", criteria.ToQueryString()),
		emptyPayload(),
		groupMemberships,
	)

	return groupMemberships, err
}
