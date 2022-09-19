package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/evangodon/jrnl/internal/api"
	"github.com/evangodon/jrnl/internal/db"
)

var ServeCmd = &cli.Command{
	Name:    "serve",
	Aliases: []string{"s"},
	Usage:   "Start the server",
	Action: func(_ *cli.Context) error {
		cfg := api.Config{
			Port: 8080,
		}

		app := &api.Application{
			Cfg:      cfg,
			DBClient: db.Connect(),
		}

		app.Serve()
		return nil
	},
}
