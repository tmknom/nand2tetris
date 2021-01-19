package main

import (
	"fmt"
	"log"
	"os"
)

func run() error {
	src := NewSrc(os.Args)
	src.Parse()
	for _, file := range src.files {
		fmt.Printf("vmファイルの変換開始：%s\n", file)
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}
