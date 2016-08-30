package stormpathweb

import (
	"net/http"

	"fmt"

	"github.com/jarias/stormpath-sdk-go"
)

type emailVerifyHandler struct {
	Application *stormpath.Application
}

func (h emailVerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx webContext) {
	if ctx.Account != nil {
		http.Redirect(w, r, Config.LoginNextURI, http.StatusFound)
		return
	}

	if IsVerifyEnabled(h.Application) {
		if r.Method == http.MethodPost {
			h.doPOST(w, r, ctx)
			return
		}
		if r.Method == http.MethodGet {
			h.doGET(w, r, ctx)
			return
		}

		methodNotAllowed(w, r, ctx)
	}
}

func (h emailVerifyHandler) doGET(w http.ResponseWriter, r *http.Request, ctx webContext) {
	contentType := ctx.ContentType

	spToken := r.URL.Query().Get("sptoken")

	if spToken != "" {
		//Validate the token
		account, err := stormpath.VerifyEmailToken(spToken)
		if err != nil {
			if contentType == stormpath.TextHTML {
				handleError(w, r, ctx.withError(nil, fmt.Errorf("This verification link is no longer valid. Please request a new link from the form below.")), h.doGET)
				return
			}
			if contentType == stormpath.ApplicationJSON {
				handleError(w, r, ctx.withError(nil, err), h.doGET)
				return
			}
		}
		var redirectSuccessURI = Config.VerifyNextURI

		if Config.RegisterAutoLoginEnabled {
			//AutoLogin
			err := saveAuthenticationResult(w, r, transientAuthenticationResult(account), h.Application)
			if err != nil {
				handleError(w, r, ctx.withError(nil, err), h.doGET)
				return
			}

			redirectSuccessURI = Config.LoginNextURI
		}
		if contentType == stormpath.TextHTML {
			http.Redirect(w, r, redirectSuccessURI, http.StatusFound)
			return
		}
		if contentType == stormpath.ApplicationJSON {
			respondJSON(w, nil, http.StatusOK)
			return
		}
	}
	model := map[string]interface{}{}

	model["error"] = ctx.Error

	if contentType == stormpath.TextHTML {
		model["loginURI"] = Config.LoginURI
		respondHTML(w, model, Config.VerifyView)
		return
	}

	if contentType == stormpath.ApplicationJSON {
		handleError(w, r, ctx.withError(nil, fmt.Errorf("sptoken parameter not provided")), h.doGET)
	}
}

func (h emailVerifyHandler) doPOST(w http.ResponseWriter, r *http.Request, ctx webContext) {
	contentType := ctx.ContentType

	data, _ := getPostedData(r)

	if data["email"] == "" {
		handleError(w, r, ctx.withError(nil, fmt.Errorf("email is required")), h.doGET)
		return
	}

	h.Application.ResendVerificationEmail(data["email"])

	if contentType == stormpath.ApplicationJSON {
		respondJSON(w, nil, http.StatusOK)
		return
	}

	if contentType == stormpath.TextHTML {
		http.Redirect(w, r, Config.VerifyNextURI, http.StatusFound)
	}
}
