package parsing

import (
	"../token"
	"fmt"
	"github.com/pkg/errors"
)

type SubroutineCall struct {
	*SubroutineCallName
	*ExpressionList
	*OpeningRoundBracket
	*ClosingRoundBracket
}

func NewSubroutineCall() *SubroutineCall {
	return &SubroutineCall{
		ExpressionList:      NewExpressionList(),
		OpeningRoundBracket: ConstOpeningRoundBracket,
		ClosingRoundBracket: ConstClosingRoundBracket,
	}
}

func (s *SubroutineCall) SetSubroutineCallName(subroutineCallName *SubroutineCallName) {
	s.SubroutineCallName = subroutineCallName
}

func (s *SubroutineCall) SetExpressionList(expressionList *ExpressionList) {
	s.ExpressionList = expressionList
}

func (s *SubroutineCall) ToXML() []string {
	result := []string{}
	result = append(result, s.SubroutineCallName.ToXML()...)
	result = append(result, s.OpeningRoundBracket.ToXML())
	result = append(result, s.ExpressionList.ToXML()...)
	result = append(result, s.ClosingRoundBracket.ToXML())
	return result
}

func (s *SubroutineCall) TermType() TermType {
	return TermSubroutineCall
}

type SubroutineCallName struct {
	*CallerName
	*Period
	*SubroutineName
}

func NewSubroutineCallName() *SubroutineCallName {
	return &SubroutineCallName{
		Period: ConstPeriod,
	}
}

func (s *SubroutineCallName) SetSubroutineName(token *token.Token) error {
	subroutineName := NewSubroutineName(token)
	if err := subroutineName.Check(); err != nil {
		return err
	}

	s.SubroutineName = subroutineName
	return nil
}

func (s *SubroutineCallName) SetCallerName(token *token.Token) error {
	callerName := NewCallerName(token)
	if err := callerName.Check(); err != nil {
		return err
	}

	s.CallerName = callerName
	return nil
}

func (s *SubroutineCallName) Check() error {
	return s.SubroutineName.Check()
}

func (s *SubroutineCallName) ToXML() []string {
	result := []string{}
	if s.CallerName != nil {
		result = append(result, s.CallerName.ToXML())
		result = append(result, s.Period.ToXML())
	}
	result = append(result, s.SubroutineName.ToXML())
	return result
}

type CallerName struct {
	*Identifier
}

func NewCallerName(token *token.Token) *CallerName {
	return &CallerName{
		Identifier: NewIdentifier("CallerName", token),
	}
}

type ExpressionList struct {
	First               *Expression
	CommaAndExpressions []*CommaAndExpression
}

func NewExpressionList() *ExpressionList {
	return &ExpressionList{
		CommaAndExpressions: []*CommaAndExpression{},
	}
}

func (e *ExpressionList) AddExpression(token *token.Token) error {
	expression := NewExpression(token)
	if err := expression.Check(); err != nil {
		return err
	}

	if e.First == nil {
		e.First = expression
	} else {
		e.CommaAndExpressions = append(e.CommaAndExpressions, NewCommaAndExpression(expression))
	}
	return nil
}

func (e *ExpressionList) ToXML() []string {
	result := []string{}
	result = append(result, "<expressionList>")

	if e.First != nil {
		result = append(result, e.First.ToXML()...)
	}

	for _, item := range e.CommaAndExpressions {
		result = append(result, item.ToXML()...)
	}
	result = append(result, "</expressionList>")
	return result
}

type CommaAndExpression struct {
	*Comma
	*Expression
}

func NewCommaAndExpression(expression *Expression) *CommaAndExpression {
	return &CommaAndExpression{
		Comma:      ConstComma,
		Expression: expression,
	}
}

func (c *CommaAndExpression) ToXML() []string {
	result := []string{}
	result = append(result, c.Comma.ToXML())
	result = append(result, c.Expression.ToXML()...)
	return result
}

