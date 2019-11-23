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
			Name:  "method, m",
			Value: "Get",
			Usage: "http method",
		},
		cli.StringFlag{
			Name:  "url, u",
			Value: "None",
			Usage: "URL to hit",
		},
		cli.StringSliceFlag{
			Name:  "header, H",
			Usage: "HTTP header",
		},
		cli.StringFlag{
			Name:  "body, b",
			Usage: "HTTP body",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output file name",
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
	option := attack.Option{
		Concurrency: c.Int("concurrency"),
		Requests:    c.Int("requests"),
		Method:      c.String("method"),
		URL:         c.String("url"),
		Header:      headers,
		Body:        c.String("body"),
		OutputFile:  c.String("output")}
	trunks.Trunks(option)
}

func main() {
	app := App()
	app.Action = Action
	app.Run(os.Args)
}
