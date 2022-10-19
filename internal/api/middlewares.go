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
	apiKeyInConfig := server.appCfg.API.Key
	if apiKeyInConfig == "" {
		server.logger.PrintFatal(
			fmt.Sprintf("API Key not found in config at path: %s", cfg.GetConfigPath()),
		)
	}

	return func(w http.ResponseWriter, req bunrouter.Request) error {
		apiKey := req.Header.Get(apiKeyheader)

		if apiKey != apiKeyInConfig {
			println("Not the right key", apiKey)
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

func corsMiddleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			return next(w, req)
		}

		h := w.Header()

		h.Set("Access-Control-Allow-Origin", origin)
		h.Set("Access-Control-Allow-Credentials", "true")

		// CORS preflight.
		if req.Method == http.MethodOptions {
			h.Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,HEAD")
			h.Set("Access-Control-Allow-Headers", "authorization,content-type,x-api-key")
			h.Set("Access-Control-Max-Age", "86400")
			return nil
		}

		return next(w, req)
	}
}
