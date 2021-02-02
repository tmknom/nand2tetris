package parsing

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestReturnStatementToCode(t *testing.T) {
	cases := []struct {
		desc            string
		returnStatement *ReturnStatement
		want            []string
	}{
		{
			desc:            "値をなにも返さないreturn文",
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
			desc: "加算を伴うExpressionを返すreturn文",
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
