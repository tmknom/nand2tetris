package parsing

import (
	"../symbol"
	"../token"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestExpressionToCode(t *testing.T) {
	cases := []struct {
		desc       string
		expression *Expression
		want       []string
	}{
		{
			desc: "IntegerConstantをひとつだけ定義",
			expression: &Expression{
				Term: NewIntegerConstantByValue("123"),
			},
			want: []string{
				"push constant 123",
			},
		},
		{
			desc: "ふたつのIntegerConstantを加算",
			expression: &Expression{
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
			want: []string{
				"push constant 2",
				"push constant 3",
				"add",
			},
		},
		{
			desc: "ふたつのIntegerConstantを乗算",
			expression: &Expression{
				Term: NewIntegerConstantByValue("2"),
				BinaryOpTerms: &BinaryOpTerms{
					Items: []*BinaryOpTerm{
						&BinaryOpTerm{
							BinaryOp: ConstAsterisk,
							Term:     NewIntegerConstantByValue("3"),
						},
					},
				},
			},
			want: []string{
				"push constant 2",
				"push constant 3",
				"call Math.multiply 2",
			},
		},
		{
			desc: "カッコを含むIntegerConstantの計算: 1 + (2 * 3)",
			expression: &Expression{
				Term: &IntegerConstant{
					Token: token.NewToken("1", token.TokenIntConst),
				},
				BinaryOpTerms: &BinaryOpTerms{
					Items: []*BinaryOpTerm{
						&BinaryOpTerm{
							BinaryOp: ConstPlus,
							Term: &GroupingExpression{
								Expression: &Expression{
									Term: NewIntegerConstantByValue("2"),
									BinaryOpTerms: &BinaryOpTerms{
										Items: []*BinaryOpTerm{
											&BinaryOpTerm{
												BinaryOp: ConstAsterisk,
												Term:     NewIntegerConstantByValue("3"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: []string{
				"push constant 1",
				"push constant 2",
				"push constant 3",
				"call Math.multiply 2",
				"add",
			},
		},
		{
			desc: "ローカル変数のVarNameをひとつだけ定義",
			expression: &Expression{
				Term: NewVarNameByValue("foo"),
			},
			want: []string{
				"push local 0",
			},
		},
		{
			desc: "引数のVarNameをひとつだけ定義",
			expression: &Expression{
				Term: NewVarNameByValue("bar"),
			},
			want: []string{
				"push argument 0",
			},
		},
		{
			desc: "GroupingExpressionをひとつだけ定義: (6 / 3)",
			expression: &Expression{
				Term: &GroupingExpression{
					Expression: &Expression{
						Term: NewIntegerConstantByValue("6"),
						BinaryOpTerms: &BinaryOpTerms{
							Items: []*BinaryOpTerm{
								&BinaryOpTerm{
									BinaryOp: ConstSlash,
									Term:     NewIntegerConstantByValue("3"),
								},
							},
						},
					},
				},
			},
			want: []string{
				"push constant 6",
				"push constant 3",
				"call Math.divide 2",
			},
		},
		{
			desc: "GroupingExpressionを複数定義: 1 + ( (2 * 3) - 4 )",
			expression: &Expression{
				Term: &IntegerConstant{
					Token: token.NewToken("1", token.TokenIntConst),
				},
				BinaryOpTerms: &BinaryOpTerms{
					Items: []*BinaryOpTerm{
						&BinaryOpTerm{
							BinaryOp: ConstPlus,
							Term: &GroupingExpression{
								Expression: &Expression{
									Term: &GroupingExpression{
										Expression: &Expression{
											Term: &IntegerConstant{
												Token: token.NewToken("2", token.TokenIntConst),
											},
											BinaryOpTerms: &BinaryOpTerms{
												Items: []*BinaryOpTerm{
													&BinaryOpTerm{
														BinaryOp: ConstAsterisk,
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
									BinaryOpTerms: &BinaryOpTerms{
										Items: []*BinaryOpTerm{
											&BinaryOpTerm{
												BinaryOp: ConstMinus,
												Term: &IntegerConstant{
													Token: token.NewToken("4", token.TokenIntConst),
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
			want: []string{
				"push constant 1",
				"push constant 2",
				"push constant 3",
				"call Math.multiply 2",
				"push constant 4",
				"sub",
				"add",
			},
		},
		{
			desc: "UnaryMinusの演算を含む: -123",
			expression: &Expression{
				Term: &UnaryOpTerm{
					UnaryOp: ConstUnaryMinus,
					Term:    NewIntegerConstantByValue("123"),
				},
			},
			want: []string{
				"push constant 123",
				"neg",
			},
		},
		{
			desc: "UnaryTildeの演算を含む: ~(2 - 3)",
			expression: &Expression{
				Term: &UnaryOpTerm{
					UnaryOp: ConstUnaryTilde,
					Term: &GroupingExpression{
						Expression: &Expression{
							Term: NewIntegerConstantByValue("2"),
							BinaryOpTerms: &BinaryOpTerms{
								Items: []*BinaryOpTerm{
									&BinaryOpTerm{
										BinaryOp: ConstMinus,
										Term:     NewIntegerConstantByValue("3"),
									},
								},
							},
						},
					},
				},
			},
			want: []string{
				"push constant 2",
				"push constant 3",
				"sub",
				"not",
			},
		},
		{
			desc: "引数なしのサブルーチン: max()",
			expression: &Expression{
				Term: &SubroutineCall{
					SubroutineCallName: &SubroutineCallName{
						SubroutineName: NewSubroutineNameByValue("max"),
					},
					ExpressionList: NewExpressionList(),
				},
			},
			want: []string{
				"call max 0",
			},
		},
		{
			desc: "引数ありのサブルーチン: max(123, foo)",
			expression: &Expression{
				Term: &SubroutineCall{
					SubroutineCallName: &SubroutineCallName{
						SubroutineName: NewSubroutineNameByValue("max"),
					},
					ExpressionList: &ExpressionList{
						First: &Expression{
							Term: NewIntegerConstantByValue("123"),
						},
						CommaAndExpressions: []*CommaAndExpression{
							NewCommaAndExpression(&Expression{
								Term: NewVarNameByValue("foo"),
							}),
						},
					},
				},
			},
			want: []string{
				"push constant 123",
				"push local 0",
				"call max 2",
			},
		},
		{
			desc: "ビルトインのサブルーチン: Output.printInt(foo)",
			expression: &Expression{
				Term: &SubroutineCall{
					SubroutineCallName: &SubroutineCallName{
						CallerName:     NewCallerNameByValue("Output"),
						SubroutineName: NewSubroutineNameByValue("printInt"),
					},
					ExpressionList: &ExpressionList{
						First: &Expression{
							Term: NewVarNameByValue("foo"),
						},
						CommaAndExpressions: []*CommaAndExpression{},
					},
				},
			},
			want: []string{
				"push local 0",
				"call Output.printInt 1",
			},
		},
	}

	// いろいろ初期化
	SetupTestForToCode()
	// シンボルテーブルのセットアップ
	symbol.GlobalSymbolTables.AddVarSymbol("foo", "int")
	symbol.GlobalSymbolTables.AddArgSymbol("bar", "int")

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.expression.ToCode()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}
