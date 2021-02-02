package main

import (
	"path/filepath"
	"strings"
)

// コマンドの入力パラメータをパースして、変換対象のvmファイル名を管理
type Arg struct {
	raw   string
	files []string
}

const DefaultArg = "Fixture/Seven/"

func NewArg(args []string) *Arg {
	arg := DefaultArg
	if len(args) >= 2 {
		arg = args[1]
	}

	if filepath.Ext(arg) == ".jack" {
		return &Arg{raw: arg, files: []string{arg}}
	}

	// jackファイルを指定していない場合は、ディレクトリが指定されたとみなす
	// ただしファイル名にIgnoreが含まれる場合は除外する
	files, _ := filepath.Glob(arg + "/*.jack")
	ignoreTestFiles := []string{}
	for _, file := range files {
		if !strings.Contains(file, "Ignore") {
			ignoreTestFiles = append(ignoreTestFiles, file)
		}
	}
	return &Arg{raw: arg, files: ignoreTestFiles}
}
