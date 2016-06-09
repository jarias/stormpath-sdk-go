package stormpath

import "net/url"

type APIKey struct {
	resource
	ID      string   `json:"id"`
	Secret  string   `json:"secret"`
	Status  string   `json:"status"`
	Account *Account `json:"account"`
	Tenant  *Tenant  `json:"tenant"`
}

type APIKeys struct {
	collectionResource
	Items []APIKey `json:"items,omitempty"`
}

type APIKeyCriteria struct {
	baseCriteria
}

func MakeAPIKeyCriteria() APIKeyCriteria {
	return APIKeyCriteria{baseCriteria{filter: url.Values{}}}
}

func MakeAPIKeysCriteria() APIKeyCriteria {
	return APIKeyCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

func GetAPIKey(href string, criteria Criteria) (*APIKey, error) {
	apiKey := &APIKey{}

	err := client.get(
		buildAbsoluteURL(href, criteria.ToQueryString()),
		apiKey,
	)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

func (k *APIKey) Delete() error {
	return client.delete(k.Href)
}

func (k *APIKey) Update() error {
	return client.post(k.Href, map[string]string{"status": k.Status}, k)
}

func (c APIKeyCriteria) WithAccount() APIKeyCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "account")
	return c
}

func (c APIKeyCriteria) WithTenat() APIKeyCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "tenataccount")
	return c
}
