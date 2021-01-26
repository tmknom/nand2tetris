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
	}

	return nil
}
