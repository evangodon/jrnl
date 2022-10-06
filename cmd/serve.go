package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/evangodon/jrnl/internal/api"
	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/evangodon/jrnl/internal/db"
)

var ServeCmd = &cli.Command{
	Name:    "serve",
	Aliases: []string{"s"},
	Usage:   "Start the server",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   8090,
			Usage:   "Port to use for server",
		},
	},
	Action: func(cCtx *cli.Context) error {
		serverCfg := api.ServerConfig{
			Port: cCtx.Int("port"),
			Env:  cfg.GetEnv(),
		}

		app := &api.Server{
			Cfg:      serverCfg,
			DBClient: db.Connect(),
			AppCfg:   cfg.GetConfig(),
		}

		err := app.Serve()
		if err != nil {
			return err
		}
		return nil
	},
}
