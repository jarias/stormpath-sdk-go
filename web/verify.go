package stormpathweb

import (
	"net/http"

	"fmt"

	"github.com/jarias/stormpath-sdk-go"
	"golang.org/x/net/context"
)

type emailVerifyHandler struct {
	Application *stormpath.Application
}

func (h emailVerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if _, ok := isAuthenticated(w, r, ctx); ok {
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

func (h emailVerifyHandler) doGET(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	contentType := ctx.Value(ResolvedContentType)

	spToken := r.URL.Query().Get("sptoken")

	if spToken != "" {
		//Validate the token
		account, err := stormpath.VerifyEmailToken(spToken)
		if err != nil {
			if contentType == stormpath.TextHTML {
				h.handleGetError(w, r, ctx, fmt.Errorf("This verification link is no longer valid. Please request a new link from the form below."))
				return
			}
			if contentType == stormpath.ApplicationJSON {
				h.handleGetError(w, r, ctx, err)
				return
			}
		}
		var redirectSuccessURI = Config.VerifyNextURI

		if Config.RegisterAutoLoginEnabled {
			//AutoLogin
			err := saveAuthenticationResult(w, r, transientAuthenticationResult(account), h.Application)
			if err != nil {
				h.handleGetError(w, r, ctx, err)
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

	if ctx.Value("error") != nil {
		model["error"] = ctx.Value("error")
	}

	if contentType == stormpath.TextHTML {
		model["loginURI"] = Config.LoginURI
		respondHTML(w, model, Config.VerifyView)
		return
	}

	if contentType == stormpath.ApplicationJSON {
		h.handleGetError(w, r, ctx, fmt.Errorf("sptoken parameter not provided"))
	}
}

func (h emailVerifyHandler) doPOST(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	contentType := ctx.Value(ResolvedContentType)

	data, _ := getPostedData(r)

	if data["email"] == "" {
		h.handlePostError(w, r, ctx, fmt.Errorf("email is required"))
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

func (h emailVerifyHandler) handleGetError(w http.ResponseWriter, r *http.Request, ctx context.Context, err error) {
	contentType := ctx.Value(ResolvedContentType)

	if contentType == stormpath.TextHTML {
		r.URL.RawQuery = ""
		h.doGET(w, r, context.WithValue(ctx, "error", buildErrorModel(err)))
		return
	}
	if contentType == stormpath.ApplicationJSON {
		badRequest(w, r, err)
	}
}

func (h emailVerifyHandler) handlePostError(w http.ResponseWriter, r *http.Request, ctx context.Context, err error) {
	contentType := ctx.Value(ResolvedContentType)

	if contentType == stormpath.TextHTML {
		h.doGET(w, r, context.WithValue(ctx, "error", buildErrorModel(err)))
		return
	}
	if contentType == stormpath.ApplicationJSON {
		badRequest(w, r, err)
	}
}
