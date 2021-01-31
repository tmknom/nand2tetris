package token

import (
	"fmt"
	"github.com/pkg/errors"
)

type Tokens struct {
	Items     []*Token
	HeadIndex int
	TailIndex int
}

func NewTokens() *Tokens {
	return &Tokens{Items: []*Token{}, HeadIndex: 0}
}

func (t *Tokens) Add(items []*Token) {
	for _, item := range items {
		t.Items = append(t.Items, item)
	}
	t.setupIndex()
}

func (t *Tokens) Advance() *Token {
	t.HeadIndex += 1
	return t.Items[t.HeadIndex-1]
}

func (t *Tokens) Backward() *Token {
	t.TailIndex -= 1
	return t.Items[t.TailIndex+1]
}

func (t *Tokens) First() *Token {
	if len(t.Items) > t.HeadIndex {
		return t.Items[t.HeadIndex]
	}
	return nil
}

func (t *Tokens) Second() *Token {
	if len(t.Items) > t.HeadIndex+1 {
		return t.Items[t.HeadIndex+1]
	}
	return nil
}

func (t *Tokens) SubList() *Tokens {
	tokens := NewTokens()
	tokens.Items = t.Items[t.HeadIndex : t.TailIndex+1]
	tokens.setupIndex()
	return tokens
}

func (t *Tokens) setupIndex() {
	t.HeadIndex = 0
	t.TailIndex = len(t.Items) - 1
}

func (t *Tokens) ToXML() []string {
	result := []string{"<tokens>"}
	for _, item := range t.Items {
		result = append(result, item.ToXML())
	}
	result = append(result, "</tokens>")
	return result
}

func (t *Tokens) Debug() string {
	result := "&Tokens{\n"
	for i, item := range t.Items {
		mark := ""
		if i == t.HeadIndex {
			mark = "<==============="
		}
		result += fmt.Sprintf("    [%d] = %s %s\n", i, item.Debug(), mark)
	}
	result += "}\n"
	result += fmt.Sprintf("HeadIndex = %d\n", t.HeadIndex)
	return result
}

func (t *Tokens) DebugForError() string {
	const indexSize = 15
	start := t.HeadIndex
	if t.HeadIndex > indexSize {
		start = t.HeadIndex - indexSize
	}
	end := len(t.Items[start:])
	if end > indexSize*2 {
		end = t.HeadIndex + indexSize
	}

	result := "&Tokens{\n"
	for i, item := range t.Items[start:end] {
		index := i + start
		mark := ""
		if index == t.HeadIndex {
			mark = "<==============="
		}
		result += fmt.Sprintf("    [%d] = %s %s\n", index, item.Debug(), mark)
	}
	result += "}\n"
	result += fmt.Sprintf("HeadIndex = %d\n", t.HeadIndex)
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

func (t *Token) CheckIntegerConstant() error {
	tokenName := fmt.Sprintf("IntegerConstant '%s'", t.Value)
	return t.CheckTokenType(TokenIntConst, tokenName)
}

func (t *Token) CheckStringConstant() error {
	tokenName := fmt.Sprintf("StringConstant '%s'", t.Value)
	return t.CheckTokenType(TokenStringConst, tokenName)
}

func (t *Token) CheckTokenType(tokenType TokenType, tokenName string) error {
	if t.TokenType == tokenType {
		return nil
	}

	message := fmt.Sprintf("error TokenType: expected = %s: got = %s", tokenName, t.Debug())
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
