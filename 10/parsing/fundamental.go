package parsing

import (
	"../token"
	"fmt"
	"github.com/pkg/errors"
)

type VarType struct {
	*token.Token
}

func NewVarType(token *token.Token) *VarType {
	return &VarType{
		Token: token,
	}
}

func (v *VarType) IsCheck() bool {
	return v.Check() == nil
}

func (v *VarType) Check() error {
	if err := v.CheckIdentifier(); err == nil {
		return nil
	}

	expected := []string{"int", "char", "boolean"}
	if err := v.CheckKeywordValue(expected...); err == nil {
		return nil
	}

	message := fmt.Sprintf("VarType: got = %s", v.Debug())
	return errors.New(message)
}

type VarNames struct {
	First            *VarName
	CommaAndVarNames []*CommaAndVarName
}

func NewVarNames() *VarNames {
	return &VarNames{
		CommaAndVarNames: []*CommaAndVarName{},
	}
}

func (v *VarNames) SetFirstVarName(token *token.Token) error {
	varName := NewVarName(token)
	if err := varName.Check(); err != nil {
		return err
	}

	v.First = varName
	return nil
}

func (v *VarNames) AddCommaAndVarName(commaToken *token.Token, varNameToken *token.Token) error {
	if err := ConstComma.Check(commaToken); err != nil {
		return err
	}

	varName := NewVarName(varNameToken)
	if err := varName.Check(); err != nil {
		return err
	}

	v.CommaAndVarNames = append(v.CommaAndVarNames, NewCommaAndVarName(varName))
	return nil
}

func (v *VarNames) ToXML() []string {
	result := []string{}
	result = append(result, v.First.ToXML())
	for _, commaAndVarName := range v.CommaAndVarNames {
		result = append(result, commaAndVarName.ToXML()...)
	}
	return result
}

type CommaAndVarName struct {
	*Comma
	*VarName
}

func NewCommaAndVarName(varName *VarName) *CommaAndVarName {
	return &CommaAndVarName{
		Comma:   ConstComma,
		VarName: varName,
	}
}

func NewCommaAndVarNameByValue(value string) *CommaAndVarName {
	return NewCommaAndVarName(NewVarNameByValue(value))
}

func (c *CommaAndVarName) ToXML() []string {
	return []string{
		c.Comma.ToXML(),
		c.VarName.ToXML(),
	}
}

type SubroutineName struct {
	*Identifier
}

func NewSubroutineName(token *token.Token) *SubroutineName {
	return &SubroutineName{
		Identifier: NewIdentifier("SubroutineName", token),
	}
}

func NewSubroutineNameByValue(value string) *SubroutineName {
	return NewSubroutineName(token.NewToken(value, token.TokenIdentifier))
}

type VarName struct {
	*Identifier
}

func NewVarName(token *token.Token) *VarName {
	return &VarName{
		Identifier: NewIdentifier("VarName", token),
	}
}

func NewVarNameByValue(value string) *VarName {
	return NewVarName(token.NewToken(value, token.TokenIdentifier))
}

type Keyword struct {
	*token.Token
}

func NewKeyword(token *token.Token) *Keyword {
	return &Keyword{
		Token: token,
	}
}

func NewKeywordByValue(value string) *Keyword {
	return NewKeyword(token.NewToken(value, token.TokenKeyword))
}

func (k *Keyword) Check(expected ...string) error {
	return k.CheckKeywordValue(expected...)
}

type Symbol struct {
	*token.Token
}

func NewSymbol(token *token.Token) *Symbol {
	return &Symbol{
		Token: token,
	}
}

func NewSymbolByValue(value string) *Symbol {
	return NewSymbol(token.NewToken(value, token.TokenSymbol))
}

func (s *Symbol) IsCheck(token *token.Token) bool {
	return s.Check(token) == nil
}

func (s *Symbol) Check(token *token.Token) error {
	return NewSymbol(token).CheckSymbolValue(s.Value)
}

type Identifier struct {
	Name string
	*token.Token
}

func NewIdentifier(name string, token *token.Token) *Identifier {
	return &Identifier{
		Name:  name,
		Token: token,
	}
}

func (i *Identifier) Check() error {
	return i.CheckTokenType(token.TokenIdentifier, "Identifier "+i.Name)
}

// よく使われるシンボル
type OpeningCurlyBracket struct {
	*Symbol
}

var ConstOpeningCurlyBracket = &OpeningCurlyBracket{
	Symbol: NewSymbolByValue("{"),
}

type ClosingCurlyBracket struct {
	*Symbol
}

var ConstClosingCurlyBracket = &ClosingCurlyBracket{
	Symbol: NewSymbolByValue("}"),
}

type OpeningRoundBracket struct {
	*Symbol
}

var ConstOpeningRoundBracket = &OpeningRoundBracket{
	Symbol: NewSymbolByValue("("),
}

type ClosingRoundBracket struct {
	*Symbol
}

var ConstClosingRoundBracket = &ClosingRoundBracket{
	Symbol: NewSymbolByValue(")"),
}

type OpeningSquareBracket struct {
	*Symbol
}

var ConstOpeningSquareBracket = &OpeningSquareBracket{
	Symbol: NewSymbolByValue("["),
}

type ClosingSquareBracket struct {
	*Symbol
}

var ConstClosingSquareBracket = &ClosingSquareBracket{
	Symbol: NewSymbolByValue("]"),
}

type Comma struct {
	*Symbol
}

var ConstComma = &Comma{
	Symbol: NewSymbolByValue(","),
}

type Period struct {
	*Symbol
}

var ConstPeriod = &Period{
	Symbol: NewSymbolByValue("."),
}

type Semicolon struct {
	*Symbol
}

var ConstSemicolon = &Semicolon{
	Symbol: NewSymbolByValue(";"),
}

type Equal struct {
	*Symbol
}

var ConstEqual = &Equal{
	Symbol: NewSymbolByValue("="),
}
