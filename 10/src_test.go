package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

const testSrcFilename = "TestSrc.jack"

func TestSrcSetup(t *testing.T) {
	cases := []struct {
		desc     string
		filename string
		length   int
		want     []string
	}{
		{
			desc:     "Setup",
			filename: "Square/SquareGame.jack",
			length:   10,
			want: []string{
				"class SquareGame {",
				"field Square square;",
				"field int direction;",
				"constructor SquareGame new() {",
				"let square = Square.new(0, 0, 30);",
				"let direction = 0;",
				"return this;",
				"}",
				"method void dispose() {",
				"do square.dispose();",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			src := NewSrc(tc.filename)
			src.Setup()

			// 全行チェックが面倒なので、lengthで指定した行数だけチェックする
			if diff := cmp.Diff(src.lines[:tc.length], tc.want); diff != "" {
				t.Errorf("failed src.lines: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestSrcDeleteCommentAndWhitespace(t *testing.T) {
	cases := []struct {
		desc      string
		line      string
		isComment bool
		want      string
	}{
		{
			desc:      "コンパイル対象の文字列",
			line:      "class Main {",
			isComment: false,
			want:      "class Main {",
		},
		{
			desc:      "コンパイル対象の文字列で一行コメントを含む",
			line:      "class Main { // コメント",
			isComment: false,
			want:      "class Main {",
		},
		{
			desc:      "コメントのみ",
			line:      "// テストコメント",
			isComment: false,
			want:      "",
		},
		{
			desc:      "空白のみ",
			line:      " ",
			isComment: false,
			want:      "",
		},
		{
			desc:      "空白＋コメント",
			line:      "  // テストコメント",
			isComment: false,
			want:      "",
		},
		{
			desc:      "空文字",
			line:      "",
			isComment: false,
			want:      "",
		},
		{
			desc:      "複数行コメントの開始（その1）",
			line:      "/*",
			isComment: false,
			want:      "",
		},
		{
			desc:      "複数行コメントの開始（その2）",
			line:      "/*",
			isComment: false,
			want:      "",
		},
		{
			desc:      "複数行コメントの開始（その3）",
			line:      "/* テスト */",
			isComment: false,
			want:      "",
		},
		{
			desc:      "複数行コメントの開始中",
			line:      " *",
			isComment: true,
			want:      "",
		},
		{
			desc:      "複数行コメントの終了",
			line:      " */",
			isComment: true,
			want:      "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			src := NewSrc(testSrcFilename)
			src.isComment = tc.isComment
			got := src.deleteCommentAndWhitespace(tc.line)

			if got != tc.want {
				t.Errorf("failed '%s': got = %s, want %s", tc.line, got, tc.want)
			}
		})
	}
}
