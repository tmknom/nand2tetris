package parsing

import (
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
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.expression.ToCode()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}
