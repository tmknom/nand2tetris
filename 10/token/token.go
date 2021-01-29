package token

import (
	"fmt"
	"github.com/pkg/errors"
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

func (t *Tokens) Debug() string {
	result := "&Tokens{\n"
	for i, item := range t.items {
		result += fmt.Sprintf("    [%d] = %s\n", i, item.Debug())
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

func (t *Token) CheckKeywordValue(expected ...string) error {
	if err := t.CheckValue("Keyword", expected...); err != nil {
		return err
	}

	return t.CheckKeyword()
}

func (t *Token) CheckSymbolValue(expected string) error {
	if err := t.CheckValue("Symbol", expected); err != nil {
		return err
	}

	return t.CheckSymbol()
}

func (t *Token) CheckValue(tokenTypeString string, expected ...string) error {
	for _, value := range expected {
		if t.Value == value {
			return nil
		}
	}

	message := fmt.Sprintf("%s expected values %v: got = %s", tokenTypeString, expected, t.Debug())
	return errors.New(message)
}

func (t *Token) CheckKeyword() error {
	tokenName := fmt.Sprintf("Keyword '%s'", t.Value)
	return t.CheckTokenType(TokenKeyword, tokenName)
}

func (t *Token) CheckSymbol() error {
	tokenName := fmt.Sprintf("Symbol '%s'", t.Value)
	return t.CheckTokenType(TokenSymbol, tokenName)
}

func (t *Token) CheckIdentifier() error {
	tokenName := fmt.Sprintf("Identifier '%s'", t.Value)
	return t.CheckTokenType(TokenIdentifier, tokenName)
}

func (t *Token) CheckTokenType(tokenType TokenType, tokenName string) error {
	if t.TokenType == tokenType {
		return nil
	}

	message := fmt.Sprintf("%s: got = %s", tokenName, t.Debug())
	return errors.New(message)
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

func (t *Token) Debug() string {
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