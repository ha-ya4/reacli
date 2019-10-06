package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"

	"github.com/urfave/cli"
)

/*
	error::{ createFlagErr, createProjectErr, createComponentErr }
	file-content::{ componentContent, tsComponentContent, testContent }
	utils::{ createEmbeddedFile, execCommand }
*/


// Create はcreateコマンドの処理を記述した構造体をリターンする関数
func commandCreate() cli.Command {

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

type cliContexter interface {
	String(name string) string
	Bool(name string) bool
}

// StringFlagの文字列に合わせて分岐する
func createAction(c cliContexter) error {

	// 新しいプロジェクトを作成する
	if c.String("project") != "" {
		project := newProject(c.String("project"), c.Bool("default"), c.Bool("typescript"))
		return project.create()
	}

	// 新しいコンポーネントファイルを作成する
	if c.String("component") != "" {
		component := newComponent(c.String("component"), c.Bool("dir"), c.Bool("typescript"))
		return component.create(c)
	}

	return fmt.Errorf("\nreacli ERR: %s\n ", createFlagErr)
}


type project struct {
	name string
	flagDefault bool
	flagTS bool
}

func newProject(n string, d, ts bool) project {
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
		return fmt.Errorf("\nreacli ERR: %s\n ", createProjectErr)
	}

	// Appファイルをクラスコンポーネントに書き換えてrender()の中身を消す
	if project.flagTS == true {
		appTSX := strings.Replace(tsComponentContent, "{$1}", "App", 3)
	  err = ioutil.WriteFile("App.tsx", []byte(appTSX), 0777)
	} else {
		appJS := strings.Replace(componentContent, "{$1}", "App", 3)
	  err = ioutil.WriteFile("App.js", []byte(appJS), 0777)
	}

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
	flagTS bool
}

type jsFile struct {
	name string
	flagTS bool
}

type cssFile struct {
	name string
}

type testFile struct {
	name string
	flagTS bool
}

func newComponent(n string, d, ts bool) component {
	return component {
		name: n,
		flagDir: d,
		flagTS: ts,
	}
}

func newJSFile(n string, c cliContexter) jsFile {
	return jsFile {
		name: n,
		flagTS: c.Bool("typescript"),
	}
}

func newCSSFile(n string) cssFile {
	return cssFile {
		name: n,
	}
}

func newTestFile(n string, c cliContexter) testFile {
	return testFile {
		name: n,
		flagTS: c.Bool("typescript"),
	}
}

// カレントディレクトリに新しいコンポーネント.js、コンポーネント.css、テストファイルを作る
func(component component) create(c cliContexter) (err error) {

	// dirフラグがあれば新しいディレクトリを作成しその中にコンポーネントを作成する
	// エラーがでたらファイル作成失敗としてリターンする
	err = component.dirFlag()
	if err != nil {
		return fmt.Errorf("\nreacli ERR: %s\n ", createComponentErr + "component.")
	}

	js := newJSFile(component.name, c)
	err = js.create()

	// tsフラグのありなしで拡張子とファイルに書き込む内容を変える
	extension := component.selectExtension()

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

func(component component) selectExtension() string {

	if component.flagTS == true {
		return ".tsx"
	}
	return ".js"
}

func(js jsFile) create() (err error) {

	componentName := js.name + js.selectExtension()
	content := js.selectContent()

	// コンポーネント名を埋め込んだJSファイル作成
	createErr := createEmbeddedFile(componentName, func() string {
		return strings.Replace(content, "{$1}", js.name, 3)
	})

	if createErr != nil {
		return fmt.Errorf("\nreacli ERR: %s\n ", createComponentErr + "js file.")
	}
	return
}

// jsファイルに書き込む内容を選択する
func(js jsFile) selectContent() string {

	if js.flagTS == true {
		return tsComponentContent
	}

	return componentContent
}

// jsファイルにつける拡張子を選択する
func(js jsFile) selectExtension() string {

	if js.flagTS == true {
		return ".tsx"
	}

	return ".js"
}

