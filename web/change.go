package stormpathweb

import (
	"fmt"
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
	"golang.org/x/net/context"
)

type changePasswordHandler struct {
	Application *stormpath.Application
}

func (h changePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx context.Context) {
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

func (h changePasswordHandler) doGET(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	contentType := ctx.Value(ResolvedContentType)

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
				"error":    ctx.Value("error"),
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

func (h changePasswordHandler) doPOST(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	contentType := ctx.Value(ResolvedContentType)

	data, _ := getPostedData(r)

	if data["sptoken"] == "" {
		h.handlePostError(w, r, ctx, fmt.Errorf("sptoken is required"))
		return
	}

	if data["password"] == "" {
		h.handlePostError(w, r, ctx, fmt.Errorf("Password is required"))
		return
	}

	if data["confirmPassword"] != "" && data["password"] != data["confirmPassword"] {
		h.handlePostError(w, r, ctx, fmt.Errorf("Password values do not match."))
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
		h.handlePostError(w, r, ctx, err)
		return
	}

	if Config.ChangePasswordAutoLoginEnabled {
		err = saveAuthenticationResult(w, r, transientAuthenticationResult(account), h.Application)
		if err != nil {
			h.handlePostError(w, r, ctx, err)
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

func (h changePasswordHandler) handlePostError(w http.ResponseWriter, r *http.Request, ctx context.Context, err error) {
	contentType := ctx.Value(ResolvedContentType)

	if contentType == stormpath.TextHTML {
		h.doGET(w, r, context.WithValue(ctx, "error", buildErrorModel(err)))
		return
	}
	if contentType == stormpath.ApplicationJSON {
		badRequest(w, r, err)
	}
}
