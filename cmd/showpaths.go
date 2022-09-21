package cmd

import (
	"fmt"

	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/evangodon/jrnl/internal/db"

	"github.com/urfave/cli/v2"
)

var ShowPathsCmd = &cli.Command{
	Name:  "showpaths",
	Usage: "Show the path to the database",
	Action: func(_ *cli.Context) error {
		fmt.Println("\nDatabase: ", db.GetDBPath())
		fmt.Println("Config:   ", cfg.GetConfigPath())

		return nil
	},
}
