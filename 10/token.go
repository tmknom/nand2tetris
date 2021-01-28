package main

import (
	"fmt"
)

type Tokens struct {
	items     []*Token
	headIndex int
	tailIndex int
}

func NewTokens() *Tokens {
	return &Tokens{items: []*Token{}, headIndex: 0}
}

func (t *Tokens) Add(items []*Token) {
	for _, item := range items {
		t.items = append(t.items, item)
	}
	t.setupIndex()
}

func (t *Tokens) Advance() *Token {
	t.headIndex += 1
	return t.items[t.headIndex-1]
}

func (t *Tokens) Backward() *Token {
	t.tailIndex -= 1
	return t.items[t.tailIndex+1]
}

func (t *Tokens) First() *Token {
	if len(t.items) > t.headIndex {
		return t.items[t.headIndex]
	}
	return nil
}

func (t *Tokens) SubList() *Tokens {
	tokens := NewTokens()
	tokens.items = t.items[t.headIndex : t.tailIndex+1]
	tokens.setupIndex()
	return tokens
}

func (t *Tokens) setupIndex() {
	t.headIndex = 0
	t.tailIndex = len(t.items) - 1
}

func (t *Tokens) ToXML() []string {
	result := []string{"<tokens>"}
	for _, item := range t.items {
		result = append(result, item.ToXML())
	}
	result = append(result, "</tokens>")
	return result
}

func (t *Tokens) debug() string {
	result := "&Tokens{\n"
	for i, item := range t.items {
		result += fmt.Sprintf("    [%d] = %s\n", i, item.debug())
	}
	result += "}\n"
	return result
}

type Token struct {
	Value     string
	TokenType TokenType
}

type TokenType int

const (
	_ TokenType = iota
	TokenKeyword
	TokenSymbol
	TokenIntConst
	TokenStringConst
	TokenIdentifier
)

func NewToken(value string, tokenType TokenType) *Token {
	return &Token{Value: value, TokenType: tokenType}
}

func (t *Token) Equals(other *Token) bool {
	return t.Value == other.Value && t.TokenType == other.TokenType
}

func (t *Token) ToXML() string {
	value := t.Value
	if t.TokenType == TokenSymbol {
		if encoded, ok := encodeChars[t.Value]; ok {
			value = encoded
		}
	}

	return fmt.Sprintf("<%s> %s </%s>", t.tokenTypeString(), value, t.tokenTypeString())
}

func (t *Token) debug() string {
	return fmt.Sprintf("&Token{Value: '%s', TokenType: %s}", t.Value, t.tokenTypeString())
}

func (t *Token) tokenTypeString() string {
	switch t.TokenType {
	case TokenKeyword:
		return "keyword"
	case TokenSymbol:
		return "symbol"
	case TokenIntConst:
		return "integerConstant"
	case TokenStringConst:
		return "stringConstant"
	case TokenIdentifier:
		return "identifier"
	default:
		return "invalid"
	}
}

var encodeChars = map[string]string{
	"<": "&lt;",
	">": "&gt;",
	"&": "&amp;",
}
