package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestNewArg(t *testing.T) {
	cases := []struct {
		desc string
		args []string
		want []string
	}{
		{
			desc: "jackファイルの指定",
			args: []string{"dummy", "foo.jack"},
			want: []string{"foo.jack"},
		},
		{
			desc: "ディレクトリの指定",
			args: []string{"dummy", "Fixture/SquareVersion10/"},
			want: []string{
				"Fixture/SquareVersion10/Main.jack",
				"Fixture/SquareVersion10/Square.jack",
				"Fixture/SquareVersion10/SquareGame.jack",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			arg := NewArg(tc.args)
			if diff := cmp.Diff(arg.files, tc.want); diff != "" {
				t.Errorf("failed arg.files: diff (-got +want):\n%s", diff)
			}
		})
	}
}
