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
	vmCode, err := ReadVmCode(file)
	if err != nil {
		return err
	}

	commands := vmCode.Commands
	err = commands.ParseAll()
	if err != nil {
		return err
	}

	translators := i.factoryTranslators(commands, file)
	assembler := translators.TranslateAll()

	dest := NewDest(file)
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
