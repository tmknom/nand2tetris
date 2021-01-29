package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParserParseClass(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*Token
		want   *Class
	}{
		{
			desc: "ClassVarDecsとSubroutineDecが初期値のクラス",
			tokens: []*Token{
				NewToken("class", TokenKeyword),
				NewToken("Main", TokenIdentifier),
				NewToken("{", TokenSymbol),
				NewToken("}", TokenSymbol),
			},
			want: &Class{
				Keyword:        NewKeywordByValue("class"),
				ClassName:      NewClassName(NewToken("Main", TokenIdentifier)),
				OpenSymbol:     ConstOpeningCurlyBracket,
				CloseSymbol:    ConstClosingCurlyBracket,
				ClassVarDecs:   NewClassVarDecs(),
				SubroutineDecs: NewSubroutineDecs(),
			},
		},
		{
			desc: "SubroutineDecsが初期値のクラス",
			tokens: []*Token{
				NewToken("class", TokenKeyword),
				NewToken("Main", TokenIdentifier),
				NewToken("{", TokenSymbol),
				NewToken("field", TokenKeyword),
				NewToken("char", TokenKeyword),
				NewToken("test", TokenIdentifier),
				NewToken(";", TokenSymbol),
				NewToken("}", TokenSymbol),
			},
			want: &Class{
				Keyword:     NewKeywordByValue("class"),
				ClassName:   NewClassName(NewToken("Main", TokenIdentifier)),
				OpenSymbol:  ConstOpeningCurlyBracket,
				CloseSymbol: ConstClosingCurlyBracket,
				ClassVarDecs: &ClassVarDecs{
					Items: []*ClassVarDec{
						&ClassVarDec{
							Keyword: NewKeywordByValue("field"),
							VarType: NewVarType(NewToken("char", TokenKeyword)),
							VarNames: &VarNames{
								First:            NewVarNameByValue("test"),
								CommaAndVarNames: []*CommaAndVarName{},
							},
							EndSymbol: ConstSemicolon,
						},
					},
				},
				SubroutineDecs: NewSubroutineDecs(),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens)
			got, err := parser.parseClass()

			if err != nil {
				t.Fatalf("failed: %+v", err)
			}

			if diff := cmp.Diff(*got, *tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseClassVarDecs(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*Token
		want   *ClassVarDecs
	}{
		{
			desc: "フィールドもスタティック変数の定義もない",
			tokens: []*Token{
				NewToken("function", TokenKeyword),
				NewToken("void", TokenKeyword),
				NewToken("main", TokenIdentifier),
			},
			want: &ClassVarDecs{
				Items: []*ClassVarDec{},
			},
		},
		{
			desc: "フィールドの定義がひとつ",
			tokens: []*Token{
				NewToken("field", TokenKeyword),
				NewToken("Array", TokenIdentifier),
				NewToken("test", TokenIdentifier),
				NewToken(";", TokenSymbol),
			},
			want: &ClassVarDecs{
				Items: []*ClassVarDec{
					&ClassVarDec{
						Keyword: NewKeywordByValue("field"),
						VarType: NewVarType(NewToken("Array", TokenIdentifier)),
						VarNames: &VarNames{
							First:            NewVarNameByValue("test"),
							CommaAndVarNames: []*CommaAndVarName{},
						},
						EndSymbol: ConstSemicolon,
					},
				},
			},
		},
		{
			desc: "スタティック変数の定義がひとつ",
			tokens: []*Token{
				NewToken("static", TokenKeyword),
				NewToken("boolean", TokenKeyword),
				NewToken("test", TokenIdentifier),
				NewToken(";", TokenSymbol),
			},
			want: &ClassVarDecs{
				Items: []*ClassVarDec{
					&ClassVarDec{
						Keyword: NewKeywordByValue("static"),
						VarType: NewVarType(NewToken("boolean", TokenKeyword)),
						VarNames: &VarNames{
							First:            NewVarNameByValue("test"),
							CommaAndVarNames: []*CommaAndVarName{},
						},
						EndSymbol: ConstSemicolon,
					},
				},
			},
		},
		{
			desc: "定義が複数",
			tokens: []*Token{
				NewToken("field", TokenKeyword),
				NewToken("int", TokenKeyword),
				NewToken("foo", TokenIdentifier),
				NewToken(",", TokenSymbol),
				NewToken("bar", TokenIdentifier),
				NewToken(",", TokenSymbol),
				NewToken("baz", TokenIdentifier),
				NewToken(";", TokenSymbol),
			},
			want: &ClassVarDecs{
				Items: []*ClassVarDec{
					&ClassVarDec{
						Keyword: NewKeywordByValue("field"),
						VarType: NewVarType(NewToken("int", TokenKeyword)),
						VarNames: &VarNames{
							First: NewVarNameByValue("foo"),
							CommaAndVarNames: []*CommaAndVarName{
								NewCommaAndVarNameByValue("bar"),
								NewCommaAndVarNameByValue("baz"),
							},
						},
						EndSymbol: ConstSemicolon,
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens)
			got, err := parser.parseClassVarDecs()

			if err != nil {
				t.Fatalf("failed: %+v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserSubroutineDecs(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*Token
		want   *SubroutineDecs
	}{
		{
			desc: "サブルーチンの定義なし",
			tokens: []*Token{
				NewToken("}", TokenSymbol),
			},
			want: NewSubroutineDecs(),
		},
		{
			desc: "function void main()",
			tokens: []*Token{
				NewToken("function", TokenKeyword),
				NewToken("void", TokenKeyword),
				NewToken("main", TokenIdentifier),
				NewToken("(", TokenSymbol),
				NewToken(")", TokenSymbol),
				NewToken("{", TokenSymbol),
			},
			want: &SubroutineDecs{
				Items: []*SubroutineDec{
					&SubroutineDec{
						Subroutine:          NewKeywordByValue("function"),
						SubroutineType:      NewSubroutineType(NewToken("void", TokenKeyword)),
						SubroutineName:      NewSubroutineNameByValue("main"),
						OpeningRoundBracket: ConstOpeningRoundBracket,
						ClosingRoundBracket: ConstClosingRoundBracket,
						ParameterList:       NewParameterList(),
						SubroutineBody:      NewSubroutineBody(),
					},
				},
			},
		},
		{
			desc: "constructor Square new(int foo, int bar)",
			tokens: []*Token{
				NewToken("constructor", TokenKeyword),
				NewToken("Square", TokenIdentifier),
				NewToken("new", TokenIdentifier),
				NewToken("(", TokenSymbol),
				NewToken("int", TokenKeyword),
				NewToken("foo", TokenIdentifier),
				NewToken(",", TokenSymbol),
				NewToken("int", TokenKeyword),
				NewToken("bar", TokenIdentifier),
				NewToken(")", TokenSymbol),
				NewToken("{", TokenSymbol),
			},
			want: &SubroutineDecs{
				Items: []*SubroutineDec{
					&SubroutineDec{
						Subroutine:          NewKeywordByValue("constructor"),
						SubroutineType:      NewSubroutineType(NewToken("Square", TokenIdentifier)),
						SubroutineName:      NewSubroutineNameByValue("new"),
						OpeningRoundBracket: ConstOpeningRoundBracket,
						ClosingRoundBracket: ConstClosingRoundBracket,
						ParameterList: &ParameterList{
							First: NewParameterByToken(
								NewToken("int", TokenKeyword),
								NewToken("foo", TokenIdentifier),
							),
							CommaAndParameters: []*CommaAndParameter{
								NewCommaAndParameterByToken(
									NewToken("int", TokenKeyword),
									NewToken("bar", TokenIdentifier),
								),
							},
						},
						SubroutineBody: NewSubroutineBody(),
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens)
			got, err := parser.parseSubroutineDecs()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}

func TestParserParameterList(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*Token
		want   *ParameterList
	}{
		{
			desc: "パラメータの定義がひとつもない",
			tokens: []*Token{
				NewToken(")", TokenSymbol),
			},
			want: &ParameterList{
				CommaAndParameters: []*CommaAndParameter{},
			},
		},
		{
			desc: "パラメータの定義がひとつ",
			tokens: []*Token{
				NewToken("Array", TokenIdentifier),
				NewToken("elements", TokenIdentifier),
				NewToken(")", TokenSymbol),
			},
			want: &ParameterList{
				First: NewParameterByToken(
					NewToken("Array", TokenIdentifier),
					NewToken("elements", TokenIdentifier),
				),
				CommaAndParameters: []*CommaAndParameter{},
			},
		},
		{
			desc: "パラメータの定義が複数",
			tokens: []*Token{
				NewToken("int", TokenKeyword),
				NewToken("foo", TokenIdentifier),
				NewToken(",", TokenSymbol),
				NewToken("boolean", TokenKeyword),
				NewToken("bar", TokenIdentifier),
				NewToken(",", TokenSymbol),
				NewToken("char", TokenKeyword),
				NewToken("baz", TokenIdentifier),
				NewToken(")", TokenSymbol),
			},
			want: &ParameterList{
				First: NewParameterByToken(
					NewToken("int", TokenKeyword),
					NewToken("foo", TokenIdentifier),
				),
				CommaAndParameters: []*CommaAndParameter{
					NewCommaAndParameterByToken(
						NewToken("boolean", TokenKeyword),
						NewToken("bar", TokenIdentifier),
					),
					NewCommaAndParameterByToken(
						NewToken("char", TokenKeyword),
						NewToken("baz", TokenIdentifier),
					),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens)
			got, err := parser.parseParameterList()

			if err != nil {
				t.Fatalf("failed: %+v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

// '{' varDec* statements* '}'
func TestSubroutineBody(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*Token
		want   *SubroutineBody
	}{
		{
			desc: "ローカル変数の定義がない",
			tokens: []*Token{
				NewToken("{", TokenSymbol),
				NewToken("return", TokenKeyword),
				NewToken(";", TokenSymbol),
				NewToken("}", TokenSymbol),
			},
			want: &SubroutineBody{
				VarDecs:             NewVarDecs(),
				OpeningCurlyBracket: ConstOpeningCurlyBracket,
				ClosingCurlyBracket: ConstClosingCurlyBracket,
			},
		},
		{
			desc: "ローカル変数の定義が一行ある",
			tokens: []*Token{
				NewToken("{", TokenSymbol),
				NewToken("var", TokenKeyword),
				NewToken("boolean", TokenKeyword),
				NewToken("foo", TokenIdentifier),
				NewToken(",", TokenSymbol),
				NewToken("bar", TokenIdentifier),
				NewToken(";", TokenSymbol),
				NewToken("return", TokenKeyword),
				NewToken(";", TokenSymbol),
				NewToken("}", TokenSymbol),
			},
			want: &SubroutineBody{
				VarDecs: &VarDecs{
					Items: []*VarDec{
						&VarDec{
							Keyword: NewKeywordByValue("var"),
							VarType: NewVarType(NewToken("boolean", TokenKeyword)),
							VarNames: &VarNames{
								First: NewVarNameByValue("foo"),
								CommaAndVarNames: []*CommaAndVarName{
									NewCommaAndVarNameByValue("bar"),
								},
							},
							Semicolon: ConstSemicolon,
						},
					},
				},
				OpeningCurlyBracket: ConstOpeningCurlyBracket,
				ClosingCurlyBracket: ConstClosingCurlyBracket,
			},
		},
		{
			desc: "ローカル変数の定義が複数行ある",
			tokens: []*Token{
				NewToken("{", TokenSymbol),
				NewToken("var", TokenKeyword),
				NewToken("int", TokenKeyword),
				NewToken("foo", TokenIdentifier),
				NewToken(";", TokenSymbol),
				NewToken("var", TokenKeyword),
				NewToken("Array", TokenIdentifier),
				NewToken("elements", TokenIdentifier),
				NewToken(";", TokenSymbol),
				NewToken("return", TokenKeyword),
				NewToken(";", TokenSymbol),
				NewToken("}", TokenSymbol),
			},
			want: &SubroutineBody{
				VarDecs: &VarDecs{
					Items: []*VarDec{
						&VarDec{
							Keyword: NewKeywordByValue("var"),
							VarType: NewVarType(NewToken("int", TokenKeyword)),
							VarNames: &VarNames{
								First:            NewVarNameByValue("foo"),
								CommaAndVarNames: []*CommaAndVarName{},
							},
							Semicolon: ConstSemicolon,
						},
						&VarDec{
							Keyword: NewKeywordByValue("var"),
							VarType: NewVarType(NewToken("Array", TokenIdentifier)),
							VarNames: &VarNames{
								First:            NewVarNameByValue("elements"),
								CommaAndVarNames: []*CommaAndVarName{},
							},
							Semicolon: ConstSemicolon,
						},
					},
				},
				OpeningCurlyBracket: ConstOpeningCurlyBracket,
				ClosingCurlyBracket: ConstClosingCurlyBracket,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens)
			got, err := parser.parseSubroutineBody()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserVarDec(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*Token
		want   *VarDec
	}{
		{
			desc: "ローカル変数の定義がひとつ",
			tokens: []*Token{
				NewToken("var", TokenKeyword),
				NewToken("Array", TokenIdentifier),
				NewToken("elements", TokenIdentifier),
				NewToken(";", TokenSymbol),
			},
			want: &VarDec{
				Keyword: NewKeywordByValue("var"),
				VarType: NewVarType(NewToken("Array", TokenIdentifier)),
				VarNames: &VarNames{
					First:            NewVarNameByValue("elements"),
					CommaAndVarNames: []*CommaAndVarName{},
				},
				Semicolon: ConstSemicolon,
			},
		},
		{
			desc: "ローカル変数の定義が複数",
			tokens: []*Token{
				NewToken("var", TokenKeyword),
				NewToken("char", TokenKeyword),
				NewToken("foo", TokenIdentifier),
				NewToken(",", TokenSymbol),
				NewToken("bar", TokenIdentifier),
				NewToken(",", TokenSymbol),
				NewToken("baz", TokenIdentifier),
				NewToken(";", TokenSymbol),
			},
			want: &VarDec{
				Keyword: NewKeywordByValue("var"),
				VarType: NewVarType(NewToken("char", TokenKeyword)),
				VarNames: &VarNames{
					First: NewVarNameByValue("foo"),
					CommaAndVarNames: []*CommaAndVarName{
						NewCommaAndVarNameByValue("bar"),
						NewCommaAndVarNameByValue("baz"),
					},
				},
				Semicolon: ConstSemicolon,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens)
			got, err := parser.parseVarDec()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}
