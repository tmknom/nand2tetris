package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 読み込んだvmファイルのコードを保持
type VmCode struct {
	org      []string
	commands []string
}

func ReadVmCode(filename string) (*VmCode, error) {
	reader := &vmCodeReader{}
	lines, commands, err := reader.read(filename)
	if err != nil {
		return nil, err
	}
	return newVmCode(lines, commands), nil
}

func newVmCode(org []string, commands []string) *VmCode {
	return &VmCode{org: org, commands: commands}
}

// デバッグ用：コマンドのみをダンプ
func (vc *VmCode) dumpCommands() {
	for i, line := range vc.commands {
		fmt.Printf("Command[%d]: %s\n", i, line)
	}
}

// デバッグ用：オリジナルのVMファイルのコードをダンプ
func (vc *VmCode) dumpOrg() {
	for i, line := range vc.org {
		fmt.Printf("VmCode[%d]: %s\n", i, line)
	}
}

type vmCodeReader struct{}

func (r *vmCodeReader) read(filename string) ([]string, []string, error) {
	lines, err := r.readLines(filename)
	if err != nil {
		return nil, nil, err
	}
	commands := r.createCommands(lines)

	return lines, commands, nil
}

func (r *vmCodeReader) createCommands(lines []string) []string {
	commands := []string{}
	for _, line := range lines {
		withoutComment := r.deleteCommentAndWhitespace(line)
		if withoutComment != "" {
			commands = append(commands, withoutComment)
		}
	}
	return commands
}

func (r *vmCodeReader) deleteCommentAndWhitespace(line string) string {
	// コメントを除外
	deletedComment := line
	if strings.Contains(line, "//") {
		deletedComment = line[:strings.Index(line, "//")]
	}

	// 空白を除去
	return strings.TrimSpace(deletedComment)
}

func (r *vmCodeReader) readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
