package stormpath

import (
	"strings"
	"time"
)

//collectionResource represent the basic attributes of collection of resources (Application, Group, Account, etc.)
type collectionResource struct {
	Href       string     `json:"href,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	ModifiedAt *time.Time `json:"modifiedAt,omitempty"`
	Offset     *int       `json:"offset,omitempty"`
	Limit      *int       `json:"limit,omitempty"`
}

func (r collectionResource) IsCacheable() bool {
	return false
}

func (r collectionResource) GetOffset() int {
	return *r.Offset
}

func (r collectionResource) GetLimit() int {
	return *r.Limit
}

//resource resprents the basic attributes of any resource (Application, Group, Account, etc.)
type resource struct {
	Href       string     `json:"href,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	ModifiedAt *time.Time `json:"modifiedAt,omitempty"`
}

func (r resource) IsCacheable() bool {
	return true
}

//Delete deletes the given account, it wont modify the calling account
func (r *resource) Delete() error {
	return client.delete(r.Href)
}

type accountStoreResource struct {
	customDataAwareResource
	Accounts *Accounts `json:"accounts,omitempty"`
}

//GetAccounts returns the accounts within a context of:
//application, directory, group, organization.
//
//See: http://docs.stormpath.com/rest/product-guide/latest/accnt_mgmt.html#how-to-search-accounts
func (r *accountStoreResource) GetAccounts(criteria Criteria) (*Accounts, error) {
	accounts := &Accounts{}

	err := client.get(
		buildAbsoluteURL(r.Accounts.Href, criteria.ToQueryString()),
		accounts,
	)

	return accounts, err
}

func GetToken(href string) string {
	return href[strings.LastIndex(href, "/")+1:]
}
