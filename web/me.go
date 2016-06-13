package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
	"golang.org/x/net/context"
)

type meHandler struct {
}

func (h meHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if r.Method == http.MethodGet {
		if ctx, ok := isAuthenticated(w, r, ctx); ok {
			w.Header().Set("Cache-Control", "no-store, no-cache")
			w.Header().Set("Pragma", "no-cache")

			account := expandAccountAttributes(ctx.Value(AccountKey).(*stormpath.Account))
			respondJSON(w, accountModel(account), http.StatusOK)
			return
		}
		unauthorizedRequest(w, r, ctx)
		return
	}

	methodNotAllowed(w, r, ctx)
}

func expandAccountAttributes(account *stormpath.Account) *stormpath.Account {
	criteria := stormpath.MakeAccountCriteria()

	for attribute, shouldExpand := range Config.MeExpand {
		switch attribute {
		case "apiKeys":
			if shouldExpand.(bool) {
				criteria = criteria.WithAPIKeys()
			}
			break
		case "applications":
			if shouldExpand.(bool) {
				criteria = criteria.WithApplications()
			}
			break
		case "customData":
			if shouldExpand.(bool) {
				criteria = criteria.WithCustomData()
			}
			break
		case "directory":
			if shouldExpand.(bool) {
				criteria = criteria.WithDirectory()
			}
			break
		case "groupMemberships":
			if shouldExpand.(bool) {
				criteria = criteria.WithGroupMemberships(stormpath.DefaultPageRequest)
			}
			break
		case "groups":
			if shouldExpand.(bool) {
				criteria = criteria.WithGroups(stormpath.DefaultPageRequest)
			}
			break
		case "providerData":
			if shouldExpand.(bool) {
				criteria = criteria.WithProviderData()
			}
			break
		case "tenant":
			if shouldExpand.(bool) {
				criteria = criteria.WithTenant()
			}
			break
		}
	}

	expandedAccount, err := stormpath.GetAccount(account.Href, criteria)
	if err != nil {
		return account
	}

	return expandedAccount
}
