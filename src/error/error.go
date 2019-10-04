package error

import ()

const (
	// CreateFlagErr createコマンドに何を作るかのflagが指定されていないときのエラー
	CreateFlagErr = "Missing create command flag. Please specify the create command flag."
	// CreateProjectErr プロジェクト作成に失敗したときのエラー
	CreateProjectErr = "Failed to create project."
	// CreateComponentErr コンポーネントの作成に失敗したときのエラー
	CreateComponentErr = "Failed to create component."
)