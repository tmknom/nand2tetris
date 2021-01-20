package main

import (
	"reflect"
	"testing"
)

func TestNewSrc(t *testing.T) {
	cases := []struct {
		desc string
		args []string
		want string
	}{
		{
			desc: "引数指定なし",
			args: []string{"dummy"},
			want: DefaultArg,
		},
		{
			desc: "引数指定あり",
			args: []string{"dummy", "foo.vm"},
			want: "foo.vm",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			src := NewSrc(tc.args)
			got := src.arg
			if got != tc.want {
				t.Errorf("failed: got = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestSrcParse(t *testing.T) {
	cases := []struct {
		desc string
		arg  string
		want []string
	}{
		{
			desc: "vmファイル",
			arg:  "foo.vm",
			want: []string{"foo.vm"},
		},
		{
			desc: "vmファイル以外",
			arg:  "foo.go",
			want: []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			src := NewSrc([]string{"dummy", tc.arg})
			src.Parse()

			if !reflect.DeepEqual(src.files, tc.want) {
				t.Errorf("failed: got = %v, want %v", src.files, tc.want)
			}
		})
	}
}
