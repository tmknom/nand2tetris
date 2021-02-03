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

var _ Term = (*SubroutineCall)(nil)

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

func (s *SubroutineCall) ToCode() []string {
	length := s.ExpressionListLength()
	callName := fmt.Sprintf("call %s %d", s.SubroutineCallName.ToCode(), length)

	result := []string{}
	result = append(result, s.ExpressionList.ToCode()...)
	result = append(result, callName)
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

func (s *SubroutineCallName) ToCode() string {
	result := ""
	if s.CallerName != nil {
		result = fmt.Sprintf("%s.", s.CallerName.Value)
	}
	result += s.SubroutineName.Value
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

func (e *ExpressionList) ExpressionListLength() int {
	if e.First == nil {
		return 0
	}
	return 1 + len(e.CommaAndExpressions)
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

func (e *ExpressionList) ToCode() []string {
	result := []string{}
	if e.First != nil {
		result = append(result, e.First.ToCode()...)
	}

	for _, item := range e.CommaAndExpressions {
		result = append(result, item.ToCode()...)
	}
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

var _ Term = (*GroupingExpression)(nil)

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

func (g *GroupingExpression) ToCode() []string {
	result := []string{}
	result = append(result, g.Expression.ToCode()...)
	return result
}

// varName '[' expression ']'
type Array struct {
	*VarName
	*Expression
	*OpeningSquareBracket
	*ClosingSquareBracket
}

var _ Term = (*Array)(nil)

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

func (a *Array) ToCode() []string {
	return []string{"Array_not_implemented"}
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
	result = append(result, "<expression>")
	result = append(result, ConstTermXMLConverter.ToTermXML(e.Term.ToXML()...)...)
	if e.BinaryOpTerms != nil {
		result = append(result, e.BinaryOpTerms.ToXML()...)
	}
	result = append(result, "</expression>")
	return result
}

func (e *Expression) ToCode() []string {
	result := []string{}
	result = append(result, e.Term.ToCode()...)
	if e.BinaryOpTerms != nil {
		result = append(result, e.BinaryOpTerms.ToCode()...)
	}
	return result
}

type BinaryOpTerms struct {
	Items []*BinaryOpTerm
}

func NewBinaryOpTerms() *BinaryOpTerms {
	return &BinaryOpTerms{
		Items: []*BinaryOpTerm{},
	}
}

func (b *BinaryOpTerms) Add(binaryOpTerm *BinaryOpTerm) {
	b.Items = append(b.Items, binaryOpTerm)
}

func (b *BinaryOpTerms) ToXML() []string {
	result := []string{}
	for _, binaryTerm := range b.Items {
		result = append(result, binaryTerm.ToXML()...)
	}
	return result
}

func (b *BinaryOpTerms) ToCode() []string {
	result := []string{}
	for _, item := range b.Items {
		result = append(result, item.ToCode()...)
	}
	return result
}

type BinaryOpTerm struct {
	BinaryOp
	Term
}

func NewBinaryOpTerm(binaryOp BinaryOp, term Term) *BinaryOpTerm {
	return &BinaryOpTerm{
		BinaryOp: binaryOp,
		Term:     term,
	}
}

func (b *BinaryOpTerm) ToXML() []string {
	result := []string{}
	result = append(result, b.BinaryOp.ToXML())
	result = append(result, ConstTermXMLConverter.ToTermXML(b.Term.ToXML()...)...)
	return result
}

func (b *BinaryOpTerm) ToCode() []string {
	result := []string{}
	result = append(result, b.Term.ToCode()...)
	result = append(result, b.BinaryOp.ToCode()...)
	return result
}

var ConstBinaryOpFactory = &BinaryOpFactory{}

type BinaryOpFactory struct{}

func (b *BinaryOpFactory) IsCheck(token *token.Token) bool {
	return b.Check(token) == nil
}

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
	case ConstPlus.Value:
		return ConstPlus, nil
	case ConstMinus.Value:
		return ConstMinus, nil
	case ConstAsterisk.Value:
		return ConstAsterisk, nil
	case ConstSlash.Value:
		return ConstSlash, nil
	case ConstAmpersand.Value:
		return ConstAmpersand, nil
	case ConstVerticalLine.Value:
		return ConstVerticalLine, nil
	case ConstLessThan.Value:
		return ConstLessThan, nil
	case ConstGreaterThan.Value:
		return ConstGreaterThan, nil
	case ConstEquals.Value:
		return ConstEquals, nil
	default:
		message := fmt.Sprintf("error create BinaryOp: got = %s", token.Debug())
		return nil, errors.New(message)
	}
}

type Plus struct {
	*Symbol
}

var _ BinaryOp = (*Plus)(nil)

var ConstPlus = &Plus{
	Symbol: NewSymbolByValue("+"),
}

func (p *Plus) OpType() BinaryOpType {
	return PlusType
}

func (p *Plus) ToCode() []string {
	return []string{"add"}
}

type Minus struct {
	*Symbol
}

var _ BinaryOp = (*Minus)(nil)

var ConstMinus = &Minus{
	Symbol: NewSymbolByValue("-"),
}

func (m *Minus) OpType() BinaryOpType {
	return MinusType
}

func (m *Minus) ToCode() []string {
	return []string{"sub"}
}

type Asterisk struct {
	*Symbol
}

var _ BinaryOp = (*Asterisk)(nil)

var ConstAsterisk = &Asterisk{
	Symbol: NewSymbolByValue("*"),
}

func (a *Asterisk) OpType() BinaryOpType {
	return AsteriskType
}

func (a *Asterisk) ToCode() []string {
	return []string{"call Math.multiply 2"}
}

type Slash struct {
	*Symbol
}

var _ BinaryOp = (*Slash)(nil)

var ConstSlash = &Slash{
	Symbol: NewSymbolByValue("/"),
}

func (s *Slash) OpType() BinaryOpType {
	return SlashType
}

func (s *Slash) ToCode() []string {
	return []string{"call Math.divide 2"}
}

type Ampersand struct {
	*Symbol
}

var _ BinaryOp = (*Ampersand)(nil)

var ConstAmpersand = &Ampersand{
	Symbol: NewSymbolByValue("&"),
}

func (a *Ampersand) OpType() BinaryOpType {
	return AmpersandType
}

func (a *Ampersand) ToCode() []string {
	return []string{"and"}
}

type VerticalLine struct {
	*Symbol
}

var _ BinaryOp = (*VerticalLine)(nil)

var ConstVerticalLine = &VerticalLine{
	Symbol: NewSymbolByValue("|"),
}

func (v *VerticalLine) OpType() BinaryOpType {
	return VerticalLineType
}

func (v *VerticalLine) ToCode() []string {
	return []string{"or"}
}

type LessThan struct {
	*Symbol
}

var _ BinaryOp = (*LessThan)(nil)

var ConstLessThan = &LessThan{
	Symbol: NewSymbolByValue("<"),
}

func (l *LessThan) OpType() BinaryOpType {
	return LessThanType
}

func (l *LessThan) ToCode() []string {
	return []string{"lt"}
}

type GreaterThan struct {
	*Symbol
}

var _ BinaryOp = (*GreaterThan)(nil)

var ConstGreaterThan = &GreaterThan{
	Symbol: NewSymbolByValue(">"),
}

func (g *GreaterThan) OpType() BinaryOpType {
	return GreaterThanType
}

func (g *GreaterThan) ToCode() []string {
	return []string{"gt"}
}

type Equals struct {
	*Symbol
}

var _ BinaryOp = (*Equals)(nil)

var ConstEquals = &Equals{
	Symbol: NewSymbolByValue("="),
}

func (e *Equals) OpType() BinaryOpType {
	return EqualsType
}

func (e *Equals) ToCode() []string {
	return []string{"eq"}
}

type BinaryOp interface {
	OpType() BinaryOpType
	ToXML() string
	ToCode() []string
}

type BinaryOpType int

const (
	_ BinaryOpType = iota
	PlusType
	MinusType
	AsteriskType
	SlashType
	AmpersandType
	VerticalLineType
	LessThanType
	GreaterThanType
	EqualsType
)

type UnaryOpTerm struct {
	UnaryOp
	Term
}

var _ Term = (*UnaryOpTerm)(nil)

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
	result = append(result, ConstTermXMLConverter.ToTermXML(u.Term.ToXML()...)...)
	return result
}

func (u *UnaryOpTerm) ToCode() []string {
	result := []string{}
	result = append(result, u.Term.ToCode()...)
	result = append(result, u.UnaryOp.ToCode()...)
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
	case ConstUnaryMinus.Value:
		return ConstUnaryMinus, nil
	case ConstUnaryTilde.Value:
		return ConstUnaryTilde, nil
	default:
		message := fmt.Sprintf("error create UnaryOp: got = %s", token.Debug())
		return nil, errors.New(message)
	}
}

