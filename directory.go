package stormpath

type Directory struct {
	Href        *string          `json:"href,omitempty"`
	Name        string           `json:"name"`
	Description *string          `json:"description,omitempty"`
	Status      *string          `json:"status,omitempty"`
	Accounts    *Link            `json:"accounts,omitempty"`
	Groups      *Link            `json:"groups,omitempty"`
	Tenant      *Link            `json:"tenant,omitempty"`
	Client      *StormpathClient `json:"-"`
}

func NewDirectory(name string) *Directory {
	return &Directory{Name: name}
}
