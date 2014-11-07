package stormpath

import "net/url"

type Directory struct {
	Href        string `json:"href,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Accounts    *link  `json:"accounts,omitempty"`
	Groups      *link  `json:"groups,omitempty"`
	Tenant      *link  `json:"tenant,omitempty"`
}

type Directories struct {
	list
	Items []Directory `json:"items"`
}

func NewDirectory(name string) *Directory {
	return &Directory{Name: name}
}

func (dir *Directory) Save() error {
	return Client.doWithResult(Client.newRequest(
		"POST",
		dir.Href,
		dir,
	), dir)
}

func (dir *Directory) Delete() error {
	return Client.do(Client.newRequest(
		"DELETE",
		dir.Href,
		emptyPayload(),
	))
}

func (dir *Directory) GetGroups(pageRequest url.Values, filter url.Values) (*Groups, error) {
	groups := &Groups{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		buildAbsoluteURL(dir.Groups.Href, requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
	), groups)

	return groups, err
}

func (dir *Directory) GetAccounts(pageRequest url.Values, filter url.Values) (*Accounts, error) {
	accounts := &Accounts{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		buildAbsoluteURL(dir.Accounts.Href, requestParams(pageRequest, filter, url.Values{})),
		emptyPayload(),
	), accounts)

	return accounts, err
}

func (dir *Directory) CreateGroup(group *Group) error {
	return Client.doWithResult(Client.newRequest(
		"POST",
		dir.Groups.Href,
		group,
	), group)
}
