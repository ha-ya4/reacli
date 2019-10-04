package cmd

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"

	"github.com/urfave/cli"

	apperr "github.com/ha-ya4/reacli/src/error"
)


// Create はcreateコマンドの処理を記述した構造体をリターンする関数
func Create() cli.Command {

	return cli.Command {
		Name: "create",
		Usage: "create the specified argument",
		Flags: createFlag(),
		Action: func(c *cli.Context) error {
			return createAction(c)
		},
	}
}

func createFlag() []cli.Flag {

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
		cli.BoolFlag {
			Name: "typescript, ts",
			Usage: "create new react project and if you need setup with typescript",
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
func createAction(c cliContexter) error {

	// 新しいプロジェクトを作成する
	if c.String("project") != "" {
		p := newProject(c.String("project"), c.Bool("default"), c.Bool("typescript"))
		return p.create()
	}

	// 新しいコンポーネントファイルを作成する
	if c.String("component") != "" {
		c := newComponent(c.String("component"), c.Bool("dir"))
		return c.create()
	}

	return fmt.Errorf("\nreacli ERR: %s\n ", apperr.CreateFlagErr)
}


type project struct {
	name string
	flagDefault bool
	flagTS bool
}

func newProject(n string, d bool, ts bool) project {
	return project {
		name: n,
		flagDefault: d,
		flagTS: ts,
	}
}

// create-react-appを使って新しいプロジェクトを作成する
func(project project) create() error {

	args := []string { "create-react-app", project.name }
	// tsフラグがあればtypescriptを導入
	if project.flagTS == true {
		args = append(args, "--typescript")
	}
	// create-react-app実行
	result := execCommand("npx", args, func() {
		fmt.Printf("\nstarting create a new project [%s]. please wait...\n", project.name)
	})

	if result != nil {
		return result
	}

	// defaultフラグがなければデフォルトプロジェクトから変更
	if project.flagDefault == false {
	  return project.setUp()
	}

	return nil
}

// デフォルトのプロジェクトを変更
func(project project) setUp() (err error) {

	// srcフォルダに移動。移動ができなければプロジェクト作成失敗のはずなのでエラーメッセージを出す
	err = os.Chdir(project.name + "/src")
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


type component struct {
	name string
	flagDir bool
}

func newComponent(n string, d bool) component {
	return component {
		name: n,
	  flagDir: d,
	}
}

// カレントディレクトリに新しいコンポーネント.js、コンポーネント.css、テストファイルを作る
func(component component) create() (err error) {

	// dirフラグがあれば新しいディレクトリを作成しその中にコンポーネントを作成する
	// エラーがでたらファイル作成失敗としてリターンする
	err = component.dirFlag()
	if err != nil {
		return fmt.Errorf("\nreacli ERR: %s\n ", apperr.CreateComponentErr)
	}

	extension := ".js"
	// コンポーネント名を埋め込んだJSファイル作成
	err = createEmbeddedFile(component.name + extension, func() string {
		return strings.Replace(componentContent, "{$1}", component.name, 3)
	})
	// コンポーネント名を埋め込んだテストファイル作成
	err = createEmbeddedFile(component.name + ".test" + extension, func() string {
		return strings.Replace(testContent, "{$1}", component.name, 2)
	})
	//　cssファイル作成
	cssFile, err := os.Create(component.name + ".css")
	cssFile.Close()

	// エラーがなければファイル作成したことを伝えるメッセージを出力する
	if err == nil {
		fmt.Printf("\ncreate a new component [%s] all ready exists\n ", component.name)
	}
	return
}

// dirフラグがあれば新しいディレクトリを作成しcdする
func(component component) dirFlag() (err error) {

	if component.flagDir == true {
		err = os.Mkdir(component.name, 0777)
		err = os.Chdir(component.name)
	}
	return
}