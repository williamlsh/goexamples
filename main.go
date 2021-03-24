package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "ginger-crouton",
				Usage: "is it in the soup?",
			},
		},
		Action: func(c *cli.Context) error {
			if !c.Bool("ginger-crouton") {
				return cli.Exit("Ginger croutons are not in the soup", 86)
			}
			return nil
		},
	}

	log.Fatal(app.Run(os.Args))
}
