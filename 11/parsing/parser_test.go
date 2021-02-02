package parsing

import (
	"../token"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"testing"
)

func TestParserParseClass(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *Class
	}{
		{
			desc: "ClassVarDecsとSubroutineDecが初期値のクラス",
			tokens: []*token.Token{
				token.NewToken("class", token.TokenKeyword),
				token.NewToken("Main", token.TokenIdentifier),
				token.NewToken("{", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &Class{
				Keyword:             NewKeywordByValue("class"),
				ClassName:           NewClassName(token.NewToken("Main", token.TokenIdentifier)),
				OpeningCurlyBracket: ConstOpeningCurlyBracket,
				ClosingCurlyBracket: ConstClosingCurlyBracket,
				ClassVarDecs:        NewClassVarDecs(),
				SubroutineDecs:      NewSubroutineDecs(),
			},
		},
		{
			desc: "SubroutineDecsが初期値のクラス",
			tokens: []*token.Token{
				token.NewToken("class", token.TokenKeyword),
				token.NewToken("Main", token.TokenIdentifier),
				token.NewToken("{", token.TokenSymbol),
				token.NewToken("field", token.TokenKeyword),
				token.NewToken("char", token.TokenKeyword),
				token.NewToken("test", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &Class{
				Keyword:             NewKeywordByValue("class"),
				ClassName:           NewClassName(token.NewToken("Main", token.TokenIdentifier)),
				OpeningCurlyBracket: ConstOpeningCurlyBracket,
				ClosingCurlyBracket: ConstClosingCurlyBracket,
				ClassVarDecs: &ClassVarDecs{
					Items: []*ClassVarDec{
						&ClassVarDec{
							Keyword: NewKeywordByValue("field"),
							VarType: NewVarType(token.NewToken("char", token.TokenKeyword)),
							VarNames: &VarNames{
								First:            NewVarNameByValue("test"),
								CommaAndVarNames: []*CommaAndVarName{},
							},
							Semicolon: ConstSemicolon,
						},
					},
				},
				SubroutineDecs: NewSubroutineDecs(),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseClass()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
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
		tokens []*token.Token
		want   *ClassVarDecs
	}{
		{
			desc: "フィールドもスタティック変数の定義もない",
			tokens: []*token.Token{
				token.NewToken("function", token.TokenKeyword),
				token.NewToken("void", token.TokenKeyword),
				token.NewToken("main", token.TokenIdentifier),
			},
			want: &ClassVarDecs{
				Items: []*ClassVarDec{},
			},
		},
		{
			desc: "フィールドの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("field", token.TokenKeyword),
				token.NewToken("Array", token.TokenIdentifier),
				token.NewToken("test", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &ClassVarDecs{
				Items: []*ClassVarDec{
					&ClassVarDec{
						Keyword: NewKeywordByValue("field"),
						VarType: NewVarType(token.NewToken("Array", token.TokenIdentifier)),
						VarNames: &VarNames{
							First:            NewVarNameByValue("test"),
							CommaAndVarNames: []*CommaAndVarName{},
						},
						Semicolon: ConstSemicolon,
					},
				},
			},
		},
		{
			desc: "スタティック変数の定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("static", token.TokenKeyword),
				token.NewToken("boolean", token.TokenKeyword),
				token.NewToken("test", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &ClassVarDecs{
				Items: []*ClassVarDec{
					&ClassVarDec{
						Keyword: NewKeywordByValue("static"),
						VarType: NewVarType(token.NewToken("boolean", token.TokenKeyword)),
						VarNames: &VarNames{
							First:            NewVarNameByValue("test"),
							CommaAndVarNames: []*CommaAndVarName{},
						},
						Semicolon: ConstSemicolon,
					},
				},
			},
		},
		{
			desc: "定義が複数",
			tokens: []*token.Token{
				token.NewToken("field", token.TokenKeyword),
				token.NewToken("int", token.TokenKeyword),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(",", token.TokenSymbol),
				token.NewToken("bar", token.TokenIdentifier),
				token.NewToken(",", token.TokenSymbol),
				token.NewToken("baz", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &ClassVarDecs{
				Items: []*ClassVarDec{
					&ClassVarDec{
						Keyword: NewKeywordByValue("field"),
						VarType: NewVarType(token.NewToken("int", token.TokenKeyword)),
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
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseClassVarDecs()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseSubroutineDecs(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *SubroutineDecs
	}{
		{
			desc: "サブルーチンの定義なし",
			tokens: []*token.Token{
				token.NewToken("}", token.TokenSymbol),
			},
			want: NewSubroutineDecs(),
		},
		{
			desc: "function void main()",
			tokens: []*token.Token{
				token.NewToken("function", token.TokenKeyword),
				token.NewToken("void", token.TokenKeyword),
				token.NewToken("main", token.TokenIdentifier),
				token.NewToken("(", token.TokenSymbol),
				token.NewToken(")", token.TokenSymbol),
				token.NewToken("{", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &SubroutineDecs{
				Items: []*SubroutineDec{
					&SubroutineDec{
						Subroutine:          NewKeywordByValue("function"),
						SubroutineType:      NewSubroutineType(token.NewToken("void", token.TokenKeyword)),
						SubroutineName:      NewSubroutineNameByValue("main"),
						OpeningRoundBracket: ConstOpeningRoundBracket,
						ClosingRoundBracket: ConstClosingRoundBracket,
						ParameterList:       NewParameterList(),
						SubroutineBody: &SubroutineBody{
							VarDecs:             NewVarDecs(),
							Statements:          NewStatements(),
							OpeningCurlyBracket: ConstOpeningCurlyBracket,
							ClosingCurlyBracket: ConstClosingCurlyBracket,
						},
					},
				},
			},
		},
		{
			desc: "constructor Square new(int foo, int bar)",
			tokens: []*token.Token{
				token.NewToken("constructor", token.TokenKeyword),
				token.NewToken("Square", token.TokenIdentifier),
				token.NewToken("new", token.TokenIdentifier),
				token.NewToken("(", token.TokenSymbol),
				token.NewToken("int", token.TokenKeyword),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(",", token.TokenSymbol),
				token.NewToken("int", token.TokenKeyword),
				token.NewToken("bar", token.TokenIdentifier),
				token.NewToken(")", token.TokenSymbol),
				token.NewToken("{", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &SubroutineDecs{
				Items: []*SubroutineDec{
					&SubroutineDec{
						Subroutine:          NewKeywordByValue("constructor"),
						SubroutineType:      NewSubroutineType(token.NewToken("Square", token.TokenIdentifier)),
						SubroutineName:      NewSubroutineNameByValue("new"),
						OpeningRoundBracket: ConstOpeningRoundBracket,
						ClosingRoundBracket: ConstClosingRoundBracket,
						ParameterList: &ParameterList{
							First: NewParameterByToken(
								token.NewToken("int", token.TokenKeyword),
								token.NewToken("foo", token.TokenIdentifier),
							),
							CommaAndParameters: []*CommaAndParameter{
								NewCommaAndParameterByToken(
									token.NewToken("int", token.TokenKeyword),
									token.NewToken("bar", token.TokenIdentifier),
								),
							},
						},
						SubroutineBody: &SubroutineBody{
							VarDecs:             NewVarDecs(),
							Statements:          NewStatements(),
							OpeningCurlyBracket: ConstOpeningCurlyBracket,
							ClosingCurlyBracket: ConstClosingCurlyBracket,
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseSubroutineDecs()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}

func TestParserParseParameterList(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *ParameterList
	}{
		{
			desc: "パラメータの定義がひとつもない",
			tokens: []*token.Token{
				token.NewToken(")", token.TokenSymbol),
			},
			want: &ParameterList{
				CommaAndParameters: []*CommaAndParameter{},
			},
		},
		{
			desc: "パラメータの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("Array", token.TokenIdentifier),
				token.NewToken("elements", token.TokenIdentifier),
				token.NewToken(")", token.TokenSymbol),
			},
			want: &ParameterList{
				First: NewParameterByToken(
					token.NewToken("Array", token.TokenIdentifier),
					token.NewToken("elements", token.TokenIdentifier),
				),
				CommaAndParameters: []*CommaAndParameter{},
			},
		},
		{
			desc: "パラメータの定義が複数",
			tokens: []*token.Token{
				token.NewToken("int", token.TokenKeyword),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(",", token.TokenSymbol),
				token.NewToken("boolean", token.TokenKeyword),
				token.NewToken("bar", token.TokenIdentifier),
				token.NewToken(",", token.TokenSymbol),
				token.NewToken("char", token.TokenKeyword),
				token.NewToken("baz", token.TokenIdentifier),
				token.NewToken(")", token.TokenSymbol),
			},
			want: &ParameterList{
				First: NewParameterByToken(
					token.NewToken("int", token.TokenKeyword),
					token.NewToken("foo", token.TokenIdentifier),
				),
				CommaAndParameters: []*CommaAndParameter{
					NewCommaAndParameterByToken(
						token.NewToken("boolean", token.TokenKeyword),
						token.NewToken("bar", token.TokenIdentifier),
					),
					NewCommaAndParameterByToken(
						token.NewToken("char", token.TokenKeyword),
						token.NewToken("baz", token.TokenIdentifier),
					),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
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
func TestParserParseSubroutineBody(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *SubroutineBody
	}{
		{
			desc: "ローカル変数の定義がない",
			tokens: []*token.Token{
				token.NewToken("{", token.TokenSymbol),
				token.NewToken("return", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &SubroutineBody{
				VarDecs: NewVarDecs(),
				Statements: &Statements{
					Items: []Statement{
						NewReturnStatement(),
					},
				},
				OpeningCurlyBracket: ConstOpeningCurlyBracket,
				ClosingCurlyBracket: ConstClosingCurlyBracket,
			},
		},
		{
			desc: "ローカル変数の定義が一行ある",
			tokens: []*token.Token{
				token.NewToken("{", token.TokenSymbol),
				token.NewToken("var", token.TokenKeyword),
				token.NewToken("boolean", token.TokenKeyword),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(",", token.TokenSymbol),
				token.NewToken("bar", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("return", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &SubroutineBody{
				VarDecs: &VarDecs{
					Items: []*VarDec{
						&VarDec{
							Keyword: NewKeywordByValue("var"),
							VarType: NewVarType(token.NewToken("boolean", token.TokenKeyword)),
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
				Statements: &Statements{
					Items: []Statement{
						NewReturnStatement(),
					},
				},
				OpeningCurlyBracket: ConstOpeningCurlyBracket,
				ClosingCurlyBracket: ConstClosingCurlyBracket,
			},
		},
		{
			desc: "ローカル変数の定義が複数行ある",
			tokens: []*token.Token{
				token.NewToken("{", token.TokenSymbol),
				token.NewToken("var", token.TokenKeyword),
				token.NewToken("int", token.TokenKeyword),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("var", token.TokenKeyword),
				token.NewToken("Array", token.TokenIdentifier),
				token.NewToken("elements", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("return", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &SubroutineBody{
				VarDecs: &VarDecs{
					Items: []*VarDec{
						&VarDec{
							Keyword: NewKeywordByValue("var"),
							VarType: NewVarType(token.NewToken("int", token.TokenKeyword)),
							VarNames: &VarNames{
								First:            NewVarNameByValue("foo"),
								CommaAndVarNames: []*CommaAndVarName{},
							},
							Semicolon: ConstSemicolon,
						},
						&VarDec{
							Keyword: NewKeywordByValue("var"),
							VarType: NewVarType(token.NewToken("Array", token.TokenIdentifier)),
							VarNames: &VarNames{
								First:            NewVarNameByValue("elements"),
								CommaAndVarNames: []*CommaAndVarName{},
							},
							Semicolon: ConstSemicolon,
						},
					},
				},
				Statements: &Statements{
					Items: []Statement{
						NewReturnStatement(),
					},
				},
				OpeningCurlyBracket: ConstOpeningCurlyBracket,
				ClosingCurlyBracket: ConstClosingCurlyBracket,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseSubroutineBody()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseVarDec(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *VarDec
	}{
		{
			desc: "ローカル変数の定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("var", token.TokenKeyword),
				token.NewToken("Array", token.TokenIdentifier),
				token.NewToken("elements", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &VarDec{
				Keyword: NewKeywordByValue("var"),
				VarType: NewVarType(token.NewToken("Array", token.TokenIdentifier)),
				VarNames: &VarNames{
					First:            NewVarNameByValue("elements"),
					CommaAndVarNames: []*CommaAndVarName{},
				},
				Semicolon: ConstSemicolon,
			},
		},
		{
			desc: "ローカル変数の定義が複数",
			tokens: []*token.Token{
				token.NewToken("var", token.TokenKeyword),
				token.NewToken("char", token.TokenKeyword),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(",", token.TokenSymbol),
				token.NewToken("bar", token.TokenIdentifier),
				token.NewToken(",", token.TokenSymbol),
				token.NewToken("baz", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &VarDec{
				Keyword: NewKeywordByValue("var"),
				VarType: NewVarType(token.NewToken("char", token.TokenKeyword)),
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
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseVarDec()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseStatements(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *Statements
	}{
		{
			desc: "Statementの定義がひとつもない",
			tokens: []*token.Token{
				token.NewToken("}", token.TokenSymbol),
			},
			want: NewStatements(),
		},
		{
			desc: "Statementの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("return", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &Statements{
				Items: []Statement{
					NewReturnStatement(),
				},
			},
		},
		{
			desc: "Statementの定義が複数",
			tokens: []*token.Token{
				token.NewToken("return", token.TokenKeyword),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("return", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &Statements{
				Items: []Statement{
					&ReturnStatement{
						StatementKeyword: NewStatementKeyword("return"),
						Expression: &Expression{
							Term: NewVarNameByValue("foo"),
						},
						Semicolon: ConstSemicolon,
					},
					NewReturnStatement(),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseStatements()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseStatement(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *ReturnStatement
	}{
		{
			desc: "return",
			tokens: []*token.Token{
				token.NewToken("return", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
			},
			want: NewReturnStatement(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseStatement()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseLetStatement(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *LetStatement
	}{
		{
			desc: "VarNameの代入",
			tokens: []*token.Token{
				token.NewToken("let", token.TokenKeyword),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken("=", token.TokenSymbol),
				token.NewToken("bar", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &LetStatement{
				StatementKeyword: NewStatementKeyword("let"),
				VarName:          NewVarNameByValue("foo"),
				Expression: &Expression{
					Term: NewVarNameByValue("bar"),
				},
				Equal:     ConstEqual,
				Semicolon: ConstSemicolon,
			},
		},
		{
			desc: "KeywordConstantの代入",
			tokens: []*token.Token{
				token.NewToken("let", token.TokenKeyword),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken("=", token.TokenSymbol),
				token.NewToken("true", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &LetStatement{
				StatementKeyword: NewStatementKeyword("let"),
				VarName:          NewVarNameByValue("foo"),
				Expression: &Expression{
					Term: ConstTrue,
				},
				Equal:     ConstEqual,
				Semicolon: ConstSemicolon,
			},
		},
		{
			desc: "配列に対する代入",
			tokens: []*token.Token{
				token.NewToken("let", token.TokenKeyword),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken("[", token.TokenSymbol),
				token.NewToken("index", token.TokenIdentifier),
				token.NewToken("]", token.TokenSymbol),
				token.NewToken("=", token.TokenSymbol),
				token.NewToken("bar", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &LetStatement{
				StatementKeyword: NewStatementKeyword("let"),
				Array: &Array{
					VarName: NewVarNameByValue("foo"),
					Expression: &Expression{
						Term: NewVarNameByValue("index"),
					},
					OpeningSquareBracket: ConstOpeningSquareBracket,
					ClosingSquareBracket: ConstClosingSquareBracket,
				},
				Expression: &Expression{
					Term: NewVarNameByValue("bar"),
				},
				Equal:     ConstEqual,
				Semicolon: ConstSemicolon,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseLetStatement()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseIfStatement(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *IfStatement
	}{
		{
			desc: "if句のみ",
			tokens: []*token.Token{
				token.NewToken("if", token.TokenKeyword),
				token.NewToken("(", token.TokenSymbol),
				token.NewToken("true", token.TokenKeyword),
				token.NewToken(")", token.TokenSymbol),
				token.NewToken("{", token.TokenSymbol),
				token.NewToken("return", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &IfStatement{
				StatementKeyword: NewStatementKeyword("if"),
				Expression: &Expression{
					Term: ConstTrue,
				},
				Statements: &Statements{
					Items: []Statement{
						NewReturnStatement(),
					},
				},
				OpeningRoundBracket: ConstOpeningRoundBracket,
				ClosingRoundBracket: ConstClosingRoundBracket,
				OpeningCurlyBracket: ConstOpeningCurlyBracket,
				ClosingCurlyBracket: ConstClosingCurlyBracket,
			},
		},
		{
			desc: "if-else句",
			tokens: []*token.Token{
				token.NewToken("if", token.TokenKeyword),
				token.NewToken("(", token.TokenSymbol),
				token.NewToken("true", token.TokenKeyword),
				token.NewToken(")", token.TokenSymbol),
				token.NewToken("{", token.TokenSymbol),
				token.NewToken("return", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
				token.NewToken("else", token.TokenKeyword),
				token.NewToken("{", token.TokenSymbol),
				token.NewToken("return", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &IfStatement{
				StatementKeyword: NewStatementKeyword("if"),
				Expression: &Expression{
					Term: ConstTrue,
				},
				Statements: &Statements{
					Items: []Statement{
						NewReturnStatement(),
					},
				},
				ElseBlock: &ElseBlock{
					Keyword: NewKeywordByValue("else"),
					Statements: &Statements{
						Items: []Statement{
							NewReturnStatement(),
						},
					},
					OpeningCurlyBracket: ConstOpeningCurlyBracket,
					ClosingCurlyBracket: ConstClosingCurlyBracket,
				},
				OpeningRoundBracket: ConstOpeningRoundBracket,
				ClosingRoundBracket: ConstClosingRoundBracket,
				OpeningCurlyBracket: ConstOpeningCurlyBracket,
				ClosingCurlyBracket: ConstClosingCurlyBracket,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseIfStatement()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseWhileStatement(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *WhileStatement
	}{
		{
			desc: "単一のwhile定義",
			tokens: []*token.Token{
				token.NewToken("while", token.TokenKeyword),
				token.NewToken("(", token.TokenSymbol),
				token.NewToken("true", token.TokenKeyword),
				token.NewToken(")", token.TokenSymbol),
				token.NewToken("{", token.TokenSymbol),
				token.NewToken("return", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &WhileStatement{
				StatementKeyword: NewStatementKeyword("while"),
				Expression: &Expression{
					Term: ConstTrue,
				},
				Statements: &Statements{
					Items: []Statement{
						NewReturnStatement(),
					},
				},
				OpeningRoundBracket: ConstOpeningRoundBracket,
				ClosingRoundBracket: ConstClosingRoundBracket,
				OpeningCurlyBracket: ConstOpeningCurlyBracket,
				ClosingCurlyBracket: ConstClosingCurlyBracket,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseWhileStatement()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseDoStatement(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *DoStatement
	}{
		{
			desc: "引数なしのサブルーチンの実行",
			tokens: []*token.Token{
				token.NewToken("do", token.TokenKeyword),
				token.NewToken("max", token.TokenIdentifier),
				token.NewToken("(", token.TokenSymbol),
				token.NewToken(")", token.TokenSymbol),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &DoStatement{
				StatementKeyword: NewStatementKeyword("do"),
				SubroutineCall: &SubroutineCall{
					SubroutineCallName: &SubroutineCallName{
						Period:         ConstPeriod,
						SubroutineName: NewSubroutineName(token.NewToken("max", token.TokenIdentifier)),
					},
					ExpressionList:      NewExpressionList(),
					OpeningRoundBracket: ConstOpeningRoundBracket,
					ClosingRoundBracket: ConstClosingRoundBracket,
				},
				Semicolon: ConstSemicolon,
			},
		},
		{
			desc: "引数ありのサブルーチンの実行",
			tokens: []*token.Token{
				token.NewToken("do", token.TokenKeyword),
				token.NewToken("run", token.TokenIdentifier),
				token.NewToken("(", token.TokenSymbol),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(")", token.TokenSymbol),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &DoStatement{
				StatementKeyword: NewStatementKeyword("do"),
				SubroutineCall: &SubroutineCall{
					SubroutineCallName: &SubroutineCallName{
						Period:         ConstPeriod,
						SubroutineName: NewSubroutineName(token.NewToken("run", token.TokenIdentifier)),
					},
					ExpressionList: &ExpressionList{
						First: &Expression{
							Term: NewVarNameByValue("foo"),
						},
						CommaAndExpressions: []*CommaAndExpression{},
					},
					OpeningRoundBracket: ConstOpeningRoundBracket,
					ClosingRoundBracket: ConstClosingRoundBracket,
				},
				Semicolon: ConstSemicolon,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseDoStatement()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseReturnStatement(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *ReturnStatement
	}{
		{
			desc: "セミコロンのみ",
			tokens: []*token.Token{
				token.NewToken("return", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
			},
			want: NewReturnStatement(),
		},
		{
			desc: "式とセミコロン",
			tokens: []*token.Token{
				token.NewToken("return", token.TokenKeyword),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &ReturnStatement{
				StatementKeyword: NewStatementKeyword("return"),
				Expression: &Expression{
					Term: NewVarNameByValue("foo"),
				},
				Semicolon: ConstSemicolon,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseReturnStatement()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseElseBlock(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *ElseBlock
	}{
		{
			desc: "else句",
			tokens: []*token.Token{
				token.NewToken("else", token.TokenKeyword),
				token.NewToken("{", token.TokenSymbol),
				token.NewToken("return", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
				token.NewToken("}", token.TokenSymbol),
			},
			want: &ElseBlock{
				Keyword: NewKeywordByValue("else"),
				Statements: &Statements{
					Items: []Statement{
						NewReturnStatement(),
					},
				},
				OpeningCurlyBracket: ConstOpeningCurlyBracket,
				ClosingCurlyBracket: ConstClosingCurlyBracket,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseElseBlock()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseSubroutineCall(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *SubroutineCall
	}{
		{
			desc: "引数なし",
			tokens: []*token.Token{
				token.NewToken("Main", token.TokenIdentifier),
				token.NewToken(".", token.TokenSymbol),
				token.NewToken("main", token.TokenIdentifier),
				token.NewToken("(", token.TokenSymbol),
				token.NewToken(")", token.TokenSymbol),
			},
			want: &SubroutineCall{
				SubroutineCallName: &SubroutineCallName{
					CallerName:     NewCallerName(token.NewToken("Main", token.TokenIdentifier)),
					Period:         ConstPeriod,
					SubroutineName: NewSubroutineName(token.NewToken("main", token.TokenIdentifier)),
				},
				ExpressionList:      NewExpressionList(),
				OpeningRoundBracket: ConstOpeningRoundBracket,
				ClosingRoundBracket: ConstClosingRoundBracket,
			},
		},
		{
			desc: "引数あり",
			tokens: []*token.Token{
				token.NewToken("join", token.TokenIdentifier),
				token.NewToken("(", token.TokenSymbol),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(",", token.TokenSymbol),
				token.NewToken("bar", token.TokenIdentifier),
				token.NewToken(",", token.TokenSymbol),
				token.NewToken("baz", token.TokenIdentifier),
				token.NewToken(")", token.TokenSymbol),
			},
			want: &SubroutineCall{
				SubroutineCallName: &SubroutineCallName{
					Period:         ConstPeriod,
					SubroutineName: NewSubroutineName(token.NewToken("join", token.TokenIdentifier)),
				},
				ExpressionList: &ExpressionList{
					First: &Expression{
						Term: NewVarNameByValue("foo"),
					},
					CommaAndExpressions: []*CommaAndExpression{
						NewCommaAndExpression(&Expression{
							Term: NewVarNameByValue("bar"),
						}),
						NewCommaAndExpression(&Expression{
							Term: NewVarNameByValue("baz"),
						}),
					},
				},
				OpeningRoundBracket: ConstOpeningRoundBracket,
				ClosingRoundBracket: ConstClosingRoundBracket,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseSubroutineCall()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseSubroutineCallName(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *SubroutineCallName
	}{
		{
			desc: "サブルーチン名のみ",
			tokens: []*token.Token{
				token.NewToken("max", token.TokenIdentifier),
				token.NewToken("(", token.TokenSymbol),
			},
			want: &SubroutineCallName{
				Period:         ConstPeriod,
				SubroutineName: NewSubroutineName(token.NewToken("max", token.TokenIdentifier)),
			},
		},
		{
			desc: "クラス名＋サブルーチン名",
			tokens: []*token.Token{
				token.NewToken("Array", token.TokenIdentifier),
				token.NewToken(".", token.TokenSymbol),
				token.NewToken("new", token.TokenIdentifier),
				token.NewToken("(", token.TokenSymbol),
			},
			want: &SubroutineCallName{
				CallerName:     NewCallerName(token.NewToken("Array", token.TokenIdentifier)),
				Period:         ConstPeriod,
				SubroutineName: NewSubroutineName(token.NewToken("new", token.TokenIdentifier)),
			},
		},
		{
			desc: "インスタンス名＋サブルーチン名",
			tokens: []*token.Token{
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(".", token.TokenSymbol),
				token.NewToken("run", token.TokenIdentifier),
				token.NewToken("(", token.TokenSymbol),
			},
			want: &SubroutineCallName{
				CallerName:     NewCallerName(token.NewToken("foo", token.TokenIdentifier)),
				Period:         ConstPeriod,
				SubroutineName: NewSubroutineName(token.NewToken("run", token.TokenIdentifier)),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseSubroutineCallName()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseExpressionList(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *ExpressionList
	}{
		{
			desc: "Expressionの定義がない",
			tokens: []*token.Token{
				token.NewToken(")", token.TokenSymbol),
			},
			want: NewExpressionList(),
		},
		{
			desc: "Identifierの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &ExpressionList{
				First: &Expression{
					Term: NewVarNameByValue("foo"),
				},
				CommaAndExpressions: []*CommaAndExpression{},
			},
		},
		{
			desc: "Keywordの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("this", token.TokenKeyword),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &ExpressionList{
				First: &Expression{
					Term: ConstThis,
				},
				CommaAndExpressions: []*CommaAndExpression{},
			},
		},
		{
			desc: "Expressionの定義が複数",
			tokens: []*token.Token{
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(",", token.TokenSymbol),
				token.NewToken("bar", token.TokenIdentifier),
				token.NewToken(",", token.TokenSymbol),
				token.NewToken("baz", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &ExpressionList{
				First: &Expression{
					Term: NewVarNameByValue("foo"),
				},
				CommaAndExpressions: []*CommaAndExpression{
					NewCommaAndExpression(&Expression{
						Term: NewVarNameByValue("bar"),
					}),
					NewCommaAndExpression(&Expression{
						Term: NewVarNameByValue("baz"),
					}),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseExpressionList()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseExpression(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *Expression
	}{
		{
			desc: "Termの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &Expression{
				Term: NewVarNameByValue("foo"),
			},
		},
		{
			desc: "Termの定義が複数",
			tokens: []*token.Token{
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken("|", token.TokenSymbol),
				token.NewToken("bar", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &Expression{
				Term: NewVarNameByValue("foo"),
				BinaryOpTerms: &BinaryOpTerms{
					Items: []*BinaryOpTerm{
						&BinaryOpTerm{
							BinaryOp: ConstVerticalLine,
							Term:     NewVarNameByValue("bar"),
						},
					},
				},
			},
		},
		{
			desc: "カッコを含む複雑なExpressionの定義",
			tokens: []*token.Token{
				token.NewToken("2", token.TokenIntConst),
				token.NewToken("*", token.TokenSymbol),
				token.NewToken("(", token.TokenSymbol),
				token.NewToken("2", token.TokenIntConst),
				token.NewToken("+", token.TokenSymbol),
				token.NewToken("3", token.TokenIntConst),
				token.NewToken(")", token.TokenSymbol),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &Expression{
				Term: &IntegerConstant{
					Token: token.NewToken("2", token.TokenIntConst),
				},
				BinaryOpTerms: &BinaryOpTerms{
					Items: []*BinaryOpTerm{
						&BinaryOpTerm{
							BinaryOp: ConstAsterisk,
							Term: &GroupingExpression{
								Expression: &Expression{
									Term: &IntegerConstant{
										Token: token.NewToken("2", token.TokenIntConst),
									},
									BinaryOpTerms: &BinaryOpTerms{
										Items: []*BinaryOpTerm{
											&BinaryOpTerm{
												BinaryOp: ConstPlus,
												Term: &IntegerConstant{
													Token: token.NewToken("3", token.TokenIntConst),
												},
											},
										},
									},
								},
								OpeningRoundBracket: ConstOpeningRoundBracket,
								ClosingRoundBracket: ConstClosingRoundBracket,
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseExpression()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseBinaryOpTerms(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *BinaryOpTerms
	}{
		{
			desc: "BinaryOpTermの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("|", token.TokenSymbol),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &BinaryOpTerms{
				Items: []*BinaryOpTerm{
					&BinaryOpTerm{
						BinaryOp: ConstVerticalLine,
						Term:     NewVarNameByValue("foo"),
					},
				},
			},
		},
		{
			desc: "BinaryOpTermの定義が複数",
			tokens: []*token.Token{
				token.NewToken("|", token.TokenSymbol),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken("|", token.TokenSymbol),
				token.NewToken("bar", token.TokenIdentifier),
				token.NewToken("|", token.TokenSymbol),
				token.NewToken("baz", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &BinaryOpTerms{
				Items: []*BinaryOpTerm{
					&BinaryOpTerm{
						BinaryOp: ConstVerticalLine,
						Term:     NewVarNameByValue("foo"),
					},
					&BinaryOpTerm{
						BinaryOp: ConstVerticalLine,
						Term:     NewVarNameByValue("bar"),
					},
					&BinaryOpTerm{
						BinaryOp: ConstVerticalLine,
						Term:     NewVarNameByValue("baz"),
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseBinaryOpTerms()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseTerm(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   Term
	}{
		{
			desc: "IntegerConstantへの分岐",
			tokens: []*token.Token{
				token.NewToken("123", token.TokenIntConst),
			},
			want: &IntegerConstant{
				Token: token.NewToken("123", token.TokenIntConst),
			},
		},
		{
			desc: "StringConstantへの分岐",
			tokens: []*token.Token{
				token.NewToken("foo bar", token.TokenStringConst),
			},
			want: &StringConstant{
				Token: token.NewToken("foo bar", token.TokenStringConst),
			},
		},
		{
			desc: "KeywordConstantへの分岐",
			tokens: []*token.Token{
				token.NewToken("true", token.TokenKeyword),
			},
			want: ConstTrue,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseTerm()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseIntegerConstant(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   Term
	}{
		{
			desc: "IntegerConstantの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("123", token.TokenIntConst),
			},
			want: &IntegerConstant{
				Token: token.NewToken("123", token.TokenIntConst),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseIntegerConstant()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseStringConstant(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   Term
	}{
		{
			desc: "StringConstantの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("foo bar", token.TokenStringConst),
			},
			want: &StringConstant{
				Token: token.NewToken("foo bar", token.TokenStringConst),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseStringConstant()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseKeywordConstant(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   Term
	}{
		{
			desc: "trueの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("true", token.TokenKeyword),
			},
			want: ConstTrue,
		},
		{
			desc: "falseの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("false", token.TokenKeyword),
			},
			want: ConstFalse,
		},
		{
			desc: "nullの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("null", token.TokenKeyword),
			},
			want: ConstNull,
		},
		{
			desc: "thisの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("this", token.TokenKeyword),
			},
			want: ConstThis,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseKeywordConstant()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseIdentifierTerm(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   Term
	}{
		{
			desc: "VarNameの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: NewVarNameByValue("foo"),
		},
		{
			desc: "Arrayの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("array", token.TokenIdentifier),
				token.NewToken("[", token.TokenSymbol),
				token.NewToken("index", token.TokenIdentifier),
				token.NewToken("]", token.TokenSymbol),
			},
			want: &Array{
				VarName: NewVarNameByValue("array"),
				Expression: &Expression{
					Term: NewVarNameByValue("index"),
				},
				OpeningSquareBracket: ConstOpeningSquareBracket,
				ClosingSquareBracket: ConstClosingSquareBracket,
			},
		},
		{
			desc: "クラス内のSubroutineCallの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("main", token.TokenIdentifier),
				token.NewToken("(", token.TokenSymbol),
				token.NewToken(")", token.TokenSymbol),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &SubroutineCall{
				SubroutineCallName: &SubroutineCallName{
					Period:         ConstPeriod,
					SubroutineName: NewSubroutineName(token.NewToken("main", token.TokenIdentifier)),
				},
				ExpressionList:      NewExpressionList(),
				OpeningRoundBracket: ConstOpeningRoundBracket,
				ClosingRoundBracket: ConstClosingRoundBracket,
			},
		},
		{
			desc: "クラス外のSubroutineCallの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("Main", token.TokenIdentifier),
				token.NewToken(".", token.TokenSymbol),
				token.NewToken("main", token.TokenIdentifier),
				token.NewToken("(", token.TokenSymbol),
				token.NewToken(")", token.TokenSymbol),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &SubroutineCall{
				SubroutineCallName: &SubroutineCallName{
					CallerName:     NewCallerName(token.NewToken("Main", token.TokenIdentifier)),
					Period:         ConstPeriod,
					SubroutineName: NewSubroutineName(token.NewToken("main", token.TokenIdentifier)),
				},
				ExpressionList:      NewExpressionList(),
				OpeningRoundBracket: ConstOpeningRoundBracket,
				ClosingRoundBracket: ConstClosingRoundBracket,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseIdentifierTerm()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseArray(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *Array
	}{
		{
			desc: "Arrayの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("array", token.TokenIdentifier),
				token.NewToken("[", token.TokenSymbol),
				token.NewToken("index", token.TokenIdentifier),
				token.NewToken("]", token.TokenSymbol),
			},
			want: &Array{
				VarName: NewVarNameByValue("array"),
				Expression: &Expression{
					Term: NewVarNameByValue("index"),
				},
				OpeningSquareBracket: ConstOpeningSquareBracket,
				ClosingSquareBracket: ConstClosingSquareBracket,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseArray()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseSymbolTerm(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   Term
	}{
		{
			desc: "Minusの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("-", token.TokenSymbol),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &UnaryOpTerm{
				UnaryOp: ConstUnaryMinus,
				Term:    NewVarNameByValue("foo"),
			},
		},
		{
			desc: "Tildeの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("~", token.TokenSymbol),
				token.NewToken("123", token.TokenIntConst),
			},
			want: &UnaryOpTerm{
				UnaryOp: ConstUnaryTilde,
				Term: &IntegerConstant{
					Token: token.NewToken("123", token.TokenIntConst),
				},
			},
		},
		{
			desc: "GroupingExpressionの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("(", token.TokenSymbol),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(")", token.TokenSymbol),
			},
			want: &GroupingExpression{
				Expression: &Expression{
					Term: NewVarNameByValue("foo"),
				},
				OpeningRoundBracket: ConstOpeningRoundBracket,
				ClosingRoundBracket: ConstClosingRoundBracket,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseSymbolTerm()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseUnaryOpTerm(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *UnaryOpTerm
	}{
		{
			desc: "Minusの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("-", token.TokenSymbol),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(";", token.TokenSymbol),
			},
			want: &UnaryOpTerm{
				UnaryOp: ConstUnaryMinus,
				Term:    NewVarNameByValue("foo"),
			},
		},
		{
			desc: "Tildeの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("~", token.TokenSymbol),
				token.NewToken("123", token.TokenIntConst),
			},
			want: &UnaryOpTerm{
				UnaryOp: ConstUnaryTilde,
				Term: &IntegerConstant{
					Token: token.NewToken("123", token.TokenIntConst),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseUnaryOpTerm()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParserParseGroupingExpression(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*token.Token
		want   *GroupingExpression
	}{
		{
			desc: "Expressionの定義がひとつ",
			tokens: []*token.Token{
				token.NewToken("(", token.TokenSymbol),
				token.NewToken("foo", token.TokenIdentifier),
				token.NewToken(")", token.TokenSymbol),
			},
			want: &GroupingExpression{
				Expression: &Expression{
					Term: NewVarNameByValue("foo"),
				},
				OpeningRoundBracket: ConstOpeningRoundBracket,
				ClosingRoundBracket: ConstClosingRoundBracket,
			},
		},
		{
			desc: "Expressionの定義がふたつ",
			tokens: []*token.Token{
				token.NewToken("(", token.TokenSymbol),
				token.NewToken("2", token.TokenIntConst),
				token.NewToken("+", token.TokenSymbol),
				token.NewToken("3", token.TokenIntConst),
				token.NewToken(")", token.TokenSymbol),
			},
			want: &GroupingExpression{
				Expression: &Expression{
					Term: &IntegerConstant{
						Token: token.NewToken("2", token.TokenIntConst),
					},
					BinaryOpTerms: &BinaryOpTerms{
						Items: []*BinaryOpTerm{
							&BinaryOpTerm{
								BinaryOp: ConstPlus,
								Term: &IntegerConstant{
									Token: token.NewToken("3", token.TokenIntConst),
								},
							},
						},
					},
				},
				OpeningRoundBracket: ConstOpeningRoundBracket,
				ClosingRoundBracket: ConstClosingRoundBracket,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := token.NewTokens()
			tokens.Add(tc.tokens)

			parser := NewParser(tokens, "Test")
			got, err := parser.parseGroupingExpression()

			if err != nil {
				t.Fatalf("failed %s: %+v", tc.desc, errors.WithMessage(err, parser.tokens.Debug()))
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}
