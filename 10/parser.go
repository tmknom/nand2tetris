package main

import (
	"fmt"
	"github.com/pkg/errors"
)

type Parser struct {
	tokens *Tokens
}

func NewParser(tokens *Tokens) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) excludeParsedTokens() {
	p.tokens = p.tokens.SubList()
}

func (p *Parser) advanceToken() *Token {
	return p.tokens.Advance()
}

func (p *Parser) backwardToken() *Token {
	return p.tokens.Backward()
}

func (p *Parser) readFirstToken() *Token {
	return p.tokens.First()
}

func (p *Parser) Parse() (*Class, error) {
	return p.parseClass()
}

// class Main { ... }
func (p *Parser) parseClass() (*Class, error) {
	class := NewClass()

	keyword := p.advanceToken()
	if err := class.CheckKeyword(keyword); err != nil {
		return nil, err
	}

	className := p.advanceToken()
	if err := class.SetClassName(className); err != nil {
		return nil, err
	}

	openSymbol := p.advanceToken()
	if err := class.CheckOpenSymbol(openSymbol); err != nil {
		return nil, err
	}

	// 閉じカッコは後ろから取得
	closeSymbol := p.backwardToken()
	if err := class.CheckCloseSymbol(closeSymbol); err != nil {
		return nil, err
	}

	// パース済みのclass要素を除外したTokensに更新
	p.excludeParsedTokens()

	classVarDecs, err := p.parseClassVarDecs()
	if err != nil {
		return nil, err
	}
	class.SetClassVarDecs(classVarDecs)

	// TODO SubroutineDec の処理
	//fmt.Printf("Tokens = %+v\n", p.tokens.debug())

	//fmt.Printf("class = %+v\n", class.debug())
	return class, nil
}

type Class struct {
	Keyword       *Keyword
	ClassName     *ClassName
	OpenSymbol    *Symbol
	CloseSymbol   *Symbol
	ClassVarDecs  *ClassVarDecs
	SubroutineDec []*Token
}

func NewClass() *Class {
	return &Class{
		Keyword:       NewKeyword("class"),
		OpenSymbol:    NewSymbol("{"),
		CloseSymbol:   NewSymbol("}"),
		SubroutineDec: []*Token{},
	}
}

func (c *Class) CheckKeyword(token *Token) error {
	return c.Keyword.Check(token)
}

type ClassName struct {
	*Identifier
}

func NewClassName(token *Token) *ClassName {
	return &ClassName{
		Identifier: NewIdentifier("ClassName", token),
	}
}

func (c *Class) SetClassName(token *Token) error {
	className := NewClassName(token)
	if err := className.Check(); err != nil {
		return err
	}

	c.ClassName = className
	return nil
}

func (c *Class) CheckOpenSymbol(token *Token) error {
	return c.OpenSymbol.Check(token)
}

func (c *Class) CheckCloseSymbol(token *Token) error {
	return c.CloseSymbol.Check(token)
}

func (c *Class) SetClassVarDecs(classVarDecs *ClassVarDecs) {
	c.ClassVarDecs = classVarDecs
}

func (c *Class) ToXML() []string {
	result := []string{}
	result = append(result, "<class>")
	result = append(result, c.Keyword.ToXML())
	result = append(result, c.ClassName.ToXML())
	result = append(result, c.OpenSymbol.ToXML())
	result = append(result, c.ClassVarDecs.ToXML()...)
	result = append(result, c.CloseSymbol.ToXML())
	result = append(result, "</class>")
	return result
}

func (c *Class) debug() string {
	return fmt.Sprintf("&Class{\n  Keyword: %s,\n  ClassName: %s},\n  OpenSymbol: %s,\n  CloseSymbol: %s,\n  ClassVarDec: &Tokens{...},\n  SubroutineDec: &Tokens{...}\n}",
		c.Keyword.debug(), c.ClassName.debug(), c.OpenSymbol.debug(), c.CloseSymbol.debug())
}

func (p *Parser) parseClassVarDecs() (*ClassVarDecs, error) {
	classVarDecs := NewClassVarDecs()

	for classVarDecs.hasClassVarDec(p.readFirstToken()) {
		classVarDec := NewClassVarDec()

		keyword := p.advanceToken()
		if err := classVarDec.SetKeyword(keyword); err != nil {
			return nil, err
		}

		varType := p.advanceToken()
		if err := classVarDec.SetVarType(varType); err != nil {
			return nil, err
		}

		varName := p.advanceToken()
		if err := classVarDec.SetFirstVarName(varName); err != nil {
			return nil, err
		}

		for p.readFirstToken().Value == "," {
			comma := p.advanceToken()
			varName := p.advanceToken()
			if err := classVarDec.AddCommaAndVarName(comma, varName); err != nil {
				return nil, err
			}
		}

		endSymbol := p.advanceToken()
		if err := classVarDec.CheckEndSymbol(endSymbol); err != nil {
			return nil, err
		}

		// パースに成功したら要素に追加
		classVarDecs.Add(classVarDec)

		// パース済みのclass要素を除外したTokensに更新
		p.excludeParsedTokens()
	}

	return classVarDecs, nil
}

