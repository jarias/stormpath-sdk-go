package stormpath

//Organization represnts the Stormpath organization resource, use for multitenancy
type Organization struct {
	accountStoreResource
	Name                       string                            `json:"name,omitempty"`
	Description                string                            `json:"description,omitempty"`
	Status                     string                            `json:"status,omitempty"`
	NameKey                    string                            `json:"nameKey,omitempty"`
	Groups                     *Groups                           `json:"groups,omitempty"`
	Tenant                     *Tenant                           `json:"tenant,omitempty"`
	AccountStoreMappings       *OrganizationAccountStoreMappings `json:"accountStoreMappings,omitempty"`
	DefaultAccountStoreMapping *OrganizationAccountStoreMapping  `json:"defaultAccountStoreMapping,omitempty"`
	DefaultGroupStoreMapping   *OrganizationAccountStoreMapping  `json:"defaultGroupStoreMapping,omitempty"`
}

//Organizations represents a paged result or applications
type Organizations struct {
	collectionResource
	Items []Organization `json:"items,omitempty"`
}

//NewOrganization creates a new organization
func NewOrganization(name string, nameKey string) *Organization {
	return &Organization{Name: name, NameKey: nameKey}
}

//CreateOrganization creates new organization for the given tenant
func (tenant *Tenant) CreateOrganization(org *Organization) error {
	return client.post(buildRelativeURL("organizations"), org, org)
}

//GetOrganization loads an organization by href and criteria
func GetOrganization(href string, criteria OrganizationCriteria) (*Organization, error) {
	organization := &Organization{}

	err := client.get(
		buildAbsoluteURL(href, criteria.toQueryString()),
		organization,
	)

	return organization, err
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (org *Organization) Refresh() error {
	return client.get(org.Href, org)
}

//Update updates the given resource, by doing a POST to the resource Href
func (org *Organization) Update() error {
	return client.post(org.Href, org, org)
}

//GetAccountStoreMappings returns all the applications account store mappings
func (org *Organization) GetAccountStoreMappings(criteria OrganizationAccountStoreMappingCriteria) (*OrganizationAccountStoreMappings, error) {
	accountStoreMappings := &OrganizationAccountStoreMappings{}

	err := client.get(
		buildAbsoluteURL(org.AccountStoreMappings.Href, criteria.toQueryString()),
		accountStoreMappings,
	)

	if err != nil {
		return nil, err
	}

	return accountStoreMappings, nil
}

func (org *Organization) GetDefaultAccountStoreMapping(criteria OrganizationAccountStoreMappingCriteria) (*OrganizationAccountStoreMapping, error) {
	err := client.get(
		buildAbsoluteURL(org.DefaultAccountStoreMapping.Href, criteria.toQueryString()),
		org.DefaultAccountStoreMapping,
	)

	if err != nil {
		return nil, err
	}

	return org.DefaultAccountStoreMapping, nil
}

//RegisterAccount registers a new account into the organization
func (org *Organization) RegisterAccount(account *Account) error {
	err := client.post(org.Accounts.Href, account, account)
	if err == nil {
		//Password should be cleanup so we don't keep an unhash password in memory
		account.Password = ""
	}
	return err
}

//RegisterSocialAccount registers a new account into the organization using an external provider Google, Facebook
func (org *Organization) RegisterSocialAccount(socialAccount *SocialAccount) (*Account, error) {
	account := &Account{}

	err := client.post(org.Accounts.Href, socialAccount, account)

	if err != nil {
		return nil, err
	}

	return account, nil
}
