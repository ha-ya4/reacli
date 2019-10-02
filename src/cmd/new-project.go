package cmd

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/urfave/cli"

	apperr "github.com/ha-ya4/reacli/src/error"
)

// CreateNewProject はnpxを使いcreate-react-appで新しいreactプロジェクトを作成する
// 引数としてプロジェクト名を受け取る
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
	projectName := c.Args().Get(0)
	if projectName == "" {
		return fmt.Errorf("\n%s\n ", apperr.ProjectNameErr)
	}

	cmd := exec.Command("npx", "create-react-ap", projectName)
	result := cmd.Start()

	fmt.Printf("\nstarting create new project [%s]. please wait...\n", projectName)

	cmd.Wait()
	if result != nil {
		return result
	}

	return nil
}