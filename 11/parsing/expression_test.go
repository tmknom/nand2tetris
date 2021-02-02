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
				Term: &IntegerConstant{
					Token: token.NewToken("123", token.TokenIntConst),
				},
			},
			want: []string{
				"push constant 123",
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
