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

	resp, err := tenant.Client.Do(NewStormpathRequestNoRedirects(GET, TenantBaseUrl+"/current", PageRequest{}, DefaultFilter{}))

	if err != nil {
		return nil, err
	}

	location := resp.Header.Get(LocationHeader)

	resp, err = tenant.Client.Do(NewStormpathRequest(GET, location, PageRequest{}, DefaultFilter{}))

	if err != nil {
		return nil, err
	}

	err = Unmarshal(resp, tenant)

	return tenant, err
}

func (tenant *Tenant) GetApplications(pageRequest PageRequest, filters DefaultFilter) (*Applications, error) {
	apps := &Applications{}

	resp, err := tenant.Client.Do(NewStormpathRequest(GET, tenant.Applications.Href, pageRequest, filters))

	if err != nil {
		return nil, err
	}

	err = Unmarshal(resp, apps)
	for _, app := range apps.Items {
		app.Client = tenant.Client
	}

	return apps, err
}

func (tenant *Tenant) GetDirectories(pageRequest PageRequest, filters DefaultFilter) (*Directories, error) {
	directories := &Directories{}

	resp, err := tenant.Client.Do(NewStormpathRequest(GET, tenant.Directories.Href, pageRequest, filters))

	if err != nil {
		return nil, err
	}

	err = Unmarshal(resp, directories)
	for _, d := range directories.Items {
		d.Client = tenant.Client
	}

	return directories, err
}
