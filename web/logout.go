package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
	"golang.org/x/net/context"
)

type logoutHandler struct {
	Application *stormpath.Application
}

func (h logoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if Config.IDSiteEnabled {
		options := stormpath.IDSiteOptions{
			CallbackURL: baseURL(r) + Config.CallbackURI,
			Logout:      true,
		}
		idSiteURL, err := h.Application.CreateIDSiteURL(options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, idSiteURL, http.StatusFound)
		return
	}

	if r.Method == http.MethodPost {
		if _, ok := isAuthenticated(w, r, ctx); ok {
			clearAuthentication(w, r, h.Application)

			if ctx.Value(ResolvedContentType) == stormpath.TextHTML {
				http.Redirect(w, r, Config.LogoutNextURI, http.StatusFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		if ctx.Value(ResolvedContentType) == stormpath.TextHTML {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	methodNotAllowed(w, r, ctx)
}
