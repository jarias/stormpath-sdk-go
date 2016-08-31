package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
)

const socialAccount = "socialAccount"

type defaultSocialHandler struct {
	Application  *stormpath.Application
	LoginHandler loginHandler
}

func (h defaultSocialHandler) authenticateSocial(w http.ResponseWriter, r *http.Request, ctx webContext, socialAccount *stormpath.SocialAccount) {
	account, err := h.Application.RegisterSocialAccount(socialAccount)

	if err != nil {
		h.LoginHandler.doGET(w, r, ctx.withError(nil, err))
		return
	}

	err = saveAuthenticationResult(w, r, transientAuthenticationResult(account), h.Application)
	if err != nil {
		h.LoginHandler.doGET(w, r, ctx.withError(nil, err))
		return
	}

	http.Redirect(w, r, Config.LoginNextURI, http.StatusFound)
	return

}
