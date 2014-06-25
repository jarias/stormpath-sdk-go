package stormpath

const (
	TenantBaseUrl  = "https://api.stormpath.com/v1/tenants"
	LocationHeader = "Location"
)

type Tenant struct {
	Href         string           `json:"href"`
	Name         string           `json:"name"`
	Key          string           `json:"key"`
	Applications Link             `json:"applications"`
	Directories  Link             `json:"directories"`
	Client       *StormpathClient `json:"-"`
}

func CurrentTenant(credentials *Credentials) (*Tenant, error) {
	tenant := &Tenant{Client: NewStormpathClient(credentials)}

	resp, err := tenant.Client.Do(NewStormpathRequestNoRedirects(GET, TenantBaseUrl+"/current", PageRequest{}, ApplicationFilter{}))

	if err != nil {
		return nil, err
	}

	location := resp.Header.Get(LocationHeader)

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
