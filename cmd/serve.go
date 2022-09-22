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
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "port",
			Value: 8080,
			Usage: "Port to use for server",
		},
	},
	Action: func(cCtx *cli.Context) error {
		cfg := api.Config{
			Port: cCtx.Int("port"),
		}

		app := &api.Application{
			Cfg:      cfg,
			DBClient: db.Connect(),
		}

		app.Serve()
		return nil
	},
}
