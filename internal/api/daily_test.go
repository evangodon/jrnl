package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/evangodon/jrnl/internal/db"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Setenv("TEST", "true")
	dbPath := db.GetDBPath()
	db.CreateNewDB(dbPath)
	t.Cleanup(func() {
		os.Remove(dbPath)
	})

	serverCfg := ServerConfig{
		Env: "TEST",
	}

	appCfg := cfg.Config{
		API: cfg.API{
			Key: "test-api-key-123",
		},
	}

	app := Server{
		Cfg:      serverCfg,
		DBClient: db.Connect(),
		AppCfg:   appCfg,
	}

	testsrv := httptest.NewServer(app.routes())
	defer testsrv.Close()

	appCfg.API.BaseURL = testsrv.URL
	client := Client{Config: appCfg}

	testEntry := db.Journal{
		Content: "A new entry for today",
	}

	t.Run("should be able to create a daily entry for today", func(t *testing.T) {
		newEntry, err := json.Marshal(testEntry)
		if err != nil {
			t.Error(err)
		}

		res, err := client.MakeRequest(http.MethodPost, "/daily/new", bytes.NewBuffer(newEntry))
		require.NoError(t, err)

		require.Equal(t, res.Status, http.StatusCreated)
	})

	t.Run("should be able to request the new entry", func(t *testing.T) {
		res, err := client.MakeRequest(http.MethodGet, "/daily/", nil)
		require.NoError(t, err)
		require.Equal(t, res.Status, http.StatusOK)

		gotEntry := struct {
			Daily db.Journal `json:"daily"`
		}{
			Daily: db.Journal{},
		}
		err = json.Unmarshal(res.Body, &gotEntry)
		require.NoError(t, err)

		require.Equal(t, testEntry.Content, gotEntry.Daily.Content)
	})
}
