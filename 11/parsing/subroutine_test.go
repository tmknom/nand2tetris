package parsing

import (
	"../symbol"
	"../token"
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
		{
			desc: "ローカル変数のあるサブルーチンの定義",
			subroutineDec: &SubroutineDec{
				ClassName:      NewClassNameByValue("Main"),
				Subroutine:     NewKeywordByValue("function"),
				SubroutineType: NewSubroutineTypeByValue("void"),
				SubroutineName: NewSubroutineNameByValue("main"),
				ParameterList:  NewParameterList(),
				SubroutineBody: &SubroutineBody{
					VarDecs: &VarDecs{
						Items: []*VarDec{
							&VarDec{
								VarType: NewVarType(token.NewToken("int", token.TokenKeyword)),
								VarNames: &VarNames{
									First:            NewVarNameByValue("foo"),
									CommaAndVarNames: []*CommaAndVarName{},
								},
							},
							&VarDec{
								VarType: NewVarType(token.NewToken("Array", token.TokenIdentifier)),
								VarNames: &VarNames{
									First:            NewVarNameByValue("bar"),
									CommaAndVarNames: []*CommaAndVarName{},
								},
							},
						},
					},
					Statements: NewStatements(),
				},
			},
			want: []string{
				"function Main.main 2",
			},
		},
		{
			desc: "field変数のないクラスのコンストラクタの定義",
			subroutineDec: &SubroutineDec{
				ClassName:      NewClassNameByValue("Square"),
				Subroutine:     NewKeywordByValue("constructor"),
				SubroutineType: NewSubroutineTypeByValue("Square"),
				SubroutineName: NewSubroutineNameByValue("new"),
				ParameterList:  NewParameterList(),
				SubroutineBody: NewSubroutineBody(),
			},
			want: []string{
				"function Square.new 0",
				"push constant 0", // クラスのフィールド変数の数
				"call Memory.alloc 1",
				"pop pointer 0",
			},
		},
		{
			desc: "ローカル変数のないメソッドの定義: method int size()",
			subroutineDec: &SubroutineDec{
				ClassName:      NewClassNameByValue("Square"),
				Subroutine:     NewKeywordByValue("method"),
				SubroutineType: NewSubroutineTypeByValue("int"),
				SubroutineName: NewSubroutineNameByValue("size"),
				ParameterList:  NewParameterList(),
				SubroutineBody: NewSubroutineBody(),
			},
			want: []string{
				"function Square.size 0",
				"push argument 0",
				"pop pointer 0",
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

func TestSubroutineDecToCodeWithField(t *testing.T) {
	cases := []struct {
		desc          string
		subroutineDec *SubroutineDec
		want          []string
	}{
		{
			desc: "field変数のあるクラスのコンストラクタの定義",
			subroutineDec: &SubroutineDec{
				ClassName:      NewClassNameByValue("Square"),
				Subroutine:     NewKeywordByValue("constructor"),
				SubroutineType: NewSubroutineTypeByValue("Square"),
				SubroutineName: NewSubroutineNameByValue("new"),
				ParameterList:  NewParameterList(),
				SubroutineBody: NewSubroutineBody(),
			},
			want: []string{
				"function Square.new 0",
				"push constant 2", // クラスのフィールド変数の数
				"call Memory.alloc 1",
				"pop pointer 0",
			},
		},
	}

	// いろいろ初期化
	SetupTestForToCode()
	// シンボルテーブルのセットアップ
	symbol.GlobalSymbolTables.AddFieldSymbol("fieldA", "int")
	symbol.GlobalSymbolTables.AddFieldSymbol("fieldB", "int")

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.subroutineDec.ToCode()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}