type Expression struct {
	*token.Token
}

func NewExpression(token *token.Token) *Expression {
	return &Expression{
		Token: token,
	}
}

func (e *Expression) IsCheck() bool {
	return e.Check() == nil
}

func (e *Expression) Check() error {
	// TODO Expression実装時にちゃんと書く
	if err := NewIdentifier("Expression", e.Token).Check(); err == nil {
		return nil
	}

	return ConstKeywordConstantFactory.Check(e.Token)
}

func (e *Expression) ToXML() []string {
	return []string{e.Token.ToXML()}
}

type KeywordConstantFactory struct{}

func (k *KeywordConstantFactory) Check(token *token.Token) error {
	if _, err := k.Create(token); err != nil {
		return err
	}
	return nil
}

func (k *KeywordConstantFactory) Create(token *token.Token) (Term, error) {
	if err := NewKeyword(token).CheckKeyword(); err != nil {
		return nil, err
	}

	switch token.Value {
	case ConstTrue.Value:
		return ConstTrue, nil
	case ConstFalse.Value:
		return ConstFalse, nil
	case ConstNull.Value:
		return ConstNull, nil
	case ConstThis.Value:
		return ConstThis, nil
	default:
		message := fmt.Sprintf("error create KeywordConstant: got = %s", token.Debug())
		return nil, errors.New(message)
	}
}

var ConstKeywordConstantFactory = &KeywordConstantFactory{}

type KeywordConstant struct {
	*Keyword
}

func NewKeywordConstant(value string) *KeywordConstant {
	return &KeywordConstant{
		Keyword: NewKeywordByValue(value),
	}
}

func (k *KeywordConstant) TermType() TermType {
	return TermKeywordConstant
}

func (k *KeywordConstant) ToXML() []string {
	return []string{k.Token.ToXML()}
}

type TrueKeywordConstant struct {
	*KeywordConstant
}

var ConstTrue = &TrueKeywordConstant{
	KeywordConstant: NewKeywordConstant("true"),
}

type FalseKeywordConstant struct {
	*KeywordConstant
}

var ConstFalse = &FalseKeywordConstant{
	KeywordConstant: NewKeywordConstant("false"),
}

type NullKeywordConstant struct {
	*KeywordConstant
}

var ConstNull = &NullKeywordConstant{
	KeywordConstant: NewKeywordConstant("null"),
}

type ThisKeywordConstant struct {
	*KeywordConstant
}

var ConstThis = &ThisKeywordConstant{
	KeywordConstant: NewKeywordConstant("this"),
}

type StringConstant struct {
	*token.Token
}

func NewStringConstant(token *token.Token) *StringConstant {
	return &StringConstant{
		Token: token,
	}
}

func (s *StringConstant) Check() error {
	return s.Token.CheckStringConstant()
}

func (s *StringConstant) TermType() TermType {
	return TermStringConstant
}

func (s *StringConstant) ToXML() []string {
	return []string{s.Token.ToXML()}
}

type IntegerConstant struct {
	*token.Token
}

func NewIntegerConstant(token *token.Token) *IntegerConstant {
	return &IntegerConstant{
		Token: token,
	}
}

func (i *IntegerConstant) Check() error {
	return i.Token.CheckIntegerConstant()
}

func (i *IntegerConstant) TermType() TermType {
	return TermIntegerConstant
}

func (i *IntegerConstant) ToXML() []string {
	return []string{i.Token.ToXML()}
}

type Term interface {
	TermType() TermType
	ToXML() []string
	Debug() string
}

type TermType int

const (
	_ TermType = iota
	TermIntegerConstant
	TermStringConstant
	TermKeywordConstant
	TermVarName
	TermSubroutineCall
	TermArray          // varName '[' expression ']'
	TermExpression     // '(' expression ')'
	TermWithUnary      // unaryOp term
	TermNotImplemented // TODO TermとExpressionを正しく実装したら消す
)
