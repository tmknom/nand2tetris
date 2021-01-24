package main

type Integrator struct {
	filenames []string
	commands  *Commands
}

func NewIntegrator(filenames []string) *Integrator {
	return &Integrator{filenames: filenames, commands: NewCommands()}
}

func (i *Integrator) Integrate() error {
	for _, file := range i.filenames {
		// ファイルを読み込んでメモリに展開
		src := NewSrc(file)
		err := src.Setup()
		if err != nil {
			return err
		}

		// Commandの生成
		i.generateCommands(src)

		// TODO あとで消す
		err = i.integrateFile(file)
		if err != nil {
			return err
		}
	}

	err := i.commands.ParseAll()
	if err != nil {
		return err
	}

	return nil
}

func (i *Integrator) integrateFile(file string) error {
	vmCode, err := ReadVmCode(file)
	if err != nil {
		return err
	}

	commands := vmCode.Commands
	err = commands.ParseAll()
	if err != nil {
		return err
	}

	translators := i.factoryTranslators(commands)
	assembler := translators.TranslateAll()

	dest := NewDest(file)
	err = dest.Write(assembler)
	if err != nil {
		return err
	}

	return nil
}

func (i *Integrator) generateCommands(src *Src) {
	for _, line := range src.lines {
		command := NewCommand(line, &src.moduleName)
		i.commands.Add(command)
	}
}

func (i *Integrator) factoryTranslators(commands *Commands) *Translators {
	translators := NewTranslators()
	for _, command := range commands.commands {
		translators.Add(command)
	}
	return translators
}
