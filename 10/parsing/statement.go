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

type LetStatement struct {
	*StatementKeyword
	*VarName
	*Array
	*Equal
	*Expression
	*Semicolon
}

var _ Statement = (*LetStatement)(nil)

func NewLetStatement() *LetStatement {
	return &LetStatement{
		StatementKeyword: NewStatementKeyword("let"),
		Equal:            ConstEqual,
		Semicolon:        ConstSemicolon,
	}
}

func (l *LetStatement) SetVarName(token *token.Token) error {
	varName, err := NewVarNameOrError(token)
	if err != nil {
		return err
	}

	l.VarName = varName
	return nil
}

func (l *LetStatement) SetArray(array *Array) {
	l.Array = array
}

func (l *LetStatement) SetExpression(expression *Expression) {
	l.Expression = expression
}

func (l *LetStatement) ToXML() []string {
	result := []string{}
	result = append(result, l.OpenTag())
	result = append(result, l.StatementKeyword.ToXML())

	if l.VarName != nil {
		result = append(result, l.VarName.ToXML()...)
	}
	if l.Array != nil {
		result = append(result, l.Array.ToXML()...)
	}

	result = append(result, l.Equal.ToXML())
	result = append(result, l.Expression.ToXML()...)
	result = append(result, l.Semicolon.ToXML())
	result = append(result, l.CloseTag())
	return result
}

type IfStatement struct {
	*StatementKeyword
	*Expression
	*Statements
	*ElseBlock
	*OpeningRoundBracket
	*ClosingRoundBracket
	*OpeningCurlyBracket
	*ClosingCurlyBracket
}

var _ Statement = (*IfStatement)(nil)

func NewIfStatement() *IfStatement {
	return &IfStatement{
		StatementKeyword:    NewStatementKeyword("if"),
		OpeningRoundBracket: ConstOpeningRoundBracket,
		ClosingRoundBracket: ConstClosingRoundBracket,
		OpeningCurlyBracket: ConstOpeningCurlyBracket,
		ClosingCurlyBracket: ConstClosingCurlyBracket,
	}
}

func (i *IfStatement) SetExpression(expression *Expression) {
	i.Expression = expression
}

func (i *IfStatement) SetStatements(statements *Statements) {
	i.Statements = statements
}

func (i *IfStatement) SetElseBlock(elseBlock *ElseBlock) {
	i.ElseBlock = elseBlock
}

func (i *IfStatement) ToXML() []string {
	result := []string{}
	result = append(result, i.OpenTag())
	result = append(result, i.StatementKeyword.ToXML())
	result = append(result, i.OpeningRoundBracket.ToXML())
	result = append(result, i.Expression.ToXML()...)
	result = append(result, i.ClosingRoundBracket.ToXML())
	result = append(result, i.OpeningCurlyBracket.ToXML())
	result = append(result, i.Statements.ToXML()...)
	result = append(result, i.ClosingCurlyBracket.ToXML())

	if i.ElseBlock != nil {
		result = append(result, i.ElseBlock.ToXML()...)
	}

	result = append(result, i.CloseTag())
	return result
}

type ElseBlock struct {
	*Keyword
	*Statements
	*OpeningCurlyBracket
	*ClosingCurlyBracket
}

func NewElseBlock() *ElseBlock {
	return &ElseBlock{
		Keyword:             NewKeywordByValue("else"),
		OpeningCurlyBracket: ConstOpeningCurlyBracket,
		ClosingCurlyBracket: ConstClosingCurlyBracket,
	}
}

func (e *ElseBlock) SetStatements(statements *Statements) {
	e.Statements = statements
}

func (e *ElseBlock) CheckElseKeyword(token *token.Token) error {
	return NewKeyword(token).Check(e.Keyword.Value)
}

func (e *ElseBlock) ToXML() []string {
	result := []string{}
	result = append(result, e.Keyword.ToXML())
	result = append(result, e.OpeningCurlyBracket.ToXML())
	result = append(result, e.Statements.ToXML()...)
	result = append(result, e.ClosingCurlyBracket.ToXML())
	return result
}

type WhileStatement struct {
	*StatementKeyword
	*Expression
	*Statements
	*OpeningRoundBracket
	*ClosingRoundBracket
	*OpeningCurlyBracket
	*ClosingCurlyBracket
}

var _ Statement = (*WhileStatement)(nil)

func NewWhileStatement() *WhileStatement {
	return &WhileStatement{
		StatementKeyword:    NewStatementKeyword("while"),
		OpeningRoundBracket: ConstOpeningRoundBracket,
		ClosingRoundBracket: ConstClosingRoundBracket,
		OpeningCurlyBracket: ConstOpeningCurlyBracket,
		ClosingCurlyBracket: ConstClosingCurlyBracket,
	}
}
func (w *WhileStatement) SetExpression(expression *Expression) {
	w.Expression = expression
}

func (w *WhileStatement) SetStatements(statements *Statements) {
	w.Statements = statements
}

func (w *WhileStatement) ToXML() []string {
	result := []string{}
	result = append(result, w.OpenTag())
	result = append(result, w.Keyword.ToXML())
	result = append(result, w.OpeningRoundBracket.ToXML())
	result = append(result, w.Expression.ToXML()...)
	result = append(result, w.ClosingRoundBracket.ToXML())
	result = append(result, w.OpeningCurlyBracket.ToXML())
	result = append(result, w.Statements.ToXML()...)
	result = append(result, w.ClosingCurlyBracket.ToXML())
	result = append(result, w.CloseTag())
	return result
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
	result = append(result, d.StatementKeyword.ToXML())
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

func (r *ReturnStatement) SetExpression(expression *Expression) {
	r.Expression = expression
}

func (r *ReturnStatement) ToXML() []string {
	result := []string{}
	result = append(result, r.OpenTag())
	result = append(result, r.StatementKeyword.ToXML())

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

func (s *StatementKeyword) CheckKeyword(token *token.Token) error {
	return NewKeyword(token).Check(s.Keyword.Value)
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
