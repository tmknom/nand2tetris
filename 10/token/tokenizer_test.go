package token

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestTokenizerTokenize(t *testing.T) {
	cases := []struct {
		desc  string
		lines []string
		want  []*Token
	}{
		{
			desc: "複数行",
			lines: []string{
				"class Main{",
				"let foo = 34;",
				"let bar = \"test string\";",
			},
			want: []*Token{
				NewToken("class", TokenKeyword),
				NewToken("Main", TokenIdentifier),
				NewToken("{", TokenSymbol),
				NewToken("let", TokenKeyword),
				NewToken("foo", TokenIdentifier),
				NewToken("=", TokenSymbol),
				NewToken("34", TokenIntConst),
				NewToken(";", TokenSymbol),
				NewToken("let", TokenKeyword),
				NewToken("bar", TokenIdentifier),
				NewToken("=", TokenSymbol),
				NewToken("test string", TokenStringConst),
				NewToken(";", TokenSymbol),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokenizer := NewTokenizer(tc.lines)
			tokenizer.Tokenize()

			if len(tokenizer.tokens.Items) > len(tc.want) {
				want := &Tokens{Items: tc.want}
				t.Fatalf("failed: size: got = %s,\nwant:%s\n", tokenizer.tokens.Debug(), want.Debug())
			}

			for i, token := range tokenizer.tokens.Items {
				if diff := cmp.Diff(token, tc.want[i]); diff != "" {
					t.Errorf("failed: token[%d]: diff (-got +want):\n%s", i, diff)
				}
			}
		})
	}
}

