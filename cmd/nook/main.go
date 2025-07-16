package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var Version = "dev"

func main() {
	app := cli.App{
		Name:     "nook",
		Usage:    "My PDS software, made in Go",
		Version:  Version,
		Commands: []*cli.Command{runCmd},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

var runCmd = &cli.Command{
	Name:    "run",
	Aliases: []string{"r"},
	Usage:   "Starts the PDS service",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "bind-addr",
			Usage:       "",
			Value:       ":8080",
			DefaultText: ":8080",
			EnvVars:     []string{"NOOK_BIND_ADDR"},
		},
	},

	Action: func(ctx *cli.Context) error {
		return nil
	},
}
