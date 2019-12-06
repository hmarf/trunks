package main

import (
	"os"
	"regexp"

	"github.com/hmarf/trunks/benche"
	"github.com/hmarf/trunks/benche/attack"
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
	app.Version = "0.1.2"
	app.Author = "hmarf"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url, u",
			Value: "None",
			Usage: "[required] string\n	 URL to hit",
		},
		cli.IntFlag{
			Name:  "concurrency, c",
			Value: 10,
			Usage: "int\n	 Concurrency Level.",
		},
		cli.IntFlag{
			Name:  "requests, r",
			Value: 100,
			Usage: "int\n	 Number of Requests.",
		},
		cli.StringFlag{
			Name:  "method, m",
			Value: "GET",
			Usage: "string\n	 http method.",
		},
		cli.StringSliceFlag{
			Name: "header, H",
			Usage: "string\n	 HTTP header",
		},
		cli.StringFlag{
			Name: "body, b",
			Usage: "string\n	 HTTP body",
		},
		cli.StringFlag{
			Name: "output, o",
			Usage: "string\n	 File name to output results",
		},
		cli.BoolFlag{
			Name:  "http2",
			Usage: "Send HTTP/2 requests when supported by the server (default false)",
		},
	}
	return app
}

func Action(c *cli.Context) {
	app := App()
	var headers []attack.Header
	if c.String("url") == "None" {
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
		OutputFile:  c.String("output"),
		Http2:       c.Bool("http2"),
	}
	benche.Trunks(option)
}

func main() {
	app := App()
	app.Action = Action
	app.Run(os.Args)
}
