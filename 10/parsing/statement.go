package parsing

import "../token"

type Statements struct {
	Items []Statement
}

func NewStatements() *Statements {
	return &Statements{
		Items: []Statement{},
	}
}

func (s *Statements) AddStatement(item Statement) {
	s.Items = append(s.Items, item)
}

func (s *Statements) ToXML() []string {
	result := []string{}
	result = append(result, "<statements>")
	for _, item := range s.Items {
		result = append(result, item.ToXML()...)
	}
	result = append(result, "</statements>")
	return result
}

func (s *Statements) IsStatementKeyword(token *token.Token) bool {
	if token == nil {
		return false
	}

	keywords := []string{
		"let",
		"if",
		"while",
		"do",
		"return",
	}

	err := NewKeyword(token).CheckKeywordValue(keywords...)
	return err == nil
}

type ReturnStatement struct {
	*Keyword
	*Expression
	*Semicolon
}

func NewReturnStatement() *ReturnStatement {
	return &ReturnStatement{
		Keyword:   NewKeywordByValue("return"),
		Semicolon: ConstSemicolon,
	}
}

func (r *ReturnStatement) SetExpression(token *token.Token) error {
	expression := NewExpression(token)
	if err := expression.Check(); err != nil {
		return err
	}

	r.Expression = expression
	return nil
}

func (r *ReturnStatement) ToXML() []string {
	result := []string{}
	result = append(result, "<returnStatement>")
	result = append(result, r.Keyword.ToXML())

	if r.Expression != nil {
		result = append(result, r.Expression.ToXML())
	}

	result = append(result, r.Semicolon.ToXML())
	result = append(result, "</returnStatement>")
	return result
}

func (r *ReturnStatement) Type() string {
	return r.Keyword.Value
}

type Statement interface {
	Type() string
	ToXML() []string
}

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
