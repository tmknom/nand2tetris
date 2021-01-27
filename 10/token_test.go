package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestTokensToXML(t *testing.T) {
	cases := []struct {
		desc  string
		token []*Token
		want  []string
	}{
		{
			desc: "複数トークンのXML変換",
			token: []*Token{
				NewToken("class", TokenKeyword),
				NewToken("Array", TokenIdentifier),
				NewToken("{", TokenSymbol),
			},
			want: []string{
				"<tokens>",
				"<keyword> class </keyword>",
				"<identifier> Array </identifier>",
				"<symbol> { </symbol>",
				"</tokens>",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := NewTokens()
			tokens.Add(tc.token)
			got := tokens.ToXML()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestTokenToXML(t *testing.T) {
	cases := []struct {
		desc  string
		token *Token
		want  string
	}{
		{
			desc:  "キーワード",
			token: NewToken("class", TokenKeyword),
			want:  "<keyword> class </keyword>",
		},
		{
			desc:  "シンボル",
			token: NewToken("+", TokenSymbol),
			want:  "<symbol> + </symbol>",
		},
		{
			desc:  "数字定値",
			token: NewToken("123", TokenIntConst),
			want:  "<integerConstant> 123 </integerConstant>",
		},
		{
			desc:  "文字列定値",
			token: NewToken("foo bar", TokenStringConst),
			want:  "<stringConstant> foo bar </stringConstant>",
		},
		{
			desc:  "識別子",
			token: NewToken("Array", TokenIdentifier),
			want:  "<identifier> Array </identifier>",
		},
		{
			desc:  "シンボルのエンコード「<」",
			token: NewToken("<", TokenSymbol),
			want:  "<symbol> &lt; </symbol>",
		},
		{
			desc:  "シンボルのエンコード「>」",
			token: NewToken(">", TokenSymbol),
			want:  "<symbol> &gt; </symbol>",
		},
		{
			desc:  "シンボルのエンコード「&」",
			token: NewToken("&", TokenSymbol),
			want:  "<symbol> &amp; </symbol>",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.token.ToXML()
			if got != tc.want {
				t.Errorf("failed: got = %s, want = %s", got, tc.want)
			}
		})
	}
}
