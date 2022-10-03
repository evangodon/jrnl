package api

import (
	"net/http"

	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/uptrace/bunrouter"
)

var apiKeyheader = "X-API-Key"

// Check if request has the api key in header
func (app *Application) checkAPIKey(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	apiKeyInConfig := cfg.GetConfig().API.Key
	if apiKeyInConfig == "" {
		panic("API Key not found in config")
	}

	return func(w http.ResponseWriter, req bunrouter.Request) error {
		apiKey := req.Header.Get(apiKeyheader)

		if apiKey != apiKeyInConfig {
			app.writeJSON(
				w,
				http.StatusUnauthorized,
				Envelope{"msg": "Unauthorized"},
				nil,
			)

			return nil
		}

		return next(w, req)
	}
}
