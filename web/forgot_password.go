package stormpathweb

import (
	"fmt"
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
	"golang.org/x/net/context"
)

type forgotPassordHandler struct {
	Application *stormpath.Application
}

func (h forgotPassordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if _, ok := isAuthenticated(w, r, ctx); ok {
		http.Redirect(w, r, Config.LoginNextURI, http.StatusFound)
		return
	}

	if Config.IDSiteEnabled {
		options := stormpath.IDSiteOptions{
			Path:        Config.IDSiteForgotURI,
			CallbackURL: baseURL(r) + Config.CallbackURI,
		}

		idSiteURL, err := h.Application.CreateIDSiteURL(options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, idSiteURL, http.StatusFound)
		return
	}

	if r.Method == http.MethodPost {
		h.doPOST(w, r, ctx)
		return
	}
	//No GET for application/json
	if r.Method == http.MethodGet && ctx.Value(ResolvedContentType) == stormpath.TextHTML {
		h.doGET(w, r, ctx)
		return
	}

	methodNotAllowed(w, r, ctx)
}

func (h forgotPassordHandler) doGET(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	model := map[string]interface{}{
		"loginURI": Config.LoginURI,
		"status":   resolveForgotPasswordStatus(r.URL.Query().Get("status")),
		"error":    ctx.Value("error"),
	}

	respondHTML(w, model, Config.ForgotPasswordView)
}

func (h forgotPassordHandler) doPOST(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	contentType := ctx.Value(ResolvedContentType)

	data, _ := getPostedData(r)

	if data["email"] == "" {
		h.handlePostError(w, r, ctx, fmt.Errorf("email is required"))
		return
	}

	h.Application.SendPasswordResetEmail(data["email"])

	if contentType == stormpath.ApplicationJSON {
		respondJSON(w, nil, http.StatusOK)
		return
	}

	if contentType == stormpath.TextHTML {
		http.Redirect(w, r, Config.ForgotPasswordNextURI, http.StatusFound)
		return
	}
}

func (h forgotPassordHandler) handlePostError(w http.ResponseWriter, r *http.Request, ctx context.Context, err error) {
	contentType := ctx.Value(ResolvedContentType)

	if contentType == stormpath.TextHTML {
		h.doGET(w, r, context.WithValue(ctx, "error", buildErrorModel(err)))
		return
	}
	if contentType == stormpath.ApplicationJSON {
		badRequest(w, r, err)
	}
}

func resolveForgotPasswordStatus(status string) string {
	statusMessage := ""

	switch status {
	case "invalid_sptoken":
		statusMessage = "The password reset link you tried to use is no longer valid. Please request a new link from the form below."
		break
	}

	return statusMessage
}