func TestTokenizerTokenizeLine(t *testing.T) {
	cases := []struct {
		desc string
		line string
		want []*Token
	}{
		{
			desc: "キーワードをひとつだけ含む",
			line: "class",
			want: []*Token{
				NewToken("class", TokenKeyword),
			},
		},
		{
			desc: "キーワードとシンボルをひとつずつ含む",
			line: "return;",
			want: []*Token{
				NewToken("return", TokenKeyword),
				NewToken(";", TokenSymbol),
			},
		},
		{
			desc: "数字定値をひとつだけ含む",
			line: "13",
			want: []*Token{
				NewToken("13", TokenIntConst),
			},
		},
		{
			desc: "文字列定値をひとつだけ含む",
			line: "\"test string\"",
			want: []*Token{
				NewToken("test string", TokenStringConst),
			},
		},
		{
			desc: "識別子をひとつだけ含む",
			line: "Array",
			want: []*Token{
				NewToken("Array", TokenIdentifier),
			},
		},
		{
			desc: "識別子とシンボルとキーワードをひとつずつ含む",
			line: "class Main{",
			want: []*Token{
				NewToken("class", TokenKeyword),
				NewToken("Main", TokenIdentifier),
				NewToken("{", TokenSymbol),
			},
		},
		{
			desc: "識別子とシンボルと数字定値をひとつずつ含む",
			line: "foo = 34",
			want: []*Token{
				NewToken("foo", TokenIdentifier),
				NewToken("=", TokenSymbol),
				NewToken("34", TokenIntConst),
			},
		},
		{
			desc: "識別子とシンボルと文字列定値をひとつずつ含む",
			line: "foo = \"test string\"",
			want: []*Token{
				NewToken("foo", TokenIdentifier),
				NewToken("=", TokenSymbol),
				NewToken("test string", TokenStringConst),
			},
		},
		{
			desc: "複数の要素を含む（その1）",
			line: "while (~(key = 0)) {",
			want: []*Token{
				NewToken("while", TokenKeyword),
				NewToken("(", TokenSymbol),
				NewToken("~", TokenSymbol),
				NewToken("(", TokenSymbol),
				NewToken("key", TokenIdentifier),
				NewToken("=", TokenSymbol),
				NewToken("0", TokenIntConst),
				NewToken(")", TokenSymbol),
				NewToken(")", TokenSymbol),
				NewToken("{", TokenSymbol),
			},
		},
		{
			desc: "複数の要素を含む（その2）",
			line: "do Screen.drawRectangle((x + size) - 1, y, x + size, y + size);",
			want: []*Token{
				NewToken("do", TokenKeyword),
				NewToken("Screen", TokenIdentifier),
				NewToken(".", TokenSymbol),
				NewToken("drawRectangle", TokenIdentifier),
				NewToken("(", TokenSymbol),
				NewToken("(", TokenSymbol),
				NewToken("x", TokenIdentifier),
				NewToken("+", TokenSymbol),
				NewToken("size", TokenIdentifier),
				NewToken(")", TokenSymbol),
				NewToken("-", TokenSymbol),
				NewToken("1", TokenIntConst),
				NewToken(",", TokenSymbol),
				NewToken("y", TokenIdentifier),
				NewToken(",", TokenSymbol),
				NewToken("x", TokenIdentifier),
				NewToken("+", TokenSymbol),
				NewToken("size", TokenIdentifier),
				NewToken(",", TokenSymbol),
				NewToken("y", TokenIdentifier),
				NewToken("+", TokenSymbol),
				NewToken("size", TokenIdentifier),
				NewToken(")", TokenSymbol),
				NewToken(";", TokenSymbol),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokenizer := NewTokenizer([]string{})
			tokens := tokenizer.tokenizeLine(tc.line)

			if diff := cmp.Diff(tokens, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestTokenizerSplitBySpaces(t *testing.T) {
	cases := []struct {
		desc string
		line string
		want []string
	}{
		{
			desc: "文字列定値を含まない",
			line: "let s = null;",
			want: []string{
				"let",
				"s",
				"=",
				"null;",
			},
		},
		{
			desc: "文字列定値を含む",
			line: "foo = \"test string\";",
			want: []string{
				"foo",
				"=",
				stringConstMarker + "test string",
				";",
			},
		},
		{
			desc: "文字列定値と複数の空白を含む",
			line: "  foo   =    \"  foo bar  baz   \"  ;   ",
			want: []string{
				"foo",
				"=",
				stringConstMarker + "  foo bar  baz   ",
				";",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokenizer := NewTokenizer([]string{})
			got := tokenizer.splitBySpaces(tc.line)

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestTokenizerSplitBySymbols(t *testing.T) {
	cases := []struct {
		desc string
		line string
		want []string
	}{
		{
			desc: "シンボルを含まない",
			line: "class",
			want: []string{
				"class",
			},
		},
		{
			desc: "最後にのみシンボルを含む",
			line: "return;",
			want: []string{
				"return",
				";",
			},
		},
		{
			desc: "最初にのみシンボルを含む",
			line: "-20",
			want: []string{
				"-",
				"20",
			},
		},
		{
			desc: "最初と最後にシンボルを含む",
			line: "(false)",
			want: []string{
				"(",
				"false",
				")",
			},
		},
		{
			desc: "複数のシンボルを含む（その1）",
			line: "while (~(key = 0)) {",
			want: []string{
				"while",
				"(",
				"~",
				"(",
				"key",
				"=",
				"0",
				")",
				")",
				"{",
			},
		},
		{
			desc: "複数のシンボルを含む（その2）",
			line: "do Screen.drawRectangle((x + size) - 1, y, x + size, y + size);",
			want: []string{
				"do Screen",
				".",
				"drawRectangle",
				"(",
				"(",
				"x",
				"+",
				"size",
				")",
				"-",
				"1",
				",",
				"y",
				",",
				"x",
				"+",
				"size",
				",",
				"y",
				"+",
				"size",
				")",
				";",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokenizer := NewTokenizer([]string{})
			got := tokenizer.splitBySymbols(tc.line)

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestTokenizerSplitBySymbol(t *testing.T) {
	cases := []struct {
		desc   string
		line   string
		symbol string
		want   []string
	}{
		{
			desc:   "指定したシンボルを含まない",
			line:   "foo",
			symbol: ";",
			want: []string{
				"foo",
			},
		},
		{
			desc:   "最後にのみ指定したシンボルを含む",
			line:   "last;",
			symbol: ";",
			want: []string{
				"last",
				";",
			},
		},
		{
			desc:   "最初にのみ指定したシンボルを含む",
			line:   "; first",
			symbol: ";",
			want: []string{
				";",
				"first",
			},
		},
		{
			desc:   "最初・途中・最後に指定したシンボルを含む",
			line:   "; x; y;",
			symbol: ";",
			want: []string{
				";",
				"x",
				";",
				"y",
				";",
			},
		},
		{
			desc:   "複数のシンボルを含む",
			line:   "; x; y; z; foo;",
			symbol: ";",
			want: []string{
				";",
				"x",
				";",
				"y",
				";",
				"z",
				";",
				"foo",
				";",
			},
		},
		{
			desc:   "最後に指定したシンボルが連続して含む",
			line:   "(~(key = 0))",
			symbol: ")",
			want: []string{
				"(~(key = 0",
				")",
				")",
			},
		},
		{
			desc:   "最初に指定したシンボルが連続して含む",
			line:   "((x + size) < 510)",
			symbol: "(",
			want: []string{
				"(",
				"(",
				"x + size) < 510)",
			},
		},
		{
			desc:   "途中に指定したシンボルが連続して含む",
			line:   "while (~(key = 0)) {",
			symbol: ")",
			want: []string{
				"while (~(key = 0",
				")",
				")",
				"{",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokenizer := NewTokenizer([]string{})
			got := tokenizer.splitBySymbol(tc.line, tc.symbol)

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestTokenizerTokenizeWord(t *testing.T) {
	cases := []struct {
		desc string
		word string
		want *Token
	}{
		{
			desc: "キーワード",
			word: "class",
			want: NewToken("class", TokenKeyword),
		},
		{
			desc: "シンボル",
			word: "+",
			want: NewToken("+", TokenSymbol),
		},
		{
			desc: "数字定値",
			word: "23",
			want: NewToken("23", TokenIntConst),
		},
		{
			desc: "識別子",
			word: "foo",
			want: NewToken("foo", TokenIdentifier),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokenizer := NewTokenizer([]string{})
			got := tokenizer.tokenizeWord(tc.word)

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}
