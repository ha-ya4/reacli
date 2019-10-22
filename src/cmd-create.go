package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/urfave/cli"
)

/*
	error::{ createFlagErr, createProjectErr, createComponentErr }
	file-content::{ componentContent, tsComponentContent, tsSfcComponentContent, sfcComponentContent, testContent }
	utils::{ createEmbeddedFile, execCommand }
*/

// Create はcreateコマンドの処理を記述した構造体をリターンする関数
func commandCreate() cli.Command {

	return cli.Command{
		Name:  "create",
		Usage: "create the specified argument",
		Flags: createFlag(),
		Action: func(c *cli.Context) error {
			return createAction(c)
		},
	}
}

func createFlag() []cli.Flag {

	return []cli.Flag{
		// プロジェクト作成
		// ディレクトリを２つ追加、logo.svg削除、App.jsをkclassに書き換える
		cli.StringFlag{
			Name:  "project, p",
			Usage: "create new react project and if you need setup",
		},
		// プロジェクトをcreate-react-appの初期状態のままにしておくかのフラグ
		cli.BoolFlag{
			Name:  "default, d",
			Usage: "default project",
		},
		// TSを使用するかのフラグ
		cli.BoolFlag{
			Name:  "typescript, ts",
			Usage: "create new react project and if you need setup with typescript",
		},
		// SCSSを使用するかのフラグ
		cli.BoolFlag{
			Name:  "scss",
			Usage: "create new react project and if you need setup with scss",
		},
		// コンポーネントファイル作成
		cli.StringFlag{
			Name:  "component, c",
			Usage: "create new component file",
		},
		// 新しいディレクトリを作成しその中にコンポーネントを作成する
		cli.BoolFlag{
			Name:  "dir",
			Usage: "Create a new directory and create components in it",
		},
		cli.BoolFlag{
			Name:  "sfc",
			Usage: "create new SFC component file",
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
		project := newProject(
			c.String("project"),
			c.Bool("default"),
			c.Bool("typescript"),
			c.Bool("scss"),
		)
		return project.create()
	}

	// 新しいコンポーネントファイルを作成する
	if c.String("component") != "" {
		component := newComponent(
			c.String("component"),
			c.Bool("dir"),
			c.Bool("typescript"),
			c.Bool("scss"),
			c.Bool("sfc"),
		)
		return component.create(c)
	}

	return fmt.Errorf("\nreacli ERR: %s\n ", createFlagErr)
}

type project struct {
	name        string
	flagDefault bool
	flagTS      bool
	flagSCSS    bool
}

func newProject(n string, d, ts, scss bool) project {
	return project{
		name:        n,
		flagDefault: d,
		flagTS:      ts,
		flagSCSS:    scss,
	}
}

// create-react-appを使って新しいプロジェクトを作成する
func (project project) create() error {

	args := project.createArgs()
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

// プロジェクト作成に使うコマンドの引数を作成する(create-react-app --typescript)
func (project project) createArgs() []string {
	args := []string{"create-react-app", project.name}
	// tsフラグがあればtypescriptを導入
	if project.flagTS == true {
		args = append(args, "--typescript")
	}
	return args
}

// デフォルトのプロジェクトを変更
func (project project) setUp() (err error) {

	// srcフォルダに移動。移動ができなければプロジェクト作成失敗のはずなのでエラーメッセージを出す
	err = os.Chdir(project.name + "/src")
	if err != nil {
		return fmt.Errorf("\nreacli ERR: %s\n ", createProjectErr)
	}

	err = project.setUpJSFile()
	err = project.setUpCSSFile()

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

// Appファイルをクラスコンポーネントに書き換えてrender()の中身を消す
// -scss jsファイルのcssインポートの部分部分をscssに変更
func (project project) setUpJSFile() error {

	// -tsフラグがあればtsの内容にする
	if project.flagTS == true {
		appTSX := strings.Replace(tsComponentContent, "{$1}", "App", 3)
		// scssフラグがあればcssインポート部分をscssに変更
		if project.flagSCSS == true {
			appTSX = strings.Replace(appTSX, ".css", ".scss", 1)
		}
		return ioutil.WriteFile("App.tsx", []byte(appTSX), 0777)
	}

	// -tsフラグがなければjsの内容にする
	appJS := strings.Replace(componentContent, "{$1}", "App", 3)
	// scssフラグがあればcssインポート部分をscssに変更
	if project.flagSCSS == true {
		appJS = strings.Replace(appJS, ".css", ".scss", 1)
	}
	return ioutil.WriteFile("App.js", []byte(appJS), 0777)
}

// scssフラグのありなしでcssファイルの拡張子を変更
func (project project) setUpCSSFile() (err error) {

	// scssフラグがなければ何もせずにリターン
	if project.flagSCSS != true {
		return nil
	}

	// scssフラグがあればAppとindexのcssをscssに変更
	appCSS := []string{"App.css", "App.scss"}
	err = execCommand("mv", appCSS)
	indexCSS := []string{"index.css", "index.scss"}
	err = execCommand("mv", indexCSS)
	return
}

type component struct {
	name     string
	flagDir  bool
	flagTS   bool
	flagSCSS bool
	flagSFC  bool
}

type jsFile struct {
	name     string
	flagTS   bool
	flagSCSS bool
	flagSFC  bool
}

type cssFile struct {
	name     string
	flagSCSS bool
}

type testFile struct {
	name   string
	flagTS bool
}

func newComponent(n string, d, ts, scss, sfc bool) component {

	return component{
		name:     n,
		flagDir:  d,
		flagTS:   ts,
		flagSCSS: scss,
		flagSFC:  sfc,
	}
}

func newJSFile(n string, c cliContexter) jsFile {

	return jsFile{
		name:     n,
		flagTS:   c.Bool("typescript"),
		flagSCSS: c.Bool("scss"),
		flagSFC:  c.Bool("sfc"),
	}
}

func newCSSFile(n string, c cliContexter) cssFile {

	return cssFile{
		name:     n,
		flagSCSS: c.Bool("scss"),
	}
}

func newTestFile(n string, c cliContexter) testFile {

	return testFile{
		name:   n,
		flagTS: c.Bool("typescript"),
	}
}

type componentFile interface {
	create() error
}

// カレントディレクトリに新しいコンポーネント.js、コンポーネント.css、テストファイルを作る
func (component component) create(c cliContexter) (err error) {

	// dirフラグがあれば新しいディレクトリを作成しその中にコンポーネントを作成する
	// エラーがでたらファイル作成失敗としてリターンする
	err = component.dirFlag()
	if err != nil {
		return fmt.Errorf("\nreacli ERR: %s\n ", createComponentErr+"component.")
	}

	// ファイル作成に必要な構造体作成
	js := newJSFile(component.name, c)
	css := newCSSFile(component.name, c)
	test := newTestFile(component.name, c)
	// forでまとめてcreateするために配列に入れる
	files := [3]componentFile{js, css, test}
	// ファイル作成時のエラーを全て拾いたいので、エラーを文字列として結合するための変数
	var createErr string
	// まとめてファイル作成。エラーハンドリングをまとめてするためにforを使う
	// エラーの場合は文字列として結合し、あとでerrors.Newする
	// これで階層化せずにエラーを全て表示できる
	for _, f := range files {
		e := f.create()
		if e != nil {
			createErr += e.Error()
		}
	}

	// ファイル作成時にエラーが出ていればerrorを生成しreturnする
	if createErr != "" {
		err = errors.New(createErr)
		return
	}

	fmt.Printf("\ncreate a new component [%s] all ready exists\n ", component.name)
	return
}

// dirフラグがあれば新しいディレクトリを作成しcdする
func (component component) dirFlag() (err error) {

	if component.flagDir == true {
		err = os.Mkdir(component.name, 0777)
		err = os.Chdir(component.name)
	}
	return
}

// jsファイル作成
func (js jsFile) create() (err error) {

	fileName := js.name + selectJSExtension(js.flagTS)
	content := js.selectContent()
	cssExtention := selectCSSExtension(js.flagSCSS)

	// コンポーネント名を埋め込んだJSファイル作成
	createErr := createEmbeddedFile(fileName, func() string {
		c := strings.Replace(content, "{$1}", js.name, 3)
		return strings.Replace(c, ".css", cssExtention, 1)
	})

	if createErr != nil {
		return fmt.Errorf("\nreacli ERR: %s\n ", createComponentErr+" js file.")
	}
	return
}

// jsファイルに書き込む内容を選択する
func (js jsFile) selectContent() string {

	if js.flagSFC == true && js.flagTS == true {
		return tsSfcComponentContent
	}

	if js.flagSFC == true {
		return sfcComponentContent
	}

	if js.flagTS == true {
		return tsComponentContent
	}

	return componentContent
}

// cssファイル作成
func (css cssFile) create() (err error) {

	fileName := css.name + selectCSSExtension(css.flagSCSS)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("\nreacli ERR: %s\n ", createComponentErr+" css file.")
	}
	file.Close()
	return
}

// testファイル作成
func (test testFile) create() (err error) {

	fileName := test.name + ".test" + selectJSExtension(test.flagTS)

	// コンポーネント名を埋め込んだtestファイル作成
	createErr := createEmbeddedFile(fileName, func() string {
		return strings.Replace(testContent, "{$1}", test.name, 3)
	})

	if createErr != nil {
		return fmt.Errorf("\nreacli ERR: %s\n ", createComponentErr+" test file.")
	}
	return
}

// jsファイルにつける拡張子を選択する
func selectJSExtension(ts bool) string {

	if ts == true {
		return ".tsx"
	}

	return ".js"
}

// cssファイルにつける拡張子を選択する
func selectCSSExtension(scss bool) string {

	if scss == true {
		return ".scss"
	}
	return ".css"
}
