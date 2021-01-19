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
	src := NewSrc(os.Args)
	src.Parse()
	for _, file := range src.files {
		err := convert(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func convert(file string) error {
	fmt.Printf("vmファイルの変換開始：%s\n", file)
	vmCode, err := ReadVmCode(file)
	if err != nil {
		return err
	}

	commands := vmCode.Commands
	commands.dump()
	return nil
}
