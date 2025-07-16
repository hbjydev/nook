package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/hbjydev/nook/pkg/server"
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

		&cli.StringFlag{
			Name:     "did",
			Required: true,
			Usage:    "The DID this Nook serves.",
			EnvVars:  []string{"NOOK_DID"},
		},

		&cli.StringFlag{
			Name:     "hostname",
			Required: true,
			Usage:    "The public-facing hostname Nook is available at",
			EnvVars:  []string{"NOOK_HOSTNAME"},
		},
	},

	Action: func(ctx *cli.Context) error {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

		srv, err := server.New(server.Args{
			Logger: logger,

			BindAddr: ctx.String("bind-addr"),
			Did:      ctx.String("did"),
			Hostname: ctx.String("hostname"),
		})
		if err != nil {
			return err
		}

		return srv.Run(ctx.Context)
	},
}
