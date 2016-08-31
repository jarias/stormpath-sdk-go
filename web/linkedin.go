package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
)

type linkedinCallbackHandler struct {
	defaultSocialHandler
}

func (h linkedinCallbackHandler) serveHTTP(w http.ResponseWriter, r *http.Request, ctx webContext) {
	if r.Method == http.MethodGet {
		code := r.URL.Query().Get("code")

		socialAccount := &stormpath.SocialAccount{
			Data: stormpath.ProviderData{
				ProviderID: "linkedin",
				Code:       code,
			},
		}

		h.authenticateSocial(w, r, ctx, socialAccount)
		return
	}

	methodNotAllowed(w, r, ctx)
}
