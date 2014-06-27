package stormpath

type List struct {
	Href   string `json:"href"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type Applications struct {
	List
	Items []Application `json:"items"`
}

type Directories struct {
	List
	Items []Directory `json:"items"`
}

type Groups struct {
	List
	Items []Groups `json:"items"`
}

type AccountStoreMappings struct {
	List
	Items []AccountStoreMapping `json:"items"`
}

type GroupMemberships struct {
	List
	Items []GroupMembership `json:"items"`
}
