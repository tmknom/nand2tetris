package main

type Integrator struct {
	filename string
}

func NewIntegrator(filename string) *Integrator {
	return &Integrator{filename: filename}
}

func (i *Integrator) Integrate() error {
	vmCode, hasInit, err := ReadVmCode(i.filename)
	if err != nil {
		return err
	}

	commands := vmCode.Commands
	err = commands.ParseAll()
	if err != nil {
		return err
	}

	translators := i.factoryTranslators(commands, i.filename, hasInit)
	assembler := translators.TranslateAll()

	dest := NewDest(i.filename)
	err = dest.Write(assembler)
	if err != nil {
		return err
	}

	return nil
}

func (i *Integrator) factoryTranslators(commands *Commands, filename string, hasInit HasInit) *Translators {
	translators := NewTranslators(filename, hasInit)
	for _, command := range commands.commands {
		translators.Add(command)
	}
	return translators
}
