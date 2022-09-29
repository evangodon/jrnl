package cmd

import (
	"fmt"

	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/evangodon/jrnl/internal/db"

	"github.com/urfave/cli/v2"
)

var ShowPathsCmd = &cli.Command{
	Name:  "showpaths",
	Usage: "Show paths to config and data",
	Action: func(_ *cli.Context) error {
		fmt.Println("\nDatabase: ", db.GetDBPath())
		fmt.Println("Config:   ", cfg.GetConfigPath())
		fmt.Println("Sending requests to: ", cfg.GetConfig().API.BaseURL)

		return nil
	},
}
