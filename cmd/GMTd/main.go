package main

import (
	"github.com/sapiens-sapide/GiveMeTime/cmd/GMTd/cmds"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "GMT"
	app.Usage = ""
	app.Commands = []cli.Command{
		{
			Name:    "serve-svg",
			Aliases: []string{"svg"},
			Usage:   "This daemon starts a HTTP server that serves the GiveMeTime watch interface as SVG.",
			Action: func(c *cli.Context) error {
				cmds.StartSVGServer()
				return nil
			},
		},
	}
	app.Run(os.Args)
}
