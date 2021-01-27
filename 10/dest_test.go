package main

import (
	"testing"
)

func TestDestTokenizedXMLFilename(t *testing.T) {
	cases := []struct {
		desc     string
		filename string
		want     string
	}{
		{
			desc:     "出力ファイル名の生成",
			filename: "Square/Test.jack",
			want:     "Square/TestT.xml",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			dest := NewDest(tc.filename)
			got := dest.tokenizedXMLFilename()
			if got != tc.want {
				t.Errorf("failed: got = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestDestParsedXMLFilename(t *testing.T) {
	cases := []struct {
		desc     string
		filename string
		want     string
	}{
		{
			desc:     "出力ファイル名の生成",
			filename: "Square/Test.jack",
			want:     "Square/Test.xml",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			dest := NewDest(tc.filename)
			got := dest.parsedXMLFilename()
			if got != tc.want {
				t.Errorf("failed: got = %s, want %s", got, tc.want)
			}
		})
	}
}
