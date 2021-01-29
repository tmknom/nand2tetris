package parsing

import "../token"

type Expression struct {
	*token.Token
}

func NewExpression(token *token.Token) *Expression {
	return &Expression{
		Token: token,
	}
}

func (e *Expression) Check() error {
	// TODO Expression実装時にちゃんと書く
	return nil
}
