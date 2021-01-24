package main

type Integrator struct {
	filenames []string
	arg       string
	commands  *Commands
}

func NewIntegrator(filenames []string, arg string) *Integrator {
	return &Integrator{filenames: filenames, arg: arg, commands: NewCommands()}
}

func (i *Integrator) Integrate() error {
	// ファイルの読み込み
	err := i.readFiles()
	if err != nil {
		return err
	}

	// コマンドのパース
	err = i.commands.Parse()
	if err != nil {
		return err
	}

	// アセンブル
	assembler := i.translate()

	// アセンブラの書き込み
	dest := NewDest(i.arg)
	err = dest.Write(assembler)
	if err != nil {
		return err
	}

	return nil
}

func (i *Integrator) readFiles() error {
	for _, file := range i.filenames {
		// ファイルを読み込んでメモリに展開
		src := NewSrc(file)
		err := src.Setup()
		if err != nil {
			return err
		}

		// Commandの生成
		i.generateCommands(src)
	}
	return nil
}

func (i *Integrator) generateCommands(src *Src) {
	for _, line := range src.lines {
		command := NewCommand(line, &src.moduleName)
		i.commands.Add(command)
	}
}

func (i *Integrator) translate() []string {
	translators := NewTranslators()
	for _, command := range i.commands.commands {
		translators.Add(command)
	}
	return translators.TranslateAll()
}
