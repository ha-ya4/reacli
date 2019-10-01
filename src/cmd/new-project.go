package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func CreateNewProject() cli.Command {
	return cli.Command {
		Name: "create",
		Usage: "create new react project",
		Action: func(c *cli.Context) error {
			return action(c)
		},
	}
}

func action(c *cli.Context) error {
	fmt.Println(c.Args().Get(0))
	return nil
}