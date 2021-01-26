package main

import "fmt"

type Integrator struct {
	filenames []string
}

func NewIntegrator(filenames []string) *Integrator {
	return &Integrator{filenames: filenames}
}

func (i *Integrator) Integrate() error {
	for _, file := range i.filenames {
		fmt.Println(file)
		err := i.integrateFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Integrator) integrateFile(file string) error {
	src := NewSrc(file)
	err := src.Setup()
	if err != nil {
		return err
	}

	tokenizer := NewTokenizer(src.lines)
	tokenizer.Tokenize()

	return nil
}
