package main

import (
	"os"
	"fmt"

	"github.com/urfave/cli"

	appcmd "github.com/ha-ya4/reacli/src/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "reacli"
	app.Usage = "project cli"
	app.Version = "0.0.1"

	app.Commands = []cli.Command {
		appcmd.CreateNewProject(),
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
