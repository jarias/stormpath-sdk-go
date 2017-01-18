package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
)

type meHandler struct {
	application *stormpath.Application
}

func (h meHandler) serveHTTP(w http.ResponseWriter, r *http.Request, ctx webContext) {
	if r.Method == http.MethodGet {
		if ctx.account != nil {
			w.Header().Set("Cache-Control", "no-store, no-cache")
			w.Header().Set("Pragma", "no-cache")

			respondJSON(w, accountModel(ctx.account), http.StatusOK)
			return
		}
		unauthorizedRequest(w, r, ctx, h.application)
		return
	}

	methodNotAllowed(w, r, ctx)
}
