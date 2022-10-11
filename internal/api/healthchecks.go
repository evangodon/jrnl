package api

import (
	"errors"
	"net/http"
	"os"

	"github.com/uptrace/bunrouter"

	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/evangodon/jrnl/internal/db"
)

func (app *Server) healthcheck(w http.ResponseWriter, _ bunrouter.Request) error {

	checks := map[string]string{
		"database": "",
		"config":   "",
	}

	dbPath := db.GetDBPath()
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		checks["database"] = "Not found at " + dbPath
	} else {
		checks["database"] = "Exists at " + dbPath
	}

	configPath := cfg.GetConfigPath()
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		checks["config"] = "Not found at " + configPath
	} else {
		checks["config"] = "Exists at " + configPath
	}

	app.JSON(
		w,
		http.StatusOK,
		Envelope{"checks": checks, "app_env": app.cfg.Env},
		nil,
	)
	return nil
}

func (app *Server) dbPath(w http.ResponseWriter, _ bunrouter.Request) error {
	dbPath := db.GetDBPath()
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("database was not found at " + dbPath))
		return nil
	}

	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte(dbPath))

	return nil
}