type UnaryMinus struct {
	*Symbol
}

var _ UnaryOp = (*UnaryMinus)(nil)

var ConstUnaryMinus = &UnaryMinus{
	Symbol: NewSymbolByValue("-"),
}

func (m *UnaryMinus) OpType() UnaryOpType {
	return UnaryMinusType
}

func (t *UnaryMinus) ToCode() []string {
	return []string{"neg"}
}

type UnaryTilde struct {
	*Symbol
}

var _ UnaryOp = (*UnaryTilde)(nil)

var ConstUnaryTilde = &UnaryTilde{
	Symbol: NewSymbolByValue("~"),
}

func (t *UnaryTilde) OpType() UnaryOpType {
	return UnaryTildeType
}

func (t *UnaryTilde) ToCode() []string {
	return []string{"not"}
}

type UnaryOp interface {
	OpType() UnaryOpType
	ToXML() string
	ToCode() []string
}

type UnaryOpType int

const (
	_ UnaryOpType = iota
	UnaryMinusType
	UnaryTildeType
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

var _ Term = (*KeywordConstant)(nil)

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

func (k *KeywordConstant) ToCode() []string {
	return []string{"KeywordConstant_not_implemented"}
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

var _ Term = (*StringConstant)(nil)

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

func (s *StringConstant) ToCode() []string {
	return []string{"StringConstant_not_implemented"}
}

type IntegerConstant struct {
	*token.Token
}

var _ Term = (*IntegerConstant)(nil)

func NewIntegerConstant(token *token.Token) *IntegerConstant {
	return &IntegerConstant{
		Token: token,
	}
}

func NewIntegerConstantByValue(value string) *IntegerConstant {
	return NewIntegerConstant(token.NewToken(value, token.TokenIntConst))
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

func (i *IntegerConstant) ToCode() []string {
	code := fmt.Sprintf("push constant %s", i.Value)
	return []string{code}
}

var ConstTermXMLConverter = &TermXMLConverter{}

type TermXMLConverter struct{}

func (t *TermXMLConverter) ToTermXML(contents ...string) []string {
	result := []string{}
	result = append(result, "<term>")
	result = append(result, contents...)
	result = append(result, "</term>")
	return result
}

type Term interface {
	TermType() TermType
	ToCode() []string
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
