package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	raw []*string
	commands []Command
}

func NewParser(raw []*string) *Parser {
	return &Parser{raw: raw}
}

func (p *Parser) Parse() error{
	for _, line := range p.raw {
		// 有効なコマンドのみ取得
		command := p.parseLine(line)
		if command == nil {
			continue
		}

		// アセンブル
		p.commands = append(p.commands, command)
		assembled, err := command.assemble()
		if err != nil {
			return err
		}
		fmt.Printf("assembled: before: %s, after: %s\n", *line, assembled)
	}
	return nil
}

func (p *Parser) parseLine(line *string) Command {
	// 空白を除去
	trimmed := strings.TrimSpace(*line)
	if len(trimmed) == 0 {
		return nil
	}

	// コメントを除外
	if len(trimmed) >= 2 && trimmed[:2] == "//" {
		return nil
	}

	prefix := trimmed[0]
	if prefix == '@' {
		return &ACommand{raw: trimmed}
	} else {
		return &CCommand{raw: trimmed, dest: "", comp: "", jump: ""}
	}
}

type CCommand struct {
	raw string
	dest string
	comp string
	jump string
}

func (c *CCommand) assemble() (string, error) {
	c.parseMnemonic()

	dest := c.assembleDest()
	jump := c.assembleJump()
	result := fmt.Sprintf("111 %s %s", dest, jump)
	fmt.Printf("C Command: before: %s comp: %s, dest: %s, jump: %s, after: %s\n", c.raw, c.comp, c.dest, c.jump, result)
	return result, nil
}

func (c *CCommand) assembleJump() string {
	jumpMap := map[string]string{
		"": "000",
		"JGT": "001",
		"JEQ": "010",
		"JGE": "011",
		"JLT": "100",
		"JNE": "101",
		"JLE": "110",
		"JMP": "111",
	}
	return jumpMap[c.jump]
}

func (c *CCommand) assembleDest() string {
	destMap := map[string]string{
		"": "000",
		"M": "001",
		"D": "010",
		"MD": "011",
		"A": "100",
		"AM": "101",
		"AD": "110",
		"AMD": "111",
	}
	return destMap[c.dest]
}

func (c *CCommand) parseMnemonic() {
	c.parseDest()
	c.parseJump()
	c.parseComp()
}

func (c *CCommand) parseDest() {
	if strings.Contains(c.raw, "=") {
		split := strings.Split(c.raw, "=")
		c.dest = split[0]
	}
}

func (c *CCommand) parseJump() {
	if strings.Contains(c.raw, ";") {
		split := strings.Split(c.raw, ";")
		c.jump = split[1]
	}
}

func (c *CCommand) parseComp() {
	compAndJump := c.raw
	if strings.Contains(c.raw, "=") {
		split := strings.Split(c.raw, "=")
		compAndJump = split[1]
	}

	if strings.Contains(compAndJump, ";") {
		split := strings.Split(compAndJump, ";")
		c.comp = split[0]
	}else{
		c.comp = compAndJump
	}
}

type ACommand struct {
	raw string
}

func (a *ACommand) assemble() (string, error) {
	withoutPrefix := a.raw[1:]
	num, err := strconv.Atoi(withoutPrefix)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("0%015b", num)
	// fmt.Printf("A Command: before: %s, after: %s\n", withoutPrefix, result)
	return result, nil
}

type Command interface {
	assemble() (string, error)
}
