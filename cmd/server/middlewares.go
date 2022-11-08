package main

import (
	"fmt"
	"net/http"

	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/uptrace/bunrouter"
)

var apiKeyheader = "X-API-Key"

// Check if request has the api key in header
func (srv *Server) validateAPIKeyMiddleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	apiKeyInConfig := srv.appCfg.API.Key
	if apiKeyInConfig == "" {
		srv.logger.PrintFatal(
			fmt.Sprintf("API Key not found in config at path: %s", cfg.GetConfigPath()),
		)
	}

	return func(w http.ResponseWriter, req bunrouter.Request) error {
		apiKey := req.Header.Get(apiKeyheader)

		if apiKey != apiKeyInConfig {
			return HTTPError{
				statusCode: http.StatusUnauthorized,
				Code:       "unauthorized",
				Message:    "Unauthorized",
			}
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

func errorHandler(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		err := next(w, req)

		switch err := err.(type) {
		case nil:
		case HTTPError:
			w.WriteHeader(err.statusCode)
			_ = bunrouter.JSON(w, err)
		default:
			httpErr := NewHTTPError(err)
			w.WriteHeader(httpErr.statusCode)
			_ = bunrouter.JSON(w, httpErr)
		}

		return err
	}
}
