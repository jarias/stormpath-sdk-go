package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
)

type logoutHandler struct {
	application *stormpath.Application
}

func (h logoutHandler) serveHTTP(w http.ResponseWriter, r *http.Request, ctx webContext) {
	if Config.IDSiteEnabled {
		options := stormpath.IDSiteOptions{
			CallbackURL: baseURL(r) + Config.CallbackURI,
			Logout:      true,
		}
		idSiteURL, err := h.application.CreateIDSiteURL(options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, idSiteURL, http.StatusFound)
		return
	}

	if r.Method == http.MethodPost {
		contentType := ctx.contentType

		if ctx.account != nil {
			clearAuthentication(w, r, h.application)

			if contentType == stormpath.TextHTML {
				http.Redirect(w, r, Config.LogoutNextURI, http.StatusFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		if contentType == stormpath.TextHTML {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	methodNotAllowed(w, r, ctx)
}
