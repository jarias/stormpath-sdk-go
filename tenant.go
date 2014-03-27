package stormpath

import (
	"encoding/json"
	"io/ioutil"
)

const tenantBaseUrl = "https://api.stormpath.com/v1/tenants"

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

	resp, err := tenant.Client.Get(tenantBaseUrl+"/current", false)

	if err != nil {
		return nil, err
	}

	location := resp.Header.Get("Location")

	resp, err = tenant.Client.Get(location, true)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, tenant)

	if err != nil {
		return nil, err
	}

	return tenant, err
}
