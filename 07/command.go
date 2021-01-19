package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Commands struct {
	commands []*Command
}

func NewCommands() *Commands {
	return &Commands{commands: []*Command{}}
}

func (cs *Commands) Add(command *Command) {
	cs.commands = append(cs.commands, command)
}

// デバッグ用：コマンドのみをダンプ
func (cs *Commands) dump() {
	for i, command := range cs.commands {
		fmt.Printf("Command[%d]: %s\n", i, command.raw)
	}
}

type Command struct {
	raw         string
	commandType CommandType
	arg1        string
	arg2        *int
}

type CommandType int

const (
	CommandArithmetic CommandType = iota
	CommandPush
	CommandPop
	CommandLabel
	CommandGoto
	CommandIf
	CommandFunction
	CommandReturn
	CommandCall
)

func NewCommand(raw string) *Command {
	return &Command{raw: raw}
}

func (c *Command) Parse() error {
	split := strings.Split(c.raw, " ")
	commandLength := len(split)
	if commandLength == 1 {
		c.parseLength1(split)
		return nil
	} else if commandLength == 2 {
		return c.parseLength2(split)
	} else if commandLength == 3 {
		return c.parseLength3(split)
	}
	return fmt.Errorf("invalid command: %s\n", c.raw)
}

func (c *Command) parseLength1(split []string) {
	c.commandType = CommandArithmetic
	c.arg1 = split[0]
}

func (c *Command) parseLength2(split []string) error {
	commandType, err := c.parseCommandType(split[0])
	if err != nil {
		return err
	}

	c.commandType = *commandType
	c.arg1 = split[1]
	return nil
}

func (c *Command) parseLength3(split []string) error {
	commandType, err := c.parseCommandType(split[0])
	if err != nil {
		return err
	}

	num, err := strconv.Atoi(split[2])
	if err != nil {
		return err
	}

	c.commandType = *commandType
	c.arg1 = split[1]
	c.arg2 = &num
	return nil
}

func (c *Command) parseCommandType(commandTypeString string) (*CommandType, error) {
	return nil, fmt.Errorf("not implemented: %s\n", commandTypeString)
}
