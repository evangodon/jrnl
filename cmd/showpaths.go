package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/evangodon/jrnl/internal/api"
	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/evangodon/jrnl/internal/logger"

	"github.com/urfave/cli/v2"
)

var ShowPathsCmd = &cli.Command{
	Name:  "showpaths",
	Usage: "Show paths to config and database",
	Action: func(_ *cli.Context) error {
		logger := logger.NewLogger(os.Stdout)

		config := fmt.Sprintf("%-12s %s", "Config", cfg.GetConfigPath())
		logger.PrintInfo(config)

		client := api.Client{
			Config: cfg.GetConfig(),
		}
		res, err := client.MakeRequest(http.MethodGet, "/dbpath", nil)
		if err != nil {
			return err
		}

		db := fmt.Sprintf("%-12s %s", "Database", string(res.Body))
		logger.PrintInfo(db)

		api := fmt.Sprintf("%-12s %s", "Server URL: ", cfg.GetConfig().API.BaseURL)
		logger.PrintInfo(api)

		return nil
	},
}
