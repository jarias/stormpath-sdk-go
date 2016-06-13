package stormpathweb

import (
	"strings"

	"github.com/jarias/stormpath-sdk-go"
)

type accountStore struct {
	Name     string   `json:"name"`
	Href     string   `json:"href"`
	Provider provider `json:"provider"`
}

type provider struct {
	Href         string `json:"href"`
	ProviderID   string `json:"providerId"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"_"`
}

func getApplicationAccountStores(application *stormpath.Application) []accountStore {
	accountStores := make([]accountStore, 0, 100)

	//TODO iterate until len(mappings.Items) == 0
	mappings, err := application.GetAccountStoreMappings(stormpath.MakeAccountStoreMappingsCriteria().Limit(100))
	if err != nil {
		stormpath.Logger.Printf("[ERROR] Error getting application's account store mappings: %s \n", err)
		return accountStores
	}

	for _, mapping := range mappings.Items {
		if strings.Contains(mapping.AccountStore.Href, "/directories/") {
			directory, err := stormpath.GetDirectory(mapping.AccountStore.Href, stormpath.MakeDirectoryCriteria().WithProvider())
			if err != nil {
				stormpath.Logger.Printf("[ERROR] Error getting directory: %s \n", err)
				continue
			}
			//TODO add SAML providers
			if directory.Provider.ProviderID != "saml" && directory.Provider.ProviderID != "stormpath" {
				accountStore := accountStore{
					Name: directory.Name,
					Href: directory.Href,
					Provider: provider{
						Href:         directory.Provider.Href,
						ProviderID:   directory.Provider.ProviderID,
						ClientID:     directory.Provider.ClientID,
						ClientSecret: directory.Provider.ClientSecret,
					},
				}

				accountStores = append(accountStores, accountStore)
			}
		}
	}

	return accountStores
}
