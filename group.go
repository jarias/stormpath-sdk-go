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

//GetGroup loads a group by href and criteria
func GetGroup(href string, criteria Criteria) (*Group, error) {
	group := &Group{}

	err := client.get(
		buildAbsoluteURL(href, criteria.ToQueryString()),
		emptyPayload(),
		group,
	)

	if err != nil {
		return nil, err
	}

	return group, nil
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (group *Group) Refresh() error {
	return client.get(group.Href, emptyPayload(), group)
}

//Update updates the given resource, by doing a POST to the resource Href
func (group *Group) Update() error {
	return client.post(group.Href, group, group)
}

//GetGroupMemberships loads the given group memeberships
func (group *Group) GetGroupMemberships(criteria Criteria) (*GroupMemberships, error) {
	groupMemberships := &GroupMemberships{}

	err := client.get(
		buildAbsoluteURL(group.Href, "accountMemberships", criteria.ToQueryString()),
		emptyPayload(),
		groupMemberships,
	)

	if err != nil {
		return nil, err
	}

	return groupMemberships, nil
}
