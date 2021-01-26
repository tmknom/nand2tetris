package main

import (
	"bufio"
	"os"
	"strings"
)

// コンパイル対象のソースファイルを読み込む
type Src struct {
	filename  string
	org       []string
	lines     []string
	isComment bool
}

func NewSrc(filename string) *Src {
	return &Src{filename: filename, isComment: false}
}

func (s *Src) Setup() error {
	// ファイルを読み込んで、一度全部メモリに展開する
	err := s.readFile(s.filename)
	if err != nil {
		return err
	}

	// コメントと空白を除外
	s.setupLines()
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
	}
}

func (s *Src) deleteCommentAndWhitespace(line string) string {
	// 複数行コメントの開始文字列を見つけたらフラグを立てる
	if strings.Contains(line, "/*") {
		s.isComment = true
	}

	// 複数行コメントの終了文字列を見つけたらフラグを反転して空文字を返す
	if s.isComment && strings.Contains(line, "*/") {
		s.isComment = false
		return ""
	}

	// 複数行コメントの途中なら空文字を返す
	if s.isComment {
		return ""
	}

	// 一行コメントを除外
	deletedComment := line
	if strings.Contains(line, "//") {
		deletedComment = line[:strings.Index(line, "//")]
	}

	// 空白を除去
	return strings.TrimSpace(deletedComment)
}
