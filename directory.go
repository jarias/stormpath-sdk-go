package stormpath

const (
	DirectoryBaseUrl = "https://api.stormpath.com/v1/directories"
)

type Directory struct {
	Href        string `json:"href,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Accounts    *Link  `json:"accounts,omitempty"`
	Groups      *Link  `json:"groups,omitempty"`
	Tenant      *Link  `json:"tenant,omitempty"`
}

func NewDirectory(name string) *Directory {
	return &Directory{Name: name}
}

func (dir *Directory) Save() error {
	return Client.DoWithResult(&StormpathRequest{
		Method:  POST,
		URL:     DirectoryBaseUrl,
		Payload: dir,
	}, dir)
}

func (dir *Directory) Delete() error {
	_, err := Client.Do(&StormpathRequest{
		Method: DELETE,
		URL:    dir.Href,
	})
	return err
}

func (dir *Directory) GetGroups(pageRequest PageRequest, filter Filter) (*Groups, error) {
	groups := &Groups{}

	err := Client.DoWithResult(&StormpathRequest{
		Method:      GET,
		URL:         dir.Groups.Href,
		PageRequest: &pageRequest,
		Filter:      filter,
	}, groups)

	return groups, err
}

func (dir *Directory) GetAccounts(pageRequest PageRequest, filter Filter) (*Accounts, error) {
	accounts := &Accounts{}

	err := Client.DoWithResult(&StormpathRequest{
		Method:      GET,
		URL:         dir.Accounts.Href,
		PageRequest: &pageRequest,
		Filter:      filter,
	}, accounts)

	return accounts, err
}
