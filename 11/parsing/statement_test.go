package parsing

import (
	"../symbol"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestLetStatementToCode(t *testing.T) {
	cases := []struct {
		desc         string
		letStatement *LetStatement
		want         []string
	}{
		{
			desc: "VarNameへの代入: let foo = 123",
			letStatement: &LetStatement{
				StatementKeyword: NewStatementKeyword("let"),
				VarName:          NewVarNameByValue("foo"),
				Expression: &Expression{
					Term: NewIntegerConstantByValue("123"),
				},
				Equal:     ConstEqual,
				Semicolon: ConstSemicolon,
			},
			want: []string{
				"push constant 123",
				"pop local 0",
			},
		},
	}

	// シンボルテーブルのセットアップ
	SetupForTestLetStatementToCode()

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.letStatement.ToCode()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func SetupForTestLetStatementToCode() {
	symbol.GlobalSymbolTables.Reset("Testing")
	symbol.GlobalSymbolTables.ResetSubroutine("TestRun")
	symbol.GlobalSymbolTables.AddVarSymbol("foo", "int")
}

func TestDoStatementToCode(t *testing.T) {
	cases := []struct {
		desc        string
		doStatement *DoStatement
		want        []string
	}{
		{
			desc: "引数なしのサブルーチンの実行: do max();",
			doStatement: &DoStatement{
				StatementKeyword: NewStatementKeyword("do"),
				SubroutineCall: &SubroutineCall{
					SubroutineCallName: &SubroutineCallName{
						Period:         ConstPeriod,
						SubroutineName: NewSubroutineNameByValue("max"),
					},
					ExpressionList:      NewExpressionList(),
					OpeningRoundBracket: ConstOpeningRoundBracket,
					ClosingRoundBracket: ConstClosingRoundBracket,
				},
				Semicolon: ConstSemicolon,
			},
			want: []string{
				"call max 0",
				"pop temp 0",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.doStatement.ToCode()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestReturnStatementToCode(t *testing.T) {
	cases := []struct {
		desc            string
		returnStatement *ReturnStatement
		want            []string
	}{
		{
			desc:            "値をなにも返さないreturn文: return;",
			returnStatement: NewReturnStatement(),
			want: []string{
				"push constant 0",
				"return",
			},
		},
		{
			desc: "IntegerConstantを返すreturn文",
			returnStatement: &ReturnStatement{
				StatementKeyword: NewStatementKeyword("return"),
				Expression: &Expression{
					Term: NewIntegerConstantByValue("123"),
				},
			},
			want: []string{
				"push constant 123",
				"return",
			},
		},
		{
			desc: "加算を伴うExpressionを返すreturn文: return 2;",
			returnStatement: &ReturnStatement{
				StatementKeyword: NewStatementKeyword("return"),
				Expression: &Expression{
					Term: NewIntegerConstantByValue("2"),
					BinaryOpTerms: &BinaryOpTerms{
						Items: []*BinaryOpTerm{
							&BinaryOpTerm{
								BinaryOp: ConstPlus,
								Term:     NewIntegerConstantByValue("3"),
							},
						},
					},
				},
			},
			want: []string{
				"push constant 2",
				"push constant 3",
				"add",
				"return",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.returnStatement.ToCode()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}
