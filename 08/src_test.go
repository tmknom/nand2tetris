package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

const testSrcFilename = "TestSrc.vm"

func TestSrcSetup(t *testing.T) {
	cases := []struct {
		desc       string
		filename   string
		lines      []string
		moduleName string
	}{
		{
			desc:     "Setup",
			filename: "StackArithmetic/SimpleAdd/SimpleAdd.vm",
			lines: []string{
				"push constant 7",
				"push constant 8",
				"add",
			},
			moduleName: "SimpleAdd",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			src := NewSrc(tc.filename)
			src.Setup()

			if diff := cmp.Diff(src.lines, tc.lines); diff != "" {
				t.Errorf("failed src.lines: diff (-got +want):\n%s", diff)
			}

			if src.moduleName != tc.moduleName {
				t.Errorf("failed src.moduleName: got = %s, want %s", src.moduleName, tc.moduleName)
			}
		})
	}
}

func TestSrcSetupLines(t *testing.T) {
	cases := []struct {
		desc string
		org  []string
		want []string
	}{
		{
			desc: "コメントと空白を除外した文字列のリストを取得",
			org: []string{
				"// テストコード",
				"",
				"push constant 7",
				"push constant 5 // コメント",
				"  // コメント",
			},
			want: []string{
				"push constant 7",
				"push constant 5",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			src := NewSrc(testSrcFilename)
			src.org = tc.org
			src.setupLines()

			if diff := cmp.Diff(src.lines, tc.want); diff != "" {
				t.Errorf("failed src.lines: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestSrcDeleteCommentAndWhitespace(t *testing.T) {
	cases := []struct {
		desc string
		line string
		want string
	}{
		{
			desc: "パース対象の文字列",
			line: "push constant 7",
			want: "push constant 7",
		},
		{
			desc: "パース対象の文字列でコメントを含む",
			line: "push constant 7 // コメント",
			want: "push constant 7",
		},
		{
			desc: "コメントのみ",
			line: "// テストコメント",
			want: "",
		},
		{
			desc: "空白のみ",
			line: " ",
			want: "",
		},
		{
			desc: "空白＋コメント",
			line: "  // テストコメント",
			want: "",
		},
		{
			desc: "空文字",
			line: "",
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			src := NewSrc(testSrcFilename)
			got := src.deleteCommentAndWhitespace(tc.line)

			if got != tc.want {
				t.Errorf("failed '%s': got = %s, want %s", tc.line, got, tc.want)
			}
		})
	}
}

func TestSrcSetupModuleName(t *testing.T) {
	cases := []struct {
		desc     string
		filename string
		want     string
	}{
		{
			desc:     "setupModuleName",
			filename: "StackArithmetic/SimpleAdd/SimpleAdd.vm",
			want:     "SimpleAdd",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			src := NewSrc(tc.filename)
			src.setupModuleName()

			if src.moduleName != tc.want {
				t.Errorf("failed src.moduleName: got = %s, want %s", src.moduleName, tc.want)
			}
		})
	}
}
