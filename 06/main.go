package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

func run() error {
	asm := newAsm()
	lines, err := asm.read()
	if err != nil {
		return err
	}

	parser := NewParser(lines)
	commands := parser.Parse()

	var assembledLines []*string
	for _, command := range commands {
		assembled, err := command.assemble()
		if err != nil {
			return err
		}
		assembledLines = append(assembledLines, &assembled)
	}

	hack := newHack(asm.filenameWithoutExt())
	err = hack.write(assembledLines)
	if err != nil {
		return err
	}
	return nil
}

type Hack struct {
	filename string
}

func newHack(filenameWithoutExt string) *Hack {
	return &Hack{filename: filenameWithoutExt + ".hack"}
}

func (h *Hack) write(lines []*string) error {
	file, err := os.Create(h.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.Write([]byte(*line+"\n"))
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

type Asm struct {
	filename string
}

func (a *Asm) filenameWithoutExt() string{
	return filepath.Base(a.filename[:len(a.filename)-len(filepath.Ext(a.filename))])
}

func (a *Asm) read() ([]*string, error) {
	file, err := os.Open(a.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var result []*string
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, &line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func newAsm() *Asm {
	filename := "add/Add.asm"
	if len(os.Args) >= 2 {
		filename = os.Args[1]
	}
	return &Asm{filename: filename}
}

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}
