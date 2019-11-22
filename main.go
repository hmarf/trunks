package main

import (
	"os"
	"strings"

	"github.com/hmarf/trunks/trunks"
	"github.com/urfave/cli"
)

func App() *cli.App {
	app := cli.NewApp()
	app.Name = "trunks"
	app.Usage = "Trunks is a simple command line tool for HTTP load testing."
	app.Version = "0.0.1"
	app.Author = "hmarf"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "concurrency, c",
			Value: 10,
			Usage: "Concurrency Level",
		},
		cli.IntFlag{
			Name:  "requests, r",
			Value: 100,
			Usage: "Number of Requests",
		},
		cli.StringFlag{
			Name:  "url, u",
			Value: "None",
			Usage: "URL to hit",
		},
	}
	return app
}

func Action(c *cli.Context) {
	app := App()
	if c.String("url") == "None" || !strings.HasPrefix(c.String("url"), "http") {
		app.Run(os.Args)
		return
	}
	trunks.Trunks(c.Int("concurrency"), c.Int("requests"), c.String("url"))
}

func main() {
	app := App()
	app.Action = Action
	app.Run(os.Args)
}
