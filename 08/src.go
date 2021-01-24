package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// パース対象のソースファイルを読み込む
type Src struct {
	filename   string
	org        []string
	lines      []string
	moduleName string
}

func NewSrc(filename string) *Src {
	return &Src{filename: filename}
}

func (s *Src) Setup() error {
	// ファイルを読み込んで、一度全部メモリに展開する
	err := s.readFile(s.filename)
	if err != nil {
		return err
	}

	// コメントと空白を除外
	s.setupLines()
	// モジュール名をセット
	s.setupModuleName()
	return nil
}

func (s *Src) readFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		s.org = append(s.org, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (s *Src) setupLines() {
	for _, line := range s.org {
		withoutComment := s.deleteCommentAndWhitespace(line)
		if withoutComment != "" {
			s.lines = append(s.lines, withoutComment)
		}

		//// 以前のテストケースも動くように「function Sys.init」が定義されてるときだけSys.initを呼ぶようフラグをセット
		//if strings.Contains(line, "function Sys.init") {
		//	s.lines = append([]string{"call Sys.init 0"}, s.lines...)
		//}
	}
}

func (s *Src) deleteCommentAndWhitespace(line string) string {
	// コメントを除外
	deletedComment := line
	if strings.Contains(line, "//") {
		deletedComment = line[:strings.Index(line, "//")]
	}

	// 空白を除去
	return strings.TrimSpace(deletedComment)
}

func (s *Src) setupModuleName() {
	s.moduleName = filepath.Base(s.filename[:len(s.filename)-len(filepath.Ext(s.filename))])
}
