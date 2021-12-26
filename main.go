package main

import (
	"jrnl/cmd"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			cmd.TodayCmd,
			cmd.ListCmd,
			cmd.ShowDbPathCmd,
			cmd.TILCmd,
			cmd.EditCmd,
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}

}
