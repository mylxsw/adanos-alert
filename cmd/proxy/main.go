package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli"
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"

func main() {
	app := &cli.App{
		Name:    "adanos-proxy",
		Usage:   "a proxy program for adanos-alert",
		Version: fmt.Sprintf("%s (%s)", Version, GitCommit[:8]),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "adanos-server",
				Value: "http://localhost:19999",
				Usage: "adanos-alert server addr",
			},
			&cli.StringFlag{
				Name:  "adanos-token",
				Value: "",
				Usage: "adanos-alert server token",
			},
			&cli.StringSliceFlag{
				Name:  "tag",
				Usage: "specify tags for alert message",
			},
			&cli.StringSliceFlag{
				Name:  "meta",
				Usage: "specify meta for alert message, meta must be key=value",
			},
			&cli.StringFlag{
				Name:  "origin",
				Usage: "specify origin for alert message",
			},
			&cli.IntFlag{
				Name:  "max-lines",
				Value: 100,
			},
		},
		Action: func(c *cli.Context) error {

			adanosServer := c.String("adanos-server")
			adanosToken := c.String("adanos-token")
			tags := c.StringSlice("tag")
			meta := c.StringSlice("meta")
			origin := c.String("origin")
			maxLines := c.Int("max-lines")

			fmt.Printf("server=%s, token=%s, tags=%v, meta=%v, origin=%s\n", adanosServer, adanosToken, tags, meta, origin)

			message := readStdin(maxLines)
			if message == "" {
				return nil
			}



			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func readStdin(maxLines int) string {

	result := ""

	reader := bufio.NewReader(os.Stdin)
	lineNo := 0
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		lineNo++
		if lineNo > maxLines {
			break
		}

		result += line + "\n"
	}

	return strings.TrimSpace(strings.TrimSuffix(result, "\n"))
}
