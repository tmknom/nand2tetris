package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type Dest struct {
	filename string
}

func NewDest(filename string) *Dest {
	withoutExt := filename[:len(filename)-len(filepath.Ext(filename))]
	return &Dest{filename: fmt.Sprintf("%s.asm", withoutExt)}
}

func (d *Dest) Write(lines []string) error {
	file, err := os.Create(d.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.Write([]byte(line + "\n"))
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}
