package main

import (
	"testing"
)

func TestDestGenerateFilename(t *testing.T) {
	cases := []struct {
		desc     string
		filename string
		want     string
	}{
		{
			desc:     "vmファイル",
			filename: "StackArithmetic/SimpleAdd/Test.vm",
			want:     "StackArithmetic/SimpleAdd/Test.asm",
		},
		{
			desc:     "ディレクトリ（最後のスラッシュなし）",
			filename: "StackArithmetic/SimpleAdd",
			want:     "StackArithmetic/SimpleAdd/SimpleAdd.asm",
		},
		{
			desc:     "ディレクトリ（最後のスラッシュあり）",
			filename: "StackArithmetic/SimpleAdd/",
			want:     "StackArithmetic/SimpleAdd/SimpleAdd.asm",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			dest := NewDest(tc.filename)
			got := dest.generateFilename()
			if got != tc.want {
				t.Errorf("failed: got = %s, want %s", got, tc.want)
			}
		})
	}
}
