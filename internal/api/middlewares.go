package api

import (
	"fmt"
	"net/http"

	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/uptrace/bunrouter"
)

var apiKeyheader = "X-API-Key"

// Check if request has the api key in header
func (server *Server) checkAPIKey(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	apiKeyInConfig := server.AppCfg.API.Key
	if apiKeyInConfig == "" {
		panic(fmt.Sprintf("API Key not found in config at path: %s", cfg.GetConfigPath()))
	}

	return func(w http.ResponseWriter, req bunrouter.Request) error {
		apiKey := req.Header.Get(apiKeyheader)

		if apiKey != apiKeyInConfig {
			server.JSON(
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
