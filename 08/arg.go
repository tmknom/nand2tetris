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

const DefaultArg = "FunctionCalls/FibonacciElement/"

func NewArg(args []string) *Arg {
	arg := DefaultArg
	if len(args) >= 2 {
		arg = args[1]
	}

	if filepath.Ext(arg) == ".vm" {
		return &Arg{raw: arg, files: []string{arg}}
	}

	// vmファイルを指定していない場合は、ディレクトリが指定されたとみなす
	// ただしファイル名にTestが含まれる場合は除外する
	files, _ := filepath.Glob(arg + "/*.vm")
	ignoreTestFiles := []string{}
	for _, file := range files {
		if !strings.Contains(file, "Test") {
			ignoreTestFiles = append(ignoreTestFiles, file)
		}
	}
	return &Arg{raw: arg, files: ignoreTestFiles}
}
