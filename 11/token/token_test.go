package token

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestTokensAdd(t *testing.T) {
	cases := []struct {
		desc       string
		tokens     []*Token
		wantTokens []*Token
		wantHead   int
	}{
		{
			desc: "トークンの追加",
			tokens: []*Token{
				NewToken("class", TokenKeyword),
				NewToken("Array", TokenIdentifier),
				NewToken("{", TokenSymbol),
			},
			wantTokens: []*Token{
				NewToken("class", TokenKeyword),
				NewToken("Array", TokenIdentifier),
				NewToken("{", TokenSymbol),
			},
			wantHead: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := NewTokens()
			tokens.Add(tc.tokens)

			if diff := cmp.Diff(tokens.Items, tc.wantTokens); diff != "" {
				t.Errorf("failed tokens: diff (-got +want):\n%s", diff)
			}

			if tokens.HeadIndex != tc.wantHead {
				t.Errorf("failed HeadIndex: got = %d, want %d", tokens.HeadIndex, tc.wantHead)
			}
		})
	}
}

func TestTokensAdvance(t *testing.T) {
	cases := []struct {
		desc       string
		tokens     []*Token
		wantFirst  *Token
		wantSecond *Token
	}{
		{
			desc: "前からトークンを取得",
			tokens: []*Token{
				NewToken("class", TokenKeyword),
				NewToken("Array", TokenIdentifier),
				NewToken("{", TokenSymbol),
			},
			wantFirst:  NewToken("class", TokenKeyword),
			wantSecond: NewToken("Array", TokenIdentifier),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := NewTokens()
			tokens.Add(tc.tokens)

			first := tokens.Advance()
			if diff := cmp.Diff(first, tc.wantFirst); diff != "" {
				t.Errorf("failed first: diff (-got +want):\n%s", diff)
			}

			second := tokens.Advance()
			if diff := cmp.Diff(second, tc.wantSecond); diff != "" {
				t.Errorf("failed second: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestTokensFirst(t *testing.T) {
	cases := []struct {
		desc             string
		tokens           []*Token
		wantFirst        *Token
		wantRetry        *Token
		wantAfterAdvance *Token
	}{
		{
			desc: "最初のトークンを参照",
			tokens: []*Token{
				NewToken("class", TokenKeyword),
				NewToken("Array", TokenIdentifier),
				NewToken("{", TokenSymbol),
			},
			wantFirst:        NewToken("class", TokenKeyword),
			wantRetry:        NewToken("class", TokenKeyword),
			wantAfterAdvance: NewToken("Array", TokenIdentifier),
		},
		{
			desc: "要素がもう存在しない",
			tokens: []*Token{
				NewToken("class", TokenKeyword),
			},
			wantFirst:        NewToken("class", TokenKeyword),
			wantRetry:        NewToken("class", TokenKeyword),
			wantAfterAdvance: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens := NewTokens()
			tokens.Add(tc.tokens)

			first := tokens.First()
			if diff := cmp.Diff(first, tc.wantFirst); diff != "" {
				t.Errorf("failed first: diff (-got +want):\n%s", diff)
			}

			retry := tokens.First()
			if diff := cmp.Diff(retry, tc.wantRetry); diff != "" {
				t.Errorf("failed retry: diff (-got +want):\n%s", diff)
			}

			tokens.Advance()
			afterAdvance := tokens.First()
			if diff := cmp.Diff(afterAdvance, tc.wantAfterAdvance); diff != "" {
				t.Errorf("failed after advance: diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestTokensToXML(t *testing.T) {
	cases := []struct {
		desc   string
		tokens []*Token
		want   []string
	}{
		{
			desc: "複数トークンのXML変換",
			tokens: []*Token{
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
			tokens.Add(tc.tokens)
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
