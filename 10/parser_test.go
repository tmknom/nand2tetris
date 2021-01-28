package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParserParse(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*Token
		want   *Class
	}{
		{
			desc: "parseClassメソッドを呼んでいることを確認",
			tokens: []*Token{
				NewToken("class", TokenKeyword),
				NewToken("Main", TokenIdentifier),
				NewToken("{", TokenSymbol),
				NewToken("}", TokenSymbol),
			},
			want: &Class{
				Keyword:       NewKeyword("class"),
				ClassName:     NewClassName(NewToken("Main", TokenIdentifier)),
				OpenSymbol:    NewSymbol("{"),
				CloseSymbol:   NewSymbol("}"),
				ClassVarDecs:  NewClassVarDecs(),
				SubroutineDec: []*Token{},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens)
			got, err := parser.Parse()

			if err != nil {
				t.Fatalf("failed: %+v", err)
			}

			if diff := cmp.Diff(*got, *tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

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
				Keyword:       NewKeyword("class"),
				ClassName:     NewClassName(NewToken("Main", TokenIdentifier)),
				OpenSymbol:    NewSymbol("{"),
				CloseSymbol:   NewSymbol("}"),
				ClassVarDecs:  NewClassVarDecs(),
				SubroutineDec: []*Token{},
			},
		},
		{
			desc: "ClassVarDecsが初期値のクラス",
			tokens: []*Token{
				NewToken("class", TokenKeyword),
				NewToken("Main", TokenIdentifier),
				NewToken("{", TokenSymbol),
				NewToken("field", TokenKeyword),
				NewToken("Array", TokenIdentifier),
				NewToken("test", TokenIdentifier),
				NewToken(";", TokenSymbol),
				NewToken("}", TokenSymbol),
			},
			want: &Class{
				Keyword:     NewKeyword("class"),
				ClassName:   NewClassName(NewToken("Main", TokenIdentifier)),
				OpenSymbol:  NewSymbol("{"),
				CloseSymbol: NewSymbol("}"),
				ClassVarDecs: &ClassVarDecs{
					Items: []*ClassVarDec{
						&ClassVarDec{
							Keyword: NewKeyword("field"),
							VarType: NewVarType(NewToken("Array", TokenIdentifier)),
							VarNames: &VarNames{
								First:            NewToken("test", TokenIdentifier),
								CommaAndVarNames: []*CommaAndVarName{},
							},
							EndSymbol: NewSymbol(";"),
						},
					},
				},
				SubroutineDec: []*Token{},
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
						Keyword: NewKeyword("field"),
						VarType: NewVarType(NewToken("Array", TokenIdentifier)),
						VarNames: &VarNames{
							First:            NewToken("test", TokenIdentifier),
							CommaAndVarNames: []*CommaAndVarName{},
						},
						EndSymbol: NewSymbol(";"),
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
						Keyword: NewKeyword("static"),
						VarType: NewVarType(NewToken("boolean", TokenKeyword)),
						VarNames: &VarNames{
							First:            NewToken("test", TokenIdentifier),
							CommaAndVarNames: []*CommaAndVarName{},
						},
						EndSymbol: NewSymbol(";"),
					},
				},
			},
		},
		{
			desc: "定義が複数",
			tokens: []*Token{
				NewToken("field", TokenKeyword),
				NewToken("Array", TokenIdentifier),
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
						Keyword: NewKeyword("field"),
						VarType: NewVarType(NewToken("Array", TokenIdentifier)),
						VarNames: &VarNames{
							First: NewToken("foo", TokenIdentifier),
							CommaAndVarNames: []*CommaAndVarName{
								&CommaAndVarName{
									Comma:   NewToken(",", TokenSymbol),
									VarName: NewToken("bar", TokenIdentifier),
								},
								&CommaAndVarName{
									Comma:   NewToken(",", TokenSymbol),
									VarName: NewToken("baz", TokenIdentifier),
								},
							},
						},
						EndSymbol: NewSymbol(";"),
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
