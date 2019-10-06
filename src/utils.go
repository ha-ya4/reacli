package main

import (
	"os"
	"os/exec"
)

// exec.Commandを実行し、stdout, stderrを出力する。
// cutin引数に関数を入れることによって、実行中の待機時間に処理を挟むことができる。
// はさみたい処理が無い場合、呼び出し側ではcutin引数にfunc(){}を指定することになってしまうので
// cutinを可変長引数にし、何も無い場合はosExec(command, args)のみで呼び出せるようにした
func execCommand(command string, args []string, cutin ...func()) error {

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	result := cmd.Start()

	// cutinのlengtが０なら関数を受け取っていないことになる
	// 0じゃなければrangeで可変引数をばらして関数を使う
	if len(cutin) != 0 {
		for _, function := range cutin {
			function()
		}
	}

	cmd.Wait()
	if result != nil {
		return result
	}
	return nil
}

// 必要な情報を埋め込んだファイルを作成する
// 埋め込む内容は外から関数として受け取る
func createEmbeddedFile(name string, replace func() string) (err error) {
	file, err := os.Create(name)
	embedded := replace()
	_, err = file.Write([]byte(embedded))
	file.Close()
	return
}