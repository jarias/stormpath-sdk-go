package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
)

type callbackHandler struct {
	Application *stormpath.Application
}

func (h callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx webContext) {
	if r.Method == http.MethodGet {
		authenticationResult, err := stormpath.NewStormpathAssertionAuthenticator(h.Application).Authenticate(r.URL.Query().Get("jwtResponse"))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		switch authenticationResult.Status {
		case "AUTHENTICATED":
			err = saveAuthenticationResult(w, r, authenticationResult, h.Application)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, Config.LoginNextURI, http.StatusFound)
			return
		case "REGISTERED":
			accountStatus := authenticationResult.Account.Status

			if accountStatus == stormpath.Enabled {
				if Config.RegisterAutoLoginEnabled {
					err := saveAuthenticationResult(w, r, authenticationResult, h.Application)
					if err != nil {
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
						return
					}
				}
				http.Redirect(w, r, Config.LoginURI+"?status=created", http.StatusFound)
				return
			} else if accountStatus == stormpath.Unverified {
				http.Redirect(w, r, Config.LoginURI+"?status=unverified", http.StatusFound)
				return
			}
			http.Redirect(w, r, Config.RegisterNextURI, http.StatusFound)
			return
		case "LOGOUT":
			clearAuthentication(w, r, h.Application)
			http.Redirect(w, r, Config.LogoutNextURI, http.StatusFound)
			return
		}
	}

	methodNotAllowed(w, r, ctx)
}