type ClassVarDecs struct {
	Items []*ClassVarDec
}

func NewClassVarDecs() *ClassVarDecs {
	return &ClassVarDecs{
		Items: []*ClassVarDec{},
	}
}

func (c *ClassVarDecs) Add(item *ClassVarDec) {
	c.Items = append(c.Items, item)
}

func (c *ClassVarDecs) ToXML() []string {
	result := []string{}
	for _, item := range c.Items {
		result = append(result, item.ToXML()...)
	}
	return result
}

func (c *ClassVarDecs) hasClassVarDec(token *Token) bool {
	if token == nil {
		return false
	}
	return token.Value == "static" || token.Value == "field"
}

type ClassVarDec struct {
	Keyword   *Token
	VarType   *Token
	VarNames  *VarNames
	EndSymbol *Symbol
}

func NewClassVarDec() *ClassVarDec {
	return &ClassVarDec{
		VarNames:  NewVarNames(),
		EndSymbol: NewSymbol(";"),
	}
}

func (c *ClassVarDec) SetKeyword(token *Token) error {
	if err := c.checkKeyword(token); err != nil {
		return err
	}

	c.Keyword = token
	return nil
}

func (c *ClassVarDec) checkKeyword(token *Token) error {
	if token.TokenType == TokenKeyword {
		return nil
	}

	message := fmt.Sprintf("Keyword: got = %s", token.debug())
	return errors.New(message)
}

func (c *ClassVarDec) SetVarType(token *Token) error {
	if err := c.checkVarType(token); err != nil {
		return err
	}

	c.VarType = token
	return nil
}

func (c *ClassVarDec) checkVarType(token *Token) error {
	if token.TokenType == TokenIdentifier {
		return nil
	}

	if token.TokenType == TokenKeyword && (token.Value == "int" || token.Value == "char" || token.Value == "boolean") {
		return nil
	}

	message := fmt.Sprintf("VarType: got = %s", token.debug())
	return errors.New(message)
}

func (c *ClassVarDec) SetFirstVarName(token *Token) error {
	return c.VarNames.SetFirst(token)
}

func (c *ClassVarDec) AddCommaAndVarName(comma *Token, varName *Token) error {
	return c.VarNames.AddCommaAndVarName(comma, varName)
}

func (c *ClassVarDec) CheckEndSymbol(token *Token) error {
	return c.EndSymbol.Check(token)
}

func (c *ClassVarDec) ToXML() []string {
	result := []string{}
	result = append(result, "<classVarDec>")
	result = append(result, c.Keyword.ToXML())
	result = append(result, c.VarType.ToXML())
	result = append(result, c.VarNames.ToXML()...)
	result = append(result, c.EndSymbol.ToXML())
	result = append(result, "</classVarDec>")
	return result
}

type VarNames struct {
	First            *Token
	CommaAndVarNames []*CommaAndVarName
}

func NewVarNames() *VarNames {
	return &VarNames{
		CommaAndVarNames: []*CommaAndVarName{},
	}
}

func (v *VarNames) SetFirst(token *Token) error {
	if err := v.checkVarName(token); err != nil {
		return err
	}

	v.First = token
	return nil
}

func (v *VarNames) AddCommaAndVarName(comma *Token, varName *Token) error {
	if err := v.checkComma(comma); err != nil {
		return err
	}

	if err := v.checkVarName(varName); err != nil {
		return err
	}

	v.CommaAndVarNames = append(v.CommaAndVarNames, NewCommaAndVarName(comma, varName))
	return nil
}

func (v *VarNames) checkVarName(token *Token) error {
	if token.TokenType == TokenIdentifier {
		return nil
	}

	message := fmt.Sprintf("VarName: got = %s", token.debug())
	return errors.New(message)
}

func (v *VarNames) checkComma(token *Token) error {
	if token.TokenType == TokenSymbol && token.Value == "," {
		return nil
	}

	message := fmt.Sprintf("Comma: got = %s", token.debug())
	return errors.New(message)
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
	Comma   *Token
	VarName *Token
}

func NewCommaAndVarName(comma *Token, varName *Token) *CommaAndVarName {
	return &CommaAndVarName{Comma: comma, VarName: varName}
}

func (c *CommaAndVarName) ToXML() []string {
	return []string{
		c.Comma.ToXML(),
		c.VarName.ToXML(),
	}
}

type Keyword struct {
	*Token
}

func NewKeyword(value string) *Keyword {
	return &Keyword{
		Token: NewToken(value, TokenKeyword),
	}
}

func (k *Keyword) Check(token *Token) error {
	return token.CheckKeyword(k.Value)
}

type Symbol struct {
	*Token
}

func NewSymbol(value string) *Symbol {
	return &Symbol{
		Token: NewToken(value, TokenSymbol),
	}
}

func (s *Symbol) Check(token *Token) error {
	return token.CheckSymbol(s.Value)
}

type Identifier struct {
	Name string
	*Token
}

func NewIdentifier(name string, token *Token) *Identifier {
	return &Identifier{
		Name:  name,
		Token: token,
	}
}

func (i *Identifier) Check() error {
	return i.CheckTokenType(TokenIdentifier, i.Name)
}
