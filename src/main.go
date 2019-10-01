package main

import (
	//"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "sample"
	app.Usage = "sample cli"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{}

	app.Run(os.Args)
}