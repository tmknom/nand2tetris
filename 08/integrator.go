package main

type Integrator struct {
	filename string
}

func NewIntegrator(filename string) *Integrator {
	return &Integrator{filename: filename}
}

func (i *Integrator) Integrate() error {
	vmCode, err := ReadVmCode(i.filename)
	if err != nil {
		return err
	}

	commands := vmCode.Commands
	err = commands.ParseAll()
	if err != nil {
		return err
	}

	translators := i.factoryTranslators(commands, i.filename)
	assembler := translators.TranslatorAll()

	dest := NewDest(i.filename)
	err = dest.Write(assembler)
	if err != nil {
		return err
	}

	return nil
}

func (i *Integrator) factoryTranslators(commands *Commands, filename string) *Translators {
	translators := NewTranslators(filename)
	for _, command := range commands.commands {
		translators.Add(command)
	}
	return translators
}
