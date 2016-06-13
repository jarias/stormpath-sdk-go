package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
	"golang.org/x/net/context"
)

const socialAccount = "socialAccount"

type defaultSocialHandler struct {
	Application  *stormpath.Application
	LoginHandler loginHandler
}

func (h defaultSocialHandler) authenticateSocial(w http.ResponseWriter, r *http.Request, ctx context.Context, socialAccount *stormpath.SocialAccount) {
	account, err := h.Application.RegisterSocialAccount(socialAccount)

	if err != nil {
		context.WithValue(ctx, "error", buildErrorModel(err))
		h.LoginHandler.doGET(w, r, ctx)
		return
	}

	err = saveAuthenticationResult(w, r, transientAuthenticationResult(account), h.Application)
	if err != nil {
		context.WithValue(ctx, "error", buildErrorModel(err))
		h.LoginHandler.doGET(w, r, ctx)
		return
	}

	http.Redirect(w, r, Config.LoginNextURI, http.StatusFound)
	return

}
