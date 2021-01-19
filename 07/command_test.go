package main

import (
	"testing"
)

func TestCommandParse1(t *testing.T) {
	cases := []struct {
		desc        string
		command     *Command
		commandType CommandType
		arg1        string
	}{
		{
			desc:        "addコマンド",
			command:     &Command{raw: "add"},
			commandType: CommandArithmetic,
			arg1:        "add",
		},
		{
			desc:        "subコマンド",
			command:     &Command{raw: "sub"},
			commandType: CommandArithmetic,
			arg1:        "sub",
		},
		{
			desc:        "eqコマンド",
			command:     &Command{raw: "eq"},
			commandType: CommandArithmetic,
			arg1:        "eq",
		},
		{
			desc:        "ltコマンド",
			command:     &Command{raw: "lt"},
			commandType: CommandArithmetic,
			arg1:        "lt",
		},
		{
			desc:        "gtコマンド",
			command:     &Command{raw: "gt"},
			commandType: CommandArithmetic,
			arg1:        "gt",
		},
		{
			desc:        "negコマンド",
			command:     &Command{raw: "neg"},
			commandType: CommandArithmetic,
			arg1:        "neg",
		},
		{
			desc:        "andコマンド",
			command:     &Command{raw: "and"},
			commandType: CommandArithmetic,
			arg1:        "and",
		},
		{
			desc:        "orコマンド",
			command:     &Command{raw: "or"},
			commandType: CommandArithmetic,
			arg1:        "or",
		},
		{
			desc:        "notコマンド",
			command:     &Command{raw: "not"},
			commandType: CommandArithmetic,
			arg1:        "not",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.command.Parse()
			got := tc.command

			if got.commandType != tc.commandType {
				t.Errorf("failed commandType: got = %d, want %d", got.commandType, tc.commandType)
			}

			if got.arg1 != tc.arg1 {
				t.Errorf("failed arg1: got = %s, want %s", got.arg1, tc.arg1)
			}
		})
	}
}

func TestCommandParse3(t *testing.T) {
	cases := []struct {
		desc        string
		command     *Command
		commandType CommandType
		arg1        string
		arg2        int
	}{
		{
			desc:        "pushコマンド",
			command:     &Command{raw: "push constant 7"},
			commandType: CommandPush,
			arg1:        "constant",
			arg2:        7,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.command.Parse()
			got := tc.command

			if got.commandType != tc.commandType {
				t.Errorf("failed commandType: got = %d, want %d", got.commandType, tc.commandType)
			}

			if got.arg1 != tc.arg1 {
				t.Errorf("failed arg1: got = %s, want %s", got.arg1, tc.arg1)
			}

			if *got.arg2 != tc.arg2 {
				t.Errorf("failed arg2: got = %d, want %d", *got.arg2, tc.arg2)
			}
		})
	}
}
