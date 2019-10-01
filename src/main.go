package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/ha-ya4/reacli/src/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "sample"
	app.Usage = "sample cli"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{}
	fmt.Println("hello")
	//cmd.Hello()

	app.Run(os.Args)
}