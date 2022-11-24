package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/evangodon/jrnl/internal/api"
	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/evangodon/jrnl/internal/db"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Setenv("TEST", "true")
	dbPath := db.GetDBPath()

	if _, err := os.Stat(dbPath); err == nil {
		err := os.Remove(dbPath)
		require.NoError(t, err)
	}

	err := db.CreateNewDB(dbPath)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = os.Remove(dbPath)
		if err != nil {
			t.Error(err)
		}
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
		cfg:      serverCfg,
		dbClient: db.Connect(),
		appCfg:   appCfg,
	}

	testsrv := httptest.NewServer(app.routes())
	defer testsrv.Close()

	appCfg.API.BaseURL = testsrv.URL
	client := api.Client{Config: appCfg}

	testEntry := db.Journal{
		Content: "A new entry for today",
		Date:    time.Now(),
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

	t.Run("should be able to update the entry for today", func(t *testing.T) {
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
		fmt.Printf("ID: %s\n", gotEntry.Daily.ID)

		updatedEntry := db.Journal{
			ID:      gotEntry.Daily.ID,
			Date:    gotEntry.Daily.Date,
			Content: gotEntry.Daily.Content + "\nand here's another thing",
		}
		payload, err := json.Marshal(updatedEntry)
		if err != nil {
			t.Error(err)
		}

		res, err = client.MakeRequest(
			http.MethodPatch,
			"/daily/"+updatedEntry.ID,
			bytes.NewBuffer(payload),
		)
		require.NoError(t, err)

		require.Equal(t, res.Status, http.StatusNoContent)
	})
}
