package main

import (
	"jrnl/cmd"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func run(args []string) error {
	app := &cli.App{
		Commands: []*cli.Command{
			cmd.TodayCmd,
			cmd.ListCmd,
			cmd.ShowDbPathCmd,
			cmd.TILCmd,
			cmd.EditCmd,
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
