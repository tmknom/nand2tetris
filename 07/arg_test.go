package main

import (
	"reflect"
	"testing"
)

func TestNewArg(t *testing.T) {
	cases := []struct {
		desc string
		args []string
		want []string
	}{
		{
			desc: "vmファイルの指定",
			args: []string{"dummy", "foo.vm"},
			want: []string{"foo.vm"},
		},
		{
			desc: "ディレクトリの指定",
			args: []string{"dummy", "StackArithmetic/SimpleAdd/"},
			want: []string{"StackArithmetic/SimpleAdd/SimpleAdd.vm"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			arg := NewArg(tc.args)
			if !reflect.DeepEqual(arg.files, tc.want) {
				t.Errorf("failed: got = %v, want %v", arg.files, tc.want)
			}
		})
	}
}
