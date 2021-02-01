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

func (s *SubroutineCall) Debug() string {
	baseIndent := 0
	indent := baseIndent + 2

	result := "\n"
	result += IndentSprintf(baseIndent, "&SubroutineCall{")
	result += s.SubroutineCallName.Debug(indent)
	result += IndentSprintf(baseIndent, "}")
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

func (s *SubroutineCallName) Debug(baseIndent int) string {
	indent := baseIndent + 2
	result := IndentSprintf(baseIndent, "&SubroutineCallName{")
	if s.CallerName != nil {
		result += s.CallerName.Debug(indent)
	}
	result += s.SubroutineName.Debug(indent)
	result += IndentSprintf(baseIndent, "}")
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

func (c *CallerName) Debug(baseIndent int) string {
	indent := baseIndent + 2
	result := IndentSprintf(baseIndent, "&CallerName{")
	result += IndentSprintf(indent, c.Token.Debug())
	result += IndentSprintf(baseIndent, "}")
	return result
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

func (e *ExpressionList) AddExpression(expression *Expression) {
	if e.First == nil {
		e.First = expression
		return
	}
	e.CommaAndExpressions = append(e.CommaAndExpressions, NewCommaAndExpression(expression))
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

// '(' expression ')'
type GroupingExpression struct {
	*Expression
	*OpeningRoundBracket
	*ClosingRoundBracket
}

func NewGroupingExpression(expression *Expression) *GroupingExpression {
	return &GroupingExpression{
		Expression:          expression,
		OpeningRoundBracket: ConstOpeningRoundBracket,
		ClosingRoundBracket: ConstClosingRoundBracket,
	}
}

func (g *GroupingExpression) TermType() TermType {
	return TermGroupingExpression
}

func (g *GroupingExpression) ToXML() []string {
	result := []string{}
	result = append(result, g.OpeningRoundBracket.ToXML())
	result = append(result, g.Expression.ToXML()...)
	result = append(result, g.ClosingRoundBracket.ToXML())
	return result
}

// varName '[' expression ']'
type Array struct {
	*VarName
	*Expression
	*OpeningSquareBracket
	*ClosingSquareBracket
}

func NewArray(varName *VarName) *Array {
	return &Array{
		VarName:              varName,
		OpeningSquareBracket: ConstOpeningSquareBracket,
		ClosingSquareBracket: ConstClosingSquareBracket,
	}
}

func NewArrayOrError(token *token.Token) (*Array, error) {
	varName := NewVarName(token)
	if err := varName.Check(); err != nil {
		return nil, err
	}
	return NewArray(varName), nil
}

func (a *Array) SetExpression(expression *Expression) {
	a.Expression = expression
}

func (a *Array) TermType() TermType {
	return TermArray
}

func (a *Array) ToXML() []string {
	result := []string{}
	result = append(result, a.VarName.ToXML()...)
	result = append(result, a.OpeningSquareBracket.ToXML())
	result = append(result, a.Expression.ToXML()...)
	result = append(result, a.ClosingSquareBracket.ToXML())
	return result
}

type Expression struct {
	Term
	*BinaryOpTerms
}

func NewExpression(term Term) *Expression {
	return &Expression{
		Term: term,
	}
}

func (e *Expression) SetBinaryOpTerms(binaryOpTerms *BinaryOpTerms) {
	e.BinaryOpTerms = binaryOpTerms
}

func (e *Expression) ToXML() []string {
	result := []string{}
	result = append(result, e.Term.ToXML()...)
	if e.BinaryOpTerms != nil {
		result = append(result, e.BinaryOpTerms.ToXML()...)
	}
	return result
}

type BinaryOpTerms struct {
	Items []*BinaryOpTerm
}

func NewBinaryOpTerms(term Term) *BinaryOpTerms {
	return &BinaryOpTerms{
		Items: []*BinaryOpTerm{},
	}
}

func (b *BinaryOpTerms) ToXML() []string {
	result := []string{}
	for _, binaryTerm := range b.Items {
		result = append(result, binaryTerm.ToXML()...)
	}
	return result
}

type BinaryOpTerm struct {
	BinaryOp
	Term
}

func (b *BinaryOpTerm) ToXML() []string {
	result := []string{}
	result = append(result, b.BinaryOp.ToXML())
	result = append(result, b.Term.ToXML()...)
	return result
}

var ConstBinaryOpFactory = &BinaryOpFactory{}

type BinaryOpFactory struct{}

func (b *BinaryOpFactory) Check(token *token.Token) error {
	if _, err := b.Create(token); err != nil {
		return err
	}
	return nil
}

func (b *BinaryOpFactory) Create(token *token.Token) (BinaryOp, error) {
	if err := NewSymbol(token).CheckSymbol(); err != nil {
		return nil, err
	}

	switch token.Value {
	case ConstVerticalLine.Value:
		return ConstVerticalLine, nil
	default:
		message := fmt.Sprintf("error create BinaryOp: got = %s", token.Debug())
		return nil, errors.New(message)
	}
}

type VerticalLine struct {
	*Symbol
}

var ConstVerticalLine = &VerticalLine{
	Symbol: NewSymbolByValue("|"),
}

func (v *VerticalLine) OpType() BinaryOpType {
	return VerticalLineType
}

type BinaryOp interface {
	OpType() BinaryOpType
	ToXML() string
}

type BinaryOpType int

const (
	_ BinaryOpType = iota
	VerticalLineType
)

type UnaryOpTerm struct {
	UnaryOp
	Term
}

func NewUnaryOpTerm(unaryOp UnaryOp) *UnaryOpTerm {
	return &UnaryOpTerm{
		UnaryOp: unaryOp,
	}
}

func (u *UnaryOpTerm) SetTerm(term Term) {
	u.Term = term
}

func (u *UnaryOpTerm) TermType() TermType {
	return TermUnaryOpTerm
}

func (u *UnaryOpTerm) ToXML() []string {
	result := []string{}
	result = append(result, u.UnaryOp.ToXML())
	result = append(result, u.Term.ToXML()...)
	return result
}

var ConstUnaryOpFactory = &UnaryOpFactory{}

type UnaryOpFactory struct{}

func (u *UnaryOpFactory) Check(token *token.Token) error {
	if _, err := u.Create(token); err != nil {
		return err
	}
	return nil
}

func (u *UnaryOpFactory) Create(token *token.Token) (UnaryOp, error) {
	if err := NewSymbol(token).CheckSymbol(); err != nil {
		return nil, err
	}

	switch token.Value {
	case ConstMinus.Value:
		return ConstMinus, nil
	case ConstTilde.Value:
		return ConstTilde, nil
	default:
		message := fmt.Sprintf("error create UnaryOp: got = %s", token.Debug())
		return nil, errors.New(message)
	}
}

type Minus struct {
	*Symbol
}

var ConstMinus = &Minus{
	Symbol: NewSymbolByValue("-"),
}

func (m *Minus) OpType() UnaryOpType {
	return MinusType
}

type Tilde struct {
	*Symbol
}

var ConstTilde = &Tilde{
	Symbol: NewSymbolByValue("~"),
}

func (t *Tilde) OpType() UnaryOpType {
	return TildeType
}

type UnaryOp interface {
	OpType() UnaryOpType
	ToXML() string
}

type UnaryOpType int

const (
	_ UnaryOpType = iota
	MinusType
	TildeType
)

var ConstKeywordConstantFactory = &KeywordConstantFactory{}

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
	TermArray              // varName '[' expression ']'
	TermGroupingExpression // '(' expression ')'
	TermUnaryOpTerm        // unaryOp term
)
