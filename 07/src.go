package main

import (
	"path/filepath"
)

// コマンドの入力パラメータをパースして、変換対象のvmファイル名を管理
type Src struct {
	arg   string
	files []string
}

const DefaultArg = "StackArithmetic/StackTest/StackTest.vm"

func NewSrc(args []string) *Src {
	arg := DefaultArg
	if len(args) >= 2 {
		arg = args[1]
	}
	return &Src{arg: arg, files: []string{}}
}

func (s *Src) Parse() {
	if filepath.Ext(s.arg) == ".vm" {
		s.files = append(s.files, s.arg)
		return
	}

	// TODO ディレクトリの読み取りはまだ未実装だよ
}
