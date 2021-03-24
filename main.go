package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "greet",
		Usage: "fight the lonliness!",
		Action: func(c *cli.Context) error {
			fmt.Printf("arg: %q\n", c.Args().Get(0))
			fmt.Println("Hello friend!")
			return nil
		},
	}

	log.Fatal(app.Run(os.Args))
}
