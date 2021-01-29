package token

import (
	"strconv"
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

func (t *Tokenizer) Tokenize() *Tokens {
	for _, line := range t.lines {
		items := t.tokenizeLine(line)
		t.tokens.Add(items)
	}
	return t.tokens
}

func (t *Tokenizer) tokenizeLine(line string) []*Token {
	// 半角空白で文字列を分割
	splitBySpace := t.splitBySpaces(line)

	// シンボルで文字列を分割
	words := []string{}
	for _, split := range splitBySpace {
		words = append(words, t.splitBySymbols(split)...)
	}

	// 分割した文字列からTokenインスタンスを生成
	result := []*Token{}
	for _, word := range words {
		token := t.tokenizeWord(word)
		result = append(result, token)
	}
	return result
}

func (t *Tokenizer) splitBySpaces(line string) []string {
	// 「"」が含まれない場合は文字列定値が含まれないので、何も考えず半角空白で分割する
	if !strings.Contains(line, "\"") {
		return strings.Fields(line)
	}

	// 「"」が含まれる場合は文字列定値が含まれる
	// string constant内の空白は分割してはいけないため,
	// 最初に「"」で分割して、最初と最後の要素をさらに半角空白で分割する
	// たとえば「foo = "test string";」は「foo」「=」「test string」「;」に分割される
	result := []string{}
	splitByDoubleQuote := strings.Split(line, "\"")
	result = append(result, strings.Fields(splitByDoubleQuote[0])...)
	result = append(result, stringConstMarker+splitByDoubleQuote[1])
	result = append(result, strings.Fields(splitByDoubleQuote[2])...)
	return result
}

func (t *Tokenizer) splitBySymbols(line string) []string {
	result := []string{line}
	if !t.containSymbol(line) {
		return result
	}

	for _, symbol := range symbolElements {
		//fmt.Printf("symbol: '%s' ---\n", symbol)
		//printSlice(result, "line")
		//fmt.Println("---")
		tmp := []string{}
		for _, element := range result {
			split := t.splitBySymbol(element, symbol)
			tmp = append(tmp, split...)
		}
		result = tmp
	}

	return result
}

func (t *Tokenizer) splitBySymbol(line string, symbol string) []string {
	// 指定したシンボルを含まない場合は即終了
	if !strings.Contains(line, symbol) {
		return []string{line}
	}

	// 指定したシンボルで入力文字列を分割
	split := strings.Split(line, symbol)
	result := []string{}
	//printSlice(split, "split")

	for _, s := range split {
		// split後の要素が空文字以外の場合は、シンボルではない
		// そこでトリムしたうえで、返り値にトリム後の文字列を追加する
		if s != "" {
			word := strings.TrimSpace(s)
			result = append(result, word)
		}

		// 分割した要素ごとにシンボルを追加する
		// なお最後の追加されるシンボルは余分なので、このあとで取り除く
		result = append(result, symbol)
	}

	// 返り値の最後の要素にシンボルが一個多い状態なので、それを取り除く
	result = result[:len(result)-1]
	//printSlice(result, "result")

	return result
}

func (t *Tokenizer) containSymbol(word string) bool {
	for _, symbol := range symbolElements {
		if strings.Contains(word, symbol) {
			return true
		}
	}
	return false
}

func (t *Tokenizer) tokenizeWord(word string) *Token {
	switch {
	case t.isKeyword(word):
		return NewToken(word, TokenKeyword)
	case t.isSymbol(word):
		return NewToken(word, TokenSymbol)
	case t.isIntConst(word):
		return NewToken(word, TokenIntConst)
	case t.isStringConst(word):
		deletedMarker := word[len(stringConstMarker):]
		return NewToken(deletedMarker, TokenStringConst)
	default:
		return NewToken(word, TokenIdentifier)
	}
}

func (t *Tokenizer) isKeyword(word string) bool {
	return t.contains(word, keywordElements)
}

func (t *Tokenizer) isSymbol(word string) bool {
	return t.contains(word, symbolElements)
}

func (t *Tokenizer) isIntConst(word string) bool {
	_, err := strconv.Atoi(word)
	return err == nil
}

func (t *Tokenizer) isStringConst(word string) bool {
	return strings.Contains(word, stringConstMarker)
}

// 文字列定値であることを示すマーカー
const stringConstMarker = "\"\"\""

func (t *Tokenizer) contains(value string, items []string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

var keywordElements = []string{
	"class",
	"constructor",
	"function",
	"method",
	"field",
	"static",
	"var",
	"int",
	"char",
	"boolean",
	"void",
	"true",
	"false",
	"null",
	"this",
	"let",
	"do",
	"if",
	"else",
	"while",
	"return",
}

var symbolElements = []string{
	"{",
	"}",
	"(",
	")",
	"[",
	"]",
	".",
	",",
	";",
	"+",
	"-",
	"*",
	"/",
	"&",
	"|",
	"<",
	">",
	"=",
	"~",
}
