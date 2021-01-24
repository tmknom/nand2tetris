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
	fmt.Printf("vmファイルの変換開始：%s\n", arg.raw)
	integrator := NewIntegrator(arg.files)
	return integrator.Integrate()
}
