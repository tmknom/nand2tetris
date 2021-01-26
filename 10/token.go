package main

type Tokens struct {
	items []*Token
}

func NewTokens() *Tokens {
	return &Tokens{items: []*Token{}}
}

func (t *Tokens) Add(items []*Token) {
	for _, item := range items {
		t.items = append(t.items, item)
	}
}

type Token struct {
	value     string
	tokenType TokenType
}

type TokenType int

const (
	TokenKeyword TokenType = iota
	TokenSymbol
	TokenIdentifier
	TokenIntConst
	TokenStringConst
)

func NewToken(value string, tokenType TokenType) *Token {
	return &Token{value: value, tokenType: tokenType}
}
