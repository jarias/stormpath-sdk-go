package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
	"golang.org/x/net/context"
)

type linkedinCallbackHandler struct {
	defaultSocialHandler
}

func (h linkedinCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx context.Context) {
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
