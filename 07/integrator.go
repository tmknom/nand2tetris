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

	translators := i.factoryTranslators(commands)
	assembler := translators.TranslatorAll()

	dest := NewDest(i.filename)
	err = dest.Write(assembler)
	if err != nil {
		return err
	}

	return nil
}

func (i *Integrator) factoryTranslators(commands *Commands) *Translators {
	translators := NewTranslators()
	for _, command := range commands.commands {
		translators.Add(command)
	}
	return translators
}
