package main

import (
	"path/filepath"
)

// コマンドの入力パラメータをパースして、変換対象のvmファイル名を管理
type Arg struct {
	raw   string
	files []string
}

const DefaultArg = "ProgramFlow/BasicLoop/"

func NewArg(args []string) *Arg {
	arg := DefaultArg
	if len(args) >= 2 {
		arg = args[1]
	}

	if filepath.Ext(arg) == ".vm" {
		return &Arg{raw: arg, files: []string{arg}}
	}

	// vmファイルを指定していない場合は、ディレクトリが指定されたとみなす
	files, _ := filepath.Glob(arg + "/*.vm")
	return &Arg{raw: arg, files: files}
}
