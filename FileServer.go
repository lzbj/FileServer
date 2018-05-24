package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		cli.Command{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start file server",

			Flags: []cli.Flag{
				cli.BoolFlag{Name: "forever, forevvarr"},
				cli.StringFlag{Name: "backend"},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("added task: ", c.FlagNames())
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
