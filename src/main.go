package main

import (
	"os"
	"fmt"

	"github.com/urfave/cli"
)

/*
  create-cmd::{ commandCreate }
*/

func main() {
	app := cli.NewApp()
	app.Name = "reacli"
	app.Usage = "project cli"
	app.Version = "0.0.1"

	app.Commands = []cli.Command {
		commandCreate(),
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
