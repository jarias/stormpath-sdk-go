package stormpath

type Group struct {
	Href        string `json:"href,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	CustomData  *Link  `json:"customData,omitempty"`
	Accounts    *Link  `json:"accounts,omitempty"`
	Tenant      *Link  `json:"tenant,omitempty"`
	Directory   *Link  `json:"directory,omitempty"`
}

func NewGroup(name string) *Group {
	return &Group{Name: name}
}
