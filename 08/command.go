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
	moduleName  *string
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

func NewCommand(raw string, moduleName *string) *Command {
	return &Command{raw: raw, moduleName: moduleName}
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
	commandTypeString := split[0]
	switch commandTypeString {
	case "return":
		c.commandType = CommandReturn
	default:
		c.commandType = CommandArithmetic
	}
	c.arg1 = commandTypeString
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
	} else if commandTypeString == "if-goto" {
		commandType = CommandIf
	} else if commandTypeString == "function" {
		commandType = CommandFunction
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
