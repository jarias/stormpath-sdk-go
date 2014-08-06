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
		newPayloadReader(dir),
	), dir)
}

func (dir *Directory) Delete() error {
	return Client.do(Client.newRequest(
		"DELETE",
		dir.Href,
		nil,
	))
}

func (dir *Directory) GetGroups(pageRequest PageRequest, filter Filter) (*Groups, error) {
	groups := &Groups{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		dir.Groups.Href+requestParams(&pageRequest, filter, url.Values{}),
		nil,
	), groups)

	return groups, err
}

func (dir *Directory) GetAccounts(pageRequest PageRequest, filter Filter) (*Accounts, error) {
	accounts := &Accounts{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		dir.Accounts.Href+requestParams(&pageRequest, filter, url.Values{}),
		nil,
	), accounts)

	return accounts, err
}

func (dir *Directory) CreateGroup(group *Group) error {
	return Client.doWithResult(Client.newRequest(
		"POST",
		dir.Groups.Href,
		newPayloadReader(group),
	), group)
}
