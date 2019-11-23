package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/hmarf/trunks/trunks"
	"github.com/hmarf/trunks/trunks/attack"
	"github.com/urfave/cli"
)

func headerSplit(header string) []string {
	re := regexp.MustCompile(`^([\w-]+):\s*(.+)`)
	return re.FindStringSubmatch(header)
}

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
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output file name",
		},
		cli.StringSliceFlag{
			Name:  "header, H",
			Usage: "HTTP header",
		},
	}
	return app
}

func Action(c *cli.Context) {
	app := App()
	var headers []attack.Header
	if c.String("url") == "None" || !strings.HasPrefix(c.String("url"), "http") {
		app.Run(os.Args)
		return
	}
	for _, header := range c.StringSlice("header") {
		h := headerSplit(header)
		if len(h) < 1 {
			return
		}
		headers = append(headers, attack.Header{Key: h[1], Value: h[2]})
	}
	outputFile := c.String("output")
	trunks.Trunks(c.Int("concurrency"), c.Int("requests"), c.String("url"), headers, outputFile)
}

func main() {
	app := App()
	app.Action = Action
	app.Run(os.Args)
}
