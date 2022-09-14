package main

import (
	"log"
	"os"

	"github.com/evangodon/jrnl/cmd"

	"github.com/urfave/cli/v2"
)

func run(args []string) error {
	app := &cli.App{
		Commands: []*cli.Command{
			cmd.EditCmd,
			cmd.ListCmd,
			cmd.NewCmd,
			cmd.ShowDBPathCmd,
			cmd.TodayCmd,
			cmd.YesterdayCmd,
		},
	}

	return app.Run(args)
}

func main() {
	err := run(os.Args)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
