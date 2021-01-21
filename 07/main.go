package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}

func run() error {
	arg := NewArg(os.Args)
	for _, file := range arg.files {
		fmt.Printf("vmファイルの変換開始：%s\n", file)
		integrator := NewIntegrator(file)
		err := integrator.Integrate()
		if err != nil {
			return err
		}
	}

	return nil
}
