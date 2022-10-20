package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/evangodon/jrnl/internal/api"
	"github.com/evangodon/jrnl/internal/cfg"
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

		srv := api.NewServer(serverCfg)

		err := srv.Serve()
		if err != nil {
			return err
		}
		return nil
	},
}
