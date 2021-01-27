package main

type Integrator struct {
	filenames []string
}

func NewIntegrator(filenames []string) *Integrator {
	return &Integrator{filenames: filenames}
}

func (i *Integrator) Integrate() error {
	for _, file := range i.filenames {
		err := i.integrateFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Integrator) integrateFile(file string) error {
	// ソースファイルの読み込み
	src := NewSrc(file)
	err := src.Setup()
	if err != nil {
		return err
	}

	// トークンに分割
	tokenizer := NewTokenizer(src.lines)
	tokens := tokenizer.Tokenize()
	xml := tokens.ToXML()

	// トークンをパース
	parser := NewParser(tokens)
	parser.Parse()

	// XMLファイルへ書き込み
	dest := NewDest(src.filename)
	err = dest.Write(xml)
	if err != nil {
		return err
	}

	return nil
}
