package main

import (
	"./io"
	"./parsing"
	"./token"
)

type Integrator struct {
	filenames []string
}

func NewIntegrator(filenames []string) *Integrator {
	return &Integrator{filenames: filenames}
}

func (i *Integrator) Integrate() error {
	for _, filename := range i.filenames {
		err := i.integrateFile(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Integrator) integrateFile(file string) error {
	// ソースファイルの読み込み
	src := io.NewSrc(file)
	err := src.Setup()
	if err != nil {
		return err
	}

	// トークンに分割
	tokenizer := token.NewTokenizer(src.Lines)
	tokens := tokenizer.Tokenize()
	tokenizedXML := tokens.ToXML()

	// トークンをパース
	parser := parsing.NewParser(tokens, src.ClassName())
	class, err := parser.Parse()
	if err != nil {
		return err
	}

	// XMLファイルへ書き込み
	dest := io.NewDest(src.Filename)
	err = dest.WriteTokenizedXML(tokenizedXML)
	if err != nil {
		return err
	}

	err = dest.WriteParsedXML(class.ToXML())
	if err != nil {
		return err
	}

	// デバッグしやすいように生成したコードを標準出力
	parser.PrintDebugCode()

	// コード生成をして書き込み
	err = dest.WriteCode(parser.CodeLines())
	if err != nil {
		return err
	}

	return nil
}
