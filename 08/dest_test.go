package main

import (
	"testing"
)

func TestNewDest(t *testing.T) {
	cases := []struct {
		desc     string
		filename string
		want     string
	}{
		{
			desc:     "asmファイルに変換",
			filename: "StackArithmetic/SimpleAdd/Test.vm",
			want:     "StackArithmetic/SimpleAdd/Test.asm",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			dest := NewDest(tc.filename)
			got := dest.filename
			if got != tc.want {
				t.Errorf("failed: got = %s, want %s", got, tc.want)
			}
		})
	}
}
