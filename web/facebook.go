package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
	"golang.org/x/net/context"
)

type facebookCallbackHandler struct {
	defaultSocialHandler
}

func (h facebookCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if r.Method == http.MethodGet {
		accessToken := r.URL.Query().Get("accessToken")

		socialAccount := &stormpath.SocialAccount{
			Data: stormpath.ProviderData{
				ProviderID:  "facebook",
				AccessToken: accessToken,
			},
		}

		h.authenticateSocial(w, r, ctx, socialAccount)
		return
	}

	methodNotAllowed(w, r, ctx)
}
