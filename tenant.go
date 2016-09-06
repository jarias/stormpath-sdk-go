package stormpath

//Tenant
//
//When you sign up for Stormpath, a private data space is created for you. This space is represented as a Tenant resource in the Stormpath REST API. Your Tenant resource can be thought of as your global starting point. You can access everything in your space by accessing your Tenant resource first and then interacting with its other linked resources (Applications, Directories, etc).
type Tenant struct {
	customDataAwareResource
	Name          string         `json:"name,omitempty"`
	Key           string         `json:"key,omitempty"`
	Accounts      *Accounts      `json:"accounts,omitempty"`
	Applications  *Applications  `json:"applications,omitempty"`
	Directories   *Directories   `json:"directories,omitempty"`
	Groups        *Groups        `json:"groups,omitempty"`
	Organizations *Organizations `json:"organizations,omitempty"`
}

//CurrentTenant retrieves the Tenant associated with the current API key.
func CurrentTenant() (*Tenant, error) {
	tenant := &Tenant{}

	err := client.get(buildRelativeURL("tenants", "current"), tenant)

	return tenant, err
}

//GetApplications retrieves the collection of all the applications associated with the Tenant.
//
//The collection can be filtered and/or paginated by passing the desire ApplicationCriteria value.
func (tenant *Tenant) GetApplications(criteria ApplicationCriteria) (*Applications, error) {
	apps := &Applications{}

	err := client.get(buildAbsoluteURL(tenant.Applications.Href, criteria.ToQueryString()), apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

//GetAccounts retrieves the collection of all the accounts associated with the Tenant.
//
//The collection can be filtered and/or paginated by passing the desire AccountCriteria value.
func (tenant *Tenant) GetAccounts(criteria AccountCriteria) (*Accounts, error) {
	accounts := &Accounts{}

	err := client.get(buildAbsoluteURL(tenant.Accounts.Href, criteria.ToQueryString()), accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

//GetGroups retrieves the collection of all the groups associated with the Tenant.
//
//The collection can be filtered and/or paginated by passing the desire GroupCriteria value.
func (tenant *Tenant) GetGroups(criteria GroupCriteria) (*Groups, error) {
	groups := &Groups{}

	err := client.get(buildAbsoluteURL(tenant.Groups.Href, criteria.ToQueryString()), groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

//GetDirectories retrieves the collection of all the directories associated with the Tenant.
//
//The collection can be filtered and/or paginated by passing the desire DirectoryCriteria value
func (tenant *Tenant) GetDirectories(criteria DirectoryCriteria) (*Directories, error) {
	directories := &Directories{}

	err := client.get(buildAbsoluteURL(tenant.Directories.Href, criteria.ToQueryString()), directories)
	if err != nil {
		return nil, err
	}

	return directories, nil
}

//GetOrganizations retrieves the collection of all the organizations associated with the Tenant.
//
//The collection can be filtered and/or paginated by passing the desire OrganizationCriteria value
func (tenant *Tenant) GetOrganizations(criteria OrganizationCriteria) (*Organizations, error) {
	organizations := &Organizations{}

	err := client.get(buildAbsoluteURL(tenant.Organizations.Href, criteria.ToQueryString()), organizations)
	if err != nil {
		return nil, err
	}

	return organizations, nil
}
