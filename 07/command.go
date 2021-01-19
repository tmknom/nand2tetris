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

func (cs *Commands) ParseAll() error {
	for _, command := range cs.commands {
		err := command.Parse()
		if err != nil {
			return err
		}
	}
	return nil
}

func (cs *Commands) ConvertAll() []string {
	result := []string{}
	for _, command := range cs.commands {
		converter := NewConverter(command.commandType, command.arg1, command.arg2)
		assembler := converter.Convert()
		result = append(result, assembler...)
	}
	return result
}

// デバッグ用：コマンドのダンプ
func (cs *Commands) dump() {
	for i, command := range cs.commands {
		fmt.Printf("Command[%d]: %s\n", i, command.tostring())
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
	var commandType CommandType
	if commandTypeString == "push" {
		commandType = CommandPush
	} else if commandTypeString == "pop" {
		commandType = CommandPop
	} else if commandTypeString == "label" {
		commandType = CommandLabel
	} else if commandTypeString == "goto" {
		commandType = CommandGoto
	} else if commandTypeString == "if" {
		commandType = CommandIf
	} else if commandTypeString == "function" {
		commandType = CommandFunction
	} else if commandTypeString == "return" {
		commandType = CommandReturn
	} else if commandTypeString == "call" {
		commandType = CommandCall
	} else {
		return nil, fmt.Errorf("not implemented: %s\n", commandTypeString)
	}
	return &commandType, nil
}

func (c *Command) tostring() string {
	split := strings.Split(c.raw, " ")
	commandTypeString := split[0]

	arg2 := "nil"
	if c.arg2 != nil {
		arg2 = strconv.Itoa(*c.arg2)
	}
	return fmt.Sprintf("&{raw: %s, commandType: %s, arg1: %s arg2: %s}", c.raw, commandTypeString, c.arg1, arg2)
}
