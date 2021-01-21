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

	converters := i.factoryConverters(commands)
	assembler := converters.ConvertAll()

	dest := NewDest(i.filename)
	err = dest.Write(assembler)
	if err != nil {
		return err
	}

	return nil
}

func (i *Integrator) factoryConverters(commands *Commands) *Converters {
	converters := NewConverters()
	for _, command := range commands.commands {
		converters.Add(command)
	}
	return converters
}
