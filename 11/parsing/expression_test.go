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
			desc: "StringConstantをひとつだけ定義",
			expression: &Expression{
				Term: NewStringConstantByValue("Hello"),
			},
			want: []string{
				"push constant 5",          // 文字列長maxLength
				"call String.new 1",        // maxLengthを指定してStringオブジェクトを生成
				"push constant 72",         // 'H'
				"call String.appendChar 2", // 'H' を引数にappendCharを実行
				"push constant 101",        // 'e'
				"call String.appendChar 2", // 'e' を引数にappendCharを実行
				"push constant 108",        // 'l'
				"call String.appendChar 2", // 'l' を引数にappendCharを実行
				"push constant 108",        // 'l'
				"call String.appendChar 2", // 'l' を引数にappendCharを実行
				"push constant 111",        // 'o'
				"call String.appendChar 2", // 'o' を引数にappendCharを実行
			},
		},
		{
			desc: "KeywordConstant「true」の定義",
			expression: &Expression{
				Term: ConstTrue,
			},
			want: []string{
				"push constant 1",
				"neg",
			},
		},
		{
			desc: "KeywordConstant「false」の定義",
			expression: &Expression{
				Term: ConstFalse,
			},
			want: []string{
				"push constant 0",
			},
		},
		{
			desc: "KeywordConstant「null」の定義",
			expression: &Expression{
				Term: ConstNull,
			},
			want: []string{
				"push constant 0",
			},
		},
		{
			desc: "KeywordConstant「this」の定義",
			expression: &Expression{
				Term: ConstThis,
			},
			want: []string{
				"push pointer 0",
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
			desc: "クラスのフィールド変数のVarNameをひとつだけ定義",
			expression: &Expression{
				Term: NewVarNameByValue("fieldA"),
			},
			want: []string{
				"push this 0",
			},
		},
		{
			desc: "クラスのStatic変数のVarNameをひとつだけ定義",
			expression: &Expression{
				Term: NewVarNameByValue("staticA"),
			},
			want: []string{
				"push static 0",
			},
		},
		{
			desc: "ローカル変数のArray要素の定義",
			expression: &Expression{
				Term: &Array{
					VarName: NewVarNameByValue("array"),
					Expression: &Expression{
						Term: NewIntegerConstantByValue("3"),
					},
				},
			},
			want: []string{
				"push local 2",    // 配列のアドレス
				"push constant 3", // 配列の添字
				"add",             // 配列要素のアドレスの算出
				"pop pointer 1",   // thatにアドレスをセット
				"push that 0",     // 配列要素に値を代入
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
						ClassName:      NewClassNameByValue("Square"),
					},
					ExpressionList: NewExpressionList(),
				},
			},
			want: []string{
				"push pointer 0",
				"call Square.max 1",
			},
		},
		{
			desc: "引数ありのサブルーチン: max(123, foo)",
			expression: &Expression{
				Term: &SubroutineCall{
					SubroutineCallName: &SubroutineCallName{
						SubroutineName: NewSubroutineNameByValue("max"),
						ClassName:      NewClassNameByValue("Square"),
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
				"push pointer 0",
				"call Square.max 3",
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
		{
			desc: "引数なしのオブジェクトのメソッド呼び出し: obj.run()",
			expression: &Expression{
				Term: &SubroutineCall{
					SubroutineCallName: &SubroutineCallName{
						CallerName:     NewCallerNameByValue("obj"),
						SubroutineName: NewSubroutineNameByValue("run"),
					},
					ExpressionList: NewExpressionList(),
				},
			},
			want: []string{
				"push local 1", // 隠れ引数として呼び出したオブジェクトのアドレスをpush
				"call Square.run 1",
			},
		},
		{
			desc: "引数ありのオブジェクトのメソッド呼び出し: obj.run(123, foo)",
			expression: &Expression{
				Term: &SubroutineCall{
					SubroutineCallName: &SubroutineCallName{
						CallerName:     NewCallerNameByValue("obj"),
						SubroutineName: NewSubroutineNameByValue("run"),
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
				"push local 1", // 隠れ引数として呼び出したオブジェクトのアドレスをpush
				"call Square.run 3",
			},
		},
	}

	// いろいろ初期化
	SetupTestForToCode()
	// シンボルテーブルのセットアップ
	symbol.GlobalSymbolTables.AddVarSymbol("foo", "int")
	symbol.GlobalSymbolTables.AddVarSymbol("obj", "Square")
	symbol.GlobalSymbolTables.AddVarSymbol("array", "Array")
	symbol.GlobalSymbolTables.AddArgSymbol("bar", "int")
	symbol.GlobalSymbolTables.AddFieldSymbol("fieldA", "int")
	symbol.GlobalSymbolTables.AddStaticSymbol("staticA", "int")

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.expression.ToCode()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}
