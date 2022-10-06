package cmd

import (
	"fmt"
	"net/http"

	"github.com/evangodon/jrnl/internal/api"
	"github.com/evangodon/jrnl/internal/cfg"

	"github.com/urfave/cli/v2"
)

var ShowPathsCmd = &cli.Command{
	Name:  "showpaths",
	Usage: "Show paths to config and database",
	Action: func(_ *cli.Context) error {

		fmt.Println("Config:   ", cfg.GetConfigPath())

		client := api.Client{
			Config: cfg.GetConfig(),
		}
		res, err := client.MakeRequest(http.MethodGet, "/dbpath", nil)
		if err != nil {
			return err
		}

		fmt.Println("Database: " + string(res.Body))
		fmt.Println("Sending requests to: ", cfg.GetConfig().API.BaseURL)

		return nil
	},
}
