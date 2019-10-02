package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli"

	apperr "github.com/ha-ya4/reacli/src/error"
)

// Create はcreateコマンドの処理を記述した構造体をリターンする関数
func Create() cli.Command {
	return cli.Command {
		Name: "create",
		Usage: "create the specified argument",
		Flags: []cli.Flag {
			flagCreateProject(),
		},
		Action: func(c *cli.Context) error {
			return action(c)
		},
	}
}

// プロジェクト名のフラグ(create -project projectname)
func flagCreateProject() cli.StringFlag {
	return cli.StringFlag {
		Name: "project, p",
		Usage: "create new react project",
	}
}

// StringFlagの文字列に合わせて分岐する
func action(c *cli.Context) error {
	// 新しいプロジェクトを作成する
	if c.String("project") != "" {
		return createNewProject(c)
	}

	return fmt.Errorf("\n%s\n ", apperr.CreateFlagErr)
}

// create-react-appを使って新しいプロジェクトを作成する
func createNewProject(c *cli.Context) error {
	projectName := c.String("project")
	cmd := exec.Command("npx", "create-react-app", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	result := cmd.Start()

	fmt.Printf("\nstarting create new project [%s]. please wait...\n", projectName)

	cmd.Wait()
	if result != nil {
		return result
	}

	return nil
}