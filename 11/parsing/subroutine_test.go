package parsing

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestSubroutineDecsToCode(t *testing.T) {
	cases := []struct {
		desc           string
		subroutineDecs *SubroutineDecs
		want           []string
	}{
		{
			desc: "サブルーチンの定義がひとつ",
			subroutineDecs: &SubroutineDecs{
				Items: []*SubroutineDec{
					&SubroutineDec{
						Subroutine:     NewKeywordByValue("function"),
						SubroutineType: NewSubroutineTypeByValue("void"),
						SubroutineName: NewSubroutineNameByValue("main"),
						ParameterList:  NewParameterList(),
						SubroutineBody: NewSubroutineBody(),
					},
				},
			},
			want: []string{
				"function main 0",
			},
		},
		{
			desc: "サブルーチンの定義が複数",
			subroutineDecs: &SubroutineDecs{
				Items: []*SubroutineDec{
					&SubroutineDec{
						Subroutine:     NewKeywordByValue("function"),
						SubroutineType: NewSubroutineTypeByValue("void"),
						SubroutineName: NewSubroutineNameByValue("foo"),
						ParameterList:  NewParameterList(),
						SubroutineBody: NewSubroutineBody(),
					},
					&SubroutineDec{
						Subroutine:     NewKeywordByValue("function"),
						SubroutineType: NewSubroutineTypeByValue("void"),
						SubroutineName: NewSubroutineNameByValue("bar"),
						ParameterList:  NewParameterList(),
						SubroutineBody: NewSubroutineBody(),
					},
					&SubroutineDec{
						Subroutine:     NewKeywordByValue("function"),
						SubroutineType: NewSubroutineTypeByValue("void"),
						SubroutineName: NewSubroutineNameByValue("baz"),
						ParameterList:  NewParameterList(),
						SubroutineBody: NewSubroutineBody(),
					},
				},
			},
			want: []string{
				"function foo 0",
				"function bar 0",
				"function baz 0",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.subroutineDecs.ToCode()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestSubroutineDecToCode(t *testing.T) {
	cases := []struct {
		desc          string
		subroutineDec *SubroutineDec
		want          []string
	}{
		{
			desc: "ローカル変数のないサブルーチンの定義",
			subroutineDec: &SubroutineDec{
				ClassName:      NewClassNameByValue("Main"),
				Subroutine:     NewKeywordByValue("function"),
				SubroutineType: NewSubroutineTypeByValue("void"),
				SubroutineName: NewSubroutineNameByValue("main"),
				ParameterList:  NewParameterList(),
				SubroutineBody: NewSubroutineBody(),
			},
			want: []string{
				"function Main.main 0",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.subroutineDec.ToCode()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}
