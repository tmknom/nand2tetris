package io

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type Dest struct {
	src string
}

func NewDest(src string) *Dest {
	return &Dest{src: src}
}

func (d *Dest) WriteTokenizedXML(lines []string) error {
	filename := d.tokenizedXMLFilename()
	return d.write(filename, lines)
}

func (d *Dest) WriteParsedXML(lines []string) error {
	filename := d.parsedXMLFilename()
	return d.write(filename, lines)
}

func (d *Dest) write(filename string, lines []string) error {
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

func (d *Dest) tokenizedXMLFilename() string {
	withoutExt := d.src[:len(d.src)-len(filepath.Ext(d.src))]
	return fmt.Sprintf("%sT.xml", withoutExt)
}

func (d *Dest) parsedXMLFilename() string {
	withoutExt := d.src[:len(d.src)-len(filepath.Ext(d.src))]
	return fmt.Sprintf("%s.xml", withoutExt)
}
