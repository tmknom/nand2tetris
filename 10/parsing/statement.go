package parsing

import (
	"../token"
	"fmt"
)

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

type DoStatement struct {
	*StatementKeyword
	*SubroutineCall
	*Semicolon
}

var _ Statement = (*DoStatement)(nil)

func NewDoStatement() *DoStatement {
	return &DoStatement{
		StatementKeyword: NewStatementKeyword("do"),
		Semicolon:        ConstSemicolon,
	}
}

func (d *DoStatement) SetSubroutineCall(subroutineCall *SubroutineCall) {
	d.SubroutineCall = subroutineCall
}

func (d *DoStatement) ToXML() []string {
	result := []string{}
	result = append(result, d.OpenTag())
	result = append(result, d.Keyword.ToXML())
	result = append(result, d.SubroutineCall.ToXML()...)
	result = append(result, d.Semicolon.ToXML())
	result = append(result, d.CloseTag())
	return result
}

type ReturnStatement struct {
	*StatementKeyword
	*Expression
	*Semicolon
}

var _ Statement = (*ReturnStatement)(nil)

func NewReturnStatement() *ReturnStatement {
	return &ReturnStatement{
		StatementKeyword: NewStatementKeyword("return"),
		Semicolon:        ConstSemicolon,
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
	result = append(result, r.OpenTag())
	result = append(result, r.Keyword.ToXML())

	if r.Expression != nil {
		result = append(result, r.Expression.ToXML()...)
	}

	result = append(result, r.Semicolon.ToXML())
	result = append(result, r.CloseTag())
	return result
}

type StatementKeyword struct {
	*Keyword
}

func NewStatementKeyword(value string) *StatementKeyword {
	return &StatementKeyword{
		Keyword: NewKeywordByValue(value),
	}
}

func (s *StatementKeyword) OpenTag() string {
	return fmt.Sprintf("<%sStatement>", s.Keyword.Value)
}

func (s *StatementKeyword) CloseTag() string {
	return fmt.Sprintf("</%sStatement>", s.Keyword.Value)
}

type Statement interface {
	ToXML() []string
	OpenTag() string
	CloseTag() string
}
