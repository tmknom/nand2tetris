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

func (p *Parser) parseClass() (*Class, error) {
	class := NewClass()

	keyword := p.advanceToken()
	if !class.Keyword.Equals(keyword) {
		return nil, fmt.Errorf("parseClass Keyword: %s", keyword.debug())
	}

	identifier := p.advanceToken()
	if identifier.TokenType != TokenIdentifier {
		return nil, fmt.Errorf("parseClass Identifier: %s", identifier.debug())
	}
	class.SetIdentifier(identifier)

	openSymbol := p.advanceToken()
	if !class.OpenSymbol.Equals(openSymbol) {
		return nil, fmt.Errorf("parseClass OpenSymbol: %s", openSymbol.debug())
	}

	// 閉じカッコは後ろから取得
	closeSymbol := p.backwardToken()
	if !class.CloseSymbol.Equals(closeSymbol) {
		return nil, fmt.Errorf("parseClass CloseSymbol: %s", closeSymbol.debug())
	}

	// パース済みのclass要素を除外したTokensに更新
	p.excludeParsedTokens()

	classVarDecs, err := p.parseClassVarDecs()
	if err != nil {
		return nil, fmt.Errorf("parseClass ClassVarDecs: %+v", err)
	}
	class.SetClassVarDecs(classVarDecs)

	// TODO SubroutineDec の処理
	//fmt.Printf("Tokens = %+v\n", p.tokens.debug())

	//fmt.Printf("class = %+v\n", class.debug())
	return class, nil
}

type Class struct {
	Keyword       *Token
	Identifier    *Token
	OpenSymbol    *Token
	CloseSymbol   *Token
	ClassVarDecs  *ClassVarDecs
	SubroutineDec *Tokens
}

func NewClass() *Class {
	return &Class{
		Keyword:       NewToken("class", TokenKeyword),
		OpenSymbol:    NewToken("{", TokenSymbol),
		CloseSymbol:   NewToken("}", TokenSymbol),
		SubroutineDec: NewTokens(),
	}
}

func (c *Class) ToXML() []string {
	result := []string{}
	result = append(result, "<class>")
	result = append(result, c.Keyword.ToXML())
	result = append(result, c.Identifier.ToXML())
	result = append(result, c.OpenSymbol.ToXML())
	result = append(result, c.ClassVarDecs.ToXML()...)
	result = append(result, c.CloseSymbol.ToXML())
	result = append(result, "</class>")
	return result
}

func (c *Class) SetIdentifier(identifier *Token) {
	c.Identifier = identifier
}

func (c *Class) SetClassVarDecs(classVarDecs *ClassVarDecs) {
	c.ClassVarDecs = classVarDecs
}

func (c *Class) debug() string {
	return fmt.Sprintf("&Class{\n  Keyword: %s,\n  Identifier: %s},\n  OpenSymbol: %s,\n  CloseSymbol: %s,\n  ClassVarDec: &Tokens{...},\n  SubroutineDec: &Tokens{...}\n}",
		c.Keyword.debug(), c.Identifier.debug(), c.OpenSymbol.debug(), c.CloseSymbol.debug())
}

func (p *Parser) parseClassVarDecs() (*ClassVarDecs, error) {
	classVarDecs := NewClassVarDecs()

	for classVarDecs.hasClassVarDec(p.readFirstToken()) {
		classVarDec := NewClassVarDec()

		keyword := p.advanceToken()
		if keyword.TokenType != TokenKeyword {
			return nil, fmt.Errorf("parseClassVarDecs Keyword: %s", keyword.debug())
		}
		classVarDec.SetKeyword(keyword)

		varType := p.advanceToken()
		if varType.TokenType != TokenKeyword && varType.TokenType != TokenIdentifier {
			return nil, fmt.Errorf("parseClassVarDecs VarType: %s", varType.debug())
		}
		classVarDec.SetVarType(varType)

		varName := p.advanceToken()
		if varName.TokenType != TokenIdentifier {
			return nil, fmt.Errorf("parseClassVarDecs VarName: %s", varName.debug())
		}
		classVarDec.AddVarName(varName)

		for p.readFirstToken().Value == "," {
			comma := p.advanceToken()
			classVarDec.AddVarName(comma)

			varName := p.advanceToken()
			if varName.TokenType != TokenIdentifier {
				return nil, fmt.Errorf("parseClassVarDecs VarName loop: %s", varName.debug())
			}
			classVarDec.AddVarName(varName)
		}

		endSymbol := p.advanceToken()
		if !classVarDec.EndSymbol.Equals(endSymbol) {
			return nil, fmt.Errorf("parseClassVarDecs EndSymbol: got = %s, want = %s", endSymbol.debug(), classVarDec.EndSymbol.debug())
		}

		// パースに成功したら要素に追加
		classVarDecs.Add(classVarDec)

		// パース済みのclass要素を除外したTokensに更新
		p.excludeParsedTokens()
	}

	return classVarDecs, nil
}

type ClassVarDecs struct {
	items []*ClassVarDec
}

func NewClassVarDecs() *ClassVarDecs {
	return &ClassVarDecs{
		items: []*ClassVarDec{},
	}
}

func (c *ClassVarDecs) Add(item *ClassVarDec) {
	c.items = append(c.items, item)
}

func (c *ClassVarDecs) ToXML() []string {
	result := []string{}
	for _, item := range c.items {
		result = append(result, item.ToXML()...)
	}
	return result
}

func (c *ClassVarDecs) hasClassVarDec(token *Token) bool {
	return token.Value == "static" || token.Value == "field"
}

type ClassVarDec struct {
	Keyword   *Token
	VarType   *Token
	VarName   *Tokens
	EndSymbol *Token
}

func NewClassVarDec() *ClassVarDec {
	return &ClassVarDec{
		VarName:   NewTokens(),
		EndSymbol: NewToken(";", TokenSymbol),
	}
}

func (c *ClassVarDec) SetKeyword(token *Token) {
	c.Keyword = token
}

func (c *ClassVarDec) SetVarType(token *Token) {
	c.VarType = token
}

func (c *ClassVarDec) AddVarName(token *Token) {
	c.VarName.Add([]*Token{token})
}

func (c *ClassVarDec) ToXML() []string {
	result := []string{}
	result = append(result, "<classVarDec>")
	result = append(result, c.Keyword.ToXML())
	result = append(result, c.VarType.ToXML())

	for _, token := range c.VarName.items {
		result = append(result, token.ToXML())
	}

	result = append(result, c.EndSymbol.ToXML())
	result = append(result, "</classVarDec>")
	return result
}
