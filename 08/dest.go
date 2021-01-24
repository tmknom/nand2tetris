package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Dest struct {
	src string
}

func NewDest(src string) *Dest {
	return &Dest{src: src}
}

func (d *Dest) Write(lines []string) error {
	filename := d.generateFilename()

	file, err := os.Create(filename)
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

func (d *Dest) generateFilename() string {
	if filepath.Ext(d.src) == ".vm" {
		withoutExt := d.src[:len(d.src)-len(filepath.Ext(d.src))]
		return fmt.Sprintf("%s.asm", withoutExt)
	}

	// vmファイルを指定していない場合は、ディレクトリが指定されたとみなす
	path := strings.TrimSuffix(d.src, "/")
	split := strings.Split(path, "/")
	return fmt.Sprintf("%s/%s.asm", path, split[len(split)-1])
}
