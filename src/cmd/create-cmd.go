package cmd

import (
	"fmt"
	"os"
	"io/ioutil"

	"github.com/urfave/cli"

	apperr "github.com/ha-ya4/reacli/src/error"
)


// Create はcreateコマンドの処理を記述した構造体をリターンする関数
func Create() cli.Command {

	return cli.Command {
		Name: "create",
		Usage: "create the specified argument",
		Flags: flag(),
		Action: func(c *cli.Context) error {
			return action(c)
		},
	}
}

func flag() []cli.Flag {
	return []cli.Flag {
		// プロジェクト作成
		// ディレクトリを２つ追加、logo.svg削除、App.jsをkclassに書き換える
		cli.StringFlag {
			Name: "project, p",
			Usage: "create new react project and if you need setup",
		},
		// プロジェクトをcreate-react-appの初期状態のままにしておくかのフラグ
		cli.BoolFlag {
			Name: "default, d",
			Usage: "default project",
		},
		// コンポーネントファイル作成
		cli.StringFlag {
			Name: "component, c",
			Usage: "create new component file",
		},
		// 新しいディレクトリを作成しその中にコンポーネントを作成する
		cli.BoolFlag {
			Name: "dir",
			Usage: "Create a new directory and create components in it",
		},
	}
}

// StringFlagの文字列に合わせて分岐する
func action(c *cli.Context) error {

	// 新しいプロジェクトを作成する
	if c.String("project") != "" {
		return createNewProject(c)
	}

	// 新しいコンポーネントファイルを作成する
	if c.String("component") != "" {
		return createNewComponent(c)
	}

	return fmt.Errorf("\nreacli ERR: %s\n ", apperr.CreateFlagErr)
}


// create-react-appを使って新しいプロジェクトを作成する
func createNewProject(c *cli.Context) error {

	projectName := c.String("project")
	args := []string { "create-react-app", projectName}
	result := execCommand("npx", args, func() {
		fmt.Printf("\nstarting create a new project [%s]. please wait...\n", projectName)
	})

	if result != nil {
		return result
	}

	// defaultフラグがなければデフォルトプロジェクトから変更
	if c.Bool("default") == false {
	  return projectSetUp(c)
	}

	return nil
}

// デフォルトのプロジェクトを変更
func projectSetUp(c *cli.Context) (err error) {

	// srcフォルダに移動。移動ができなければプロジェクト作成失敗のはずなのでエラーメッセージを出す
	err = os.Chdir(c.String("project") + "/src")
	if err != nil {
		return fmt.Errorf("\nreacli ERR: %s\n ", apperr.CreateProjectErr)
	}

	// App.jsをクラスコンポーネントに書き換えてrender()の中身をdivのみにする
	err = ioutil.WriteFile("App.js", []byte(appComponentContent), 0777)
	// componentsフォルダ作成
	err = os.Mkdir("components", 0777)
	// viewsフォルダ作成
	err = os.Mkdir("views", 0777)
	// 最初から用意されてるreactのロゴを削除
	err = os.Remove("logo.svg")

	if err == nil {
		fmt.Println("\nproject setup OK!\n ")
	}
	return
}

// カレントディレクトリに新しいコンポーネント.js、コンポーネント.css、テストファイルを作る
func createNewComponent(c *cli.Context) (err error) {

	componentName := c.String("component")

	// dirフラグがあれば新しいディレクトリを作成しその中にコンポーネントを作成する
	if c.Bool("dir") == true {
		err = os.Mkdir(componentName, 0777)
		err = os.Chdir(componentName)
		if err != nil {
			return fmt.Errorf("\nreacli ERR: %s\n ", apperr.CreateComponentErr)
		}
	}

	jsFile, err := os.Create(componentName + ".js")
	_, err = jsFile.Write([]byte(componentContent))
	jsFile.Close()

	cssFile, err := os.Create(componentName + ".css")
	cssFile.Close()

	testFile, err := os.Create(componentName + ".test.js")
	_, err = testFile.Write([]byte(testContent))
	testFile.Close()

	if err == nil {
		fmt.Printf("\ncreate a new [component] %s all ready exists\n ", componentName)
	}
	return
}