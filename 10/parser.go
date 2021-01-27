package main

import (
	"fmt"
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

func (p *Parser) Parse() (*Class, error) {
	return p.parseClass()
}

func (p *Parser) parseClass() (*Class, error) {
	class := NewClass()

	keyword := p.tokens.Advance()
	if !class.Keyword.Equals(keyword) {
		return nil, fmt.Errorf("parseClass keyword: %s", keyword.debug())
	}

	identifier := p.tokens.Advance()
	if identifier.TokenType != TokenIdentifier {
		return nil, fmt.Errorf("parseClass identifier: %s", identifier.debug())
	}
	class.SetIdentifier(identifier)

	openSymbol := p.tokens.Advance()
	if !class.OpenSymbol.Equals(openSymbol) {
		return nil, fmt.Errorf("parseClass openSymbol: %s", openSymbol.debug())
	}

	// 閉じカッコは後ろから取得
	closeSymbol := p.tokens.Backward()
	if !class.CloseSymbol.Equals(closeSymbol) {
		return nil, fmt.Errorf("parseClass closeSymbol: %s", closeSymbol.debug())
	}

	// パース済みのclass要素を除外したTokensに更新
	p.excludeParsedTokens()
	//fmt.Printf("Tokens = %+v\n", p.tokens.debug())

	// TODO ClassVarDec の処理

	// TODO SubroutineDec の処理

	//fmt.Printf("class = %+v\n", class.debug())
	return class, nil
}

type Class struct {
	Keyword       *Token
	Identifier    *Token
	OpenSymbol    *Token
	CloseSymbol   *Token
	ClassVarDec   *Tokens
	SubroutineDec *Tokens
}

func NewClass() *Class {
	return &Class{
		Keyword:       NewToken("class", TokenKeyword),
		OpenSymbol:    NewToken("{", TokenSymbol),
		CloseSymbol:   NewToken("}", TokenSymbol),
		ClassVarDec:   NewTokens(),
		SubroutineDec: NewTokens(),
	}
}

func (c *Class) ToXML() []string {
	const indent = "  "

	result := []string{}
	result = append(result, "<class>")
	result = append(result, indent+c.Keyword.ToXML())
	result = append(result, indent+c.Identifier.ToXML())
	result = append(result, indent+c.OpenSymbol.ToXML())
	result = append(result, indent+c.CloseSymbol.ToXML())
	result = append(result, "</class>")
	return result
}

func (c *Class) SetIdentifier(identifier *Token) {
	c.Identifier = identifier
}

func (c *Class) debug() string {
	return fmt.Sprintf("&Class{\n  Keyword: %s,\n  Identifier: %s},\n  OpenSymbol: %s,\n  CloseSymbol: %s,\n  ClassVarDec: &Tokens{...},\n  SubroutineDec: &Tokens{...}\n}",
		c.Keyword.debug(), c.Identifier.debug(), c.OpenSymbol.debug(), c.CloseSymbol.debug())
}
