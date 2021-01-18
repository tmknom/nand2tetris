package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	raw         []*string
	symbolTable *SymbolTable
	nextAddress int
}

func NewParser(raw []*string, symbolTable *SymbolTable) *Parser {
	return &Parser{raw: raw, symbolTable: symbolTable, nextAddress: 0}
}

func (p *Parser) Parse() []Command {
	var commands []Command
	for _, line := range p.raw {
		// 有効なコマンドのみ取得
		command := p.parseLine(*line)
		if command == nil {
			continue
		}

		commands = append(commands, command)
	}
	return commands
}

func (p *Parser) incrementAddress() {
	p.nextAddress += 1
}

func (p *Parser) parseLine(line string) Command {
	// コメントを除外
	deletedComment := line
	if strings.Contains(line, "//") {
		deletedComment = line[:strings.Index(line, "//")]
	}

	// 空白を除去
	trimmed := strings.TrimSpace(deletedComment)
	if len(trimmed) == 0 {
		return nil
	}

	prefix := trimmed[0]

	// ラベルシンボルを見つけたら、シンボルテーブルに追加
	if prefix == '(' {
		symbol := trimmed[1 : len(trimmed)-1]
		p.symbolTable.AddEntry(symbol, p.nextAddress)
		fmt.Printf("parse LCommand: %s, symbol: %s, address: %d\n", trimmed, symbol, p.symbolTable.Address(symbol))
		return nil
	}

	// AコマンドとCコマンドをパース
	address := p.nextAddress
	if prefix == '@' {
		fmt.Printf("parse ACommand[%d]: %s\n", address, trimmed)
		p.incrementAddress()
		return &ACommand{mnemonic: trimmed, address: address}
	} else {
		fmt.Printf("parse CCommand[%d]: %s\n", address, trimmed)
		p.incrementAddress()
		return &CCommand{mnemonic: trimmed, dest: "", comp: "", jump: "", address: address}
	}
}

type CCommand struct {
	mnemonic string
	dest     string
	comp     string
	jump     string
	address  int
}

func (c *CCommand) assemble() (string, error) {
	c.parseMnemonic()

	dest := c.assembleDest()
	jump := c.assembleJump()
	comp := c.assembleComp()
	result := fmt.Sprintf("111%s%s%s", comp, dest, jump)
	// fmt.Printf("C Command: before: %s comp: %s, dest: %s, jump: %s, after: %s\n", c.mnemonic, c.comp, c.dest, c.jump, result)
	return result, nil
}

func (c *CCommand) assembleComp() string {
	compMap := map[string]string{
		"0":   "0101010",
		"1":   "0111111",
		"-1":  "0111010",
		"D":   "0001100",
		"A":   "0110000",
		"M":   "1110000",
		"!D":  "0001101",
		"!A":  "0110001",
		"!M":  "1110001",
		"-D":  "0001111",
		"-A":  "0110011",
		"-M":  "1110011",
		"D+1": "0011111",
		"A+1": "0110111",
		"M+1": "1110111",
		"D-1": "0001110",
		"A-1": "0110010",
		"M-1": "1110010",
		"D+A": "0000010",
		"D+M": "1000010",
		"D-A": "0010011",
		"D-M": "1010011",
		"A-D": "0000111",
		"M-D": "1000111",
		"D&A": "0000000",
		"D&M": "1000000",
		"D|A": "0010101",
		"D|M": "1010101",
	}
	return compMap[c.comp]
}

func (c *CCommand) assembleJump() string {
	jumpMap := map[string]string{
		"":    "000",
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
		"":    "000",
		"M":   "001",
		"D":   "010",
		"MD":  "011",
		"A":   "100",
		"AM":  "101",
		"AD":  "110",
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
	if strings.Contains(c.mnemonic, "=") {
		split := strings.Split(c.mnemonic, "=")
		c.dest = split[0]
	}
}

func (c *CCommand) parseJump() {
	if strings.Contains(c.mnemonic, ";") {
		split := strings.Split(c.mnemonic, ";")
		c.jump = split[1]
	}
}

func (c *CCommand) parseComp() {
	compAndJump := c.mnemonic
	if strings.Contains(c.mnemonic, "=") {
		split := strings.Split(c.mnemonic, "=")
		compAndJump = split[1]
	}

	if strings.Contains(compAndJump, ";") {
		split := strings.Split(compAndJump, ";")
		c.comp = split[0]
	} else {
		c.comp = compAndJump
	}
}

type ACommand struct {
	mnemonic string
	address  int
}

func (a *ACommand) assemble() (string, error) {
	withoutPrefix := a.mnemonic[1:]
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
