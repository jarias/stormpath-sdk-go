package stormpath

const (
	TENANT_BASE_URL = "https://api.stormpath.com/v1/tenants"
	LOCATION_HEADER = "Location"
)

type Tenant struct {
	Href         string
	Name         string
	Key          string
	Applications struct {
		Href string
	}
	Directories struct {
		Href string
	}
	Client *StormpathClient
}

func CurrentTenant(credentials *Credentials) (*Tenant, error) {
	tenant := &Tenant{Client: NewStormpathClient(credentials)}

	resp, err := tenant.Client.Do(NewStormpathRequestNoRedirects(GET, TENANT_BASE_URL+"/current", PageRequest{}, ApplicationFilter{}))

	if err != nil {
		return nil, err
	}

	location := resp.Header.Get(LOCATION_HEADER)

	resp, err = tenant.Client.Do(NewStormpathRequest(GET, location, PageRequest{}, ApplicationFilter{}))

	if err != nil {
		return nil, err
	}

	err = Unmarshal(resp, tenant)

	return tenant, err
}

func (tenant *Tenant) GetApplications(pageRequest PageRequest, filters ApplicationFilter) (*Applications, error) {
	apps := &Applications{}

	resp, err := tenant.Client.Do(NewStormpathRequest(GET, tenant.Applications.Href, pageRequest, filters))

	if err != nil {
		return nil, err
	}

	err = Unmarshal(resp, apps)

	return apps, err
}
