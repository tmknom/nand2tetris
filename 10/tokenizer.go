package main

import (
	"fmt"
	"strings"
)

type Tokenizer struct {
	lines  []string
	tokens *Tokens
}

func NewTokenizer(lines []string) *Tokenizer {
	tokens := NewTokens()
	return &Tokenizer{lines: lines, tokens: tokens}
}

func (t *Tokenizer) Tokenize() {
	for _, line := range t.lines {
		items := t.tokenizeLine(line)
		t.tokens.Add(items)
	}

	for _, item := range t.tokens.items {
		fmt.Printf("%+v\n", item)
	}
}

func (t *Tokenizer) tokenizeLine(line string) []*Token {
	split := strings.Split(line, " ")

	result := []*Token{}
	for _, word := range split {
		token := NewToken(word, TokenKeyword)
		result = append(result, token)
	}
	return result
}

func (t *Tokenizer) tokenizeWord(word string) *Token {
	if t.isKeyword(word) {
		return NewToken(word, TokenKeyword)
	}
	return NewToken(word, TokenStringConst)
}

func (t *Tokenizer) isKeyword(word string) bool {
	keywords := []string{
		"class",
		"constructor",
		"method",
		"function",
		"var",
		"let",
		"field",
		"static",
		"void",
		"boolean",
		"int",
		"return",
		"false",
		"true",
		"if",
		"else",
		"do",
		"while",
	}

	return t.contains(word, keywords)
}

func (t *Tokenizer) contains(value string, items []string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}
