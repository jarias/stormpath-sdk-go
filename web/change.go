package stormpathweb

import (
	"fmt"
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
)

type changePasswordHandler struct {
	Application *stormpath.Application
}

func (h changePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx webContext) {
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

func (h changePasswordHandler) doGET(w http.ResponseWriter, r *http.Request, ctx webContext) {
	contentType := ctx.ContentType

	sptoken := r.URL.Query().Get("sptoken")

	if sptoken != "" {
		_, err := h.Application.ValidatePasswordResetToken(sptoken)
		if err != nil {
			if contentType == stormpath.TextHTML {
				http.Redirect(w, r, Config.ChangePasswordErrorURI, http.StatusFound)
				return
			}
			if contentType == stormpath.ApplicationJSON {
				badRequest(w, r, err)
				return
			}
		}
		if contentType == stormpath.TextHTML {
			model := map[string]interface{}{
				"loginURI": Config.LoginURI,
				"error":    ctx.Error,
			}
			respondHTML(w, model, Config.ChangePasswordView)
			return
		}
		if contentType == stormpath.ApplicationJSON {
			respondJSON(w, nil, http.StatusOK)
			return
		}
	}
	if contentType == stormpath.TextHTML {
		http.Redirect(w, r, Config.ForgotPasswordURI, http.StatusFound)
		return
	}
	if contentType == stormpath.ApplicationJSON {
		badRequest(w, r, fmt.Errorf("sptoken parameter not provided"))
	}
}

func (h changePasswordHandler) doPOST(w http.ResponseWriter, r *http.Request, ctx webContext) {
	contentType := ctx.ContentType

	data, _ := getPostedData(r)

	if data["sptoken"] == "" {
		handleError(w, r, ctx.withError(nil, fmt.Errorf("sptoken is required")), h.doGET)
		return
	}

	if data["password"] == "" {
		handleError(w, r, ctx.withError(nil, fmt.Errorf("Password is required")), h.doGET)
		return
	}

	if data["confirmPassword"] != "" && data["password"] != data["confirmPassword"] {
		handleError(w, r, ctx.withError(data, fmt.Errorf("Password values do not match.")), h.doGET)
		return
	}

	_, err := h.Application.ValidatePasswordResetToken(data["sptoken"])
	if err != nil {
		if contentType == stormpath.TextHTML {
			http.Redirect(w, r, Config.ChangePasswordErrorURI, http.StatusFound)
			return
		}
		if contentType == stormpath.ApplicationJSON {
			badRequest(w, r, err)
			return
		}
	}

	account, err := h.Application.ResetPassword(data["sptoken"], data["password"])
	if err != nil {
		handleError(w, r, ctx.withError(nil, err), h.doGET)
		return
	}

	if Config.ChangePasswordAutoLoginEnabled {
		err = saveAuthenticationResult(w, r, transientAuthenticationResult(account), h.Application)
		if err != nil {
			handleError(w, r, ctx.withError(nil, err), h.doGET)
			return
		}

		if contentType == stormpath.TextHTML {
			http.Redirect(w, r, Config.LoginNextURI, http.StatusFound)
			return
		}
		if contentType == stormpath.ApplicationJSON {
			respondJSON(w, accountModel(account), http.StatusOK)
			return
		}
	}
	if contentType == stormpath.TextHTML {
		http.Redirect(w, r, Config.ChangePasswordNextURI, http.StatusFound)
		return
	}
	if contentType == stormpath.ApplicationJSON {
		respondJSON(w, nil, http.StatusOK)
		return
	}
}
