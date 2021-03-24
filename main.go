package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	var test int
	flags := []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "test",
			Destination: &test,
		}),
		&cli.StringFlag{Name: "load"},
	}
	app := &cli.App{
		Flags: flags,
		Action: func(c *cli.Context) error {
			fmt.Println("yaml ist rad")
			fmt.Println("test:", test)
			return nil
		},
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("load")),
	}

	log.Fatal(app.Run(os.Args))
}
