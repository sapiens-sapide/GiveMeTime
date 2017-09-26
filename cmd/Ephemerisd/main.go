package main

import (
	"github.com/sapiens-sapide/GiveMeTime/cmd/Ephemerisd/cmds"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "GiveMeTime ephemeris"
	app.Usage = ""
	app.Commands = []cli.Command{
		{
			Name:    "serve-ephemeris",
			Aliases: []string{"serve"},
			Usage:   "This daemon starts a HTTP server that serves the GiveMeTime ephemeris service.",
			Action: func(c *cli.Context) error {
				cmds.StartEphemerisServer()
				return nil
			},
		},
	}
	app.Run(os.Args)
}
