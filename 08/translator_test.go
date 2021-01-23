package main

import (
	"github.com/google/go-cmp/cmp"
	"strconv"

	"strings"
	"testing"
)

const (
	testFilename = "Dummy.vm"
	testPC       = 100
)

var testModuleName = "TestModule" // 定数だとアドレス参照できなかったのでvarで定義

func TestNewTranslators(t *testing.T) {
	cases := []struct {
		desc     string
		filename string
		want     string
	}{
		{
			desc:     "TODO",
			filename: "StackArithmetic/SimpleAdd/Test.vm",
			want:     "Test",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			dest := NewTranslators(tc.filename, false)
			got := dest.moduleName
			if got != tc.want {
				t.Errorf("failed: got = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestTranslatorsTranslateAll(t *testing.T) {
	cases := []struct {
		desc    string
		hasInit HasInit
		want    int
	}{
		{
			desc:    "hasInitがtrue",
			hasInit: true,
			want:    34,
		},
		{
			desc:    "hasInitがfalse",
			hasInit: false,
			want:    33,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translators := NewTranslators(testFilename, tc.hasInit)
			assembler := translators.TranslateAll()

			var containInit HasInit = false
			for _, line := range assembler {
				if strings.Contains(line, "@Sys.init") {
					containInit = true
					break
				}
			}

			if containInit != tc.hasInit {
				t.Errorf("failed hasInit: got = %t, want = %t", containInit, tc.hasInit)
			}

			length := len(assembler)
			if length != tc.want {
				t.Errorf("failed assembler length: got = %d, hasInit = %d", length, tc.want)
			}
		})
	}
}

func TestTranslatorArithmetic(t *testing.T) {
	cases := []struct {
		desc string
		arg1 string
		want []string
	}{
		{
			desc: "add",
			arg1: "add",
			want: []string{
				"@SP",
				"AM=M-1",
				"D=M",
				"@SP",
				"AM=M-1",
				"M=M+D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc: "sub",
			arg1: "sub",
			want: []string{
				"@SP",
				"AM=M-1",
				"D=M",
				"@SP",
				"AM=M-1",
				"M=M-D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc: "neg",
			arg1: "neg",
			want: []string{
				"@SP",
				"AM=M-1",
				"M=-M",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc: "eq",
			arg1: "eq",
			want: []string{
				"@114",
				"D=A",
				"@R14",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@SP",
				"AM=M-1",
				"D=M-D",
				"@TRUE",
				"D;JEQ",
				"@FALSE",
				"0;JMP",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc: "lt",
			arg1: "lt",
			want: []string{
				"@114",
				"D=A",
				"@R14",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@SP",
				"AM=M-1",
				"D=M-D",
				"@TRUE",
				"D;JLT",
				"@FALSE",
				"0;JMP",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc: "gt",
			arg1: "gt",
			want: []string{
				"@114",
				"D=A",
				"@R14",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@SP",
				"AM=M-1",
				"D=M-D",
				"@TRUE",
				"D;JGT",
				"@FALSE",
				"0;JMP",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc: "and",
			arg1: "and",
			want: []string{
				"@SP",
				"AM=M-1",
				"D=M",
				"@SP",
				"AM=M-1",
				"M=D&M",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc: "or",
			arg1: "or",
			want: []string{
				"@SP",
				"AM=M-1",
				"D=M",
				"@SP",
				"AM=M-1",
				"M=D|M",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc: "not",
			arg1: "not",
			want: []string{
				"@SP",
				"AM=M-1",
				"M=!M",
				"@SP",
				"M=M+1",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translator := NewTranslator(testPC, CommandArithmetic, tc.arg1, nil, &testModuleName)
			got := translator.Translate()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}

func TestTranslatorCompareBinary(t *testing.T) {
	cases := []struct {
		desc            string
		pc              int
		returnAddress   int
		beforeJumpStep  []string
		afterReturnStep []string
	}{
		{
			desc:          "プログラムカウンタがゼロの場合",
			pc:            0,
			returnAddress: 14,
			beforeJumpStep: []string{
				"@14",
				"D=A",
				"@R14",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@SP",
				"AM=M-1",
				"D=M-D",
				"@TRUE",
				"D;JEQ",
				"@FALSE",
				"0;JMP",
			},
			afterReturnStep: []string{
				// Dレジスタに -1 がセットされる
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc:          "プログラムカウンタがゼロ以外の場合",
			pc:            50,
			returnAddress: 64,
			beforeJumpStep: []string{
				"@64",
				"D=A",
				"@R14",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@SP",
				"AM=M-1",
				"D=M-D",
				"@TRUE",
				"D;JEQ",
				"@FALSE",
				"0;JMP",
			},
			afterReturnStep: []string{
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translator := NewTranslator(tc.pc, CommandArithmetic, "eq", nil, &testModuleName)
			got := translator.Translate()

			gotReturnAddress, _ := strconv.Atoi((got[0])[1:])
			if gotReturnAddress != tc.returnAddress {
				t.Errorf("failed returnAddress: got = %d, want %d", gotReturnAddress, tc.returnAddress)
			}

			beforeJumpIndex := tc.returnAddress - tc.pc
			if diff := cmp.Diff(got[:beforeJumpIndex], tc.beforeJumpStep); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}

			afterReturnIndex := tc.returnAddress - tc.pc
			if diff := cmp.Diff(got[afterReturnIndex:], tc.afterReturnStep); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}

func TestTranslatorPush(t *testing.T) {
	cases := []struct {
		desc        string
		commandType CommandType
		arg1        string
		arg2        int
		want        []string
	}{
		{
			desc:        "push constant 100",
			commandType: CommandPush,
			arg1:        "constant",
			arg2:        100,
			want: []string{
				"@100",
				"D=A",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc:        "push local 10",
			commandType: CommandPush,
			arg1:        "local",
			arg2:        10,
			want: []string{
				"@10",
				"D=A",
				"@LCL",
				"A=D+M",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc:        "push argument 10",
			commandType: CommandPush,
			arg1:        "argument",
			arg2:        11,
			want: []string{
				"@11",
				"D=A",
				"@ARG",
				"A=D+M",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc:        "push this 12",
			commandType: CommandPush,
			arg1:        "this",
			arg2:        12,
			want: []string{
				"@12",
				"D=A",
				"@THIS",
				"A=D+M",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc:        "push that 13",
			commandType: CommandPush,
			arg1:        "that",
			arg2:        13,
			want: []string{
				"@13",
				"D=A",
				"@THAT",
				"A=D+M",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc:        "push temp 6",
			commandType: CommandPush,
			arg1:        "temp",
			arg2:        6,
			want: []string{
				"@11",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc:        "push pointer 1",
			commandType: CommandPush,
			arg1:        "pointer",
			arg2:        1,
			want: []string{
				"@4",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
		{
			desc:        "push static 14",
			commandType: CommandPush,
			arg1:        "static",
			arg2:        14,
			want: []string{
				"@30",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translator := NewTranslator(testPC, tc.commandType, tc.arg1, &tc.arg2, &testModuleName)
			got := translator.Translate()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}

func TestTranslatorPop(t *testing.T) {
	cases := []struct {
		desc        string
		commandType CommandType
		arg1        string
		arg2        int
		want        []string
	}{
		{
			desc:        "pop local 10",
			commandType: CommandPop,
			arg1:        "local",
			arg2:        10,
			want: []string{
				"@10",
				"D=A",
				"@LCL",
				"D=D+M",
				"@R14",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@R14",
				"A=M",
				"M=D",
			},
		},
		{
			desc:        "pop argument 11",
			commandType: CommandPop,
			arg1:        "argument",
			arg2:        11,
			want: []string{
				"@11",
				"D=A",
				"@ARG",
				"D=D+M",
				"@R14",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@R14",
				"A=M",
				"M=D",
			},
		},
		{
			desc:        "pop this 12",
			commandType: CommandPop,
			arg1:        "this",
			arg2:        12,
			want: []string{
				"@12",
				"D=A",
				"@THIS",
				"D=D+M",
				"@R14",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@R14",
				"A=M",
				"M=D",
			},
		},
		{
			desc:        "pop that 13",
			commandType: CommandPop,
			arg1:        "that",
			arg2:        13,
			want: []string{
				"@13",
				"D=A",
				"@THAT",
				"D=D+M",
				"@R14",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@R14",
				"A=M",
				"M=D",
			},
		},
		{
			desc:        "pop temp 6",
			commandType: CommandPop,
			arg1:        "temp",
			arg2:        6,
			want: []string{
				"@SP",
				"AM=M-1",
				"D=M",
				"@11",
				"M=D",
			},
		},
		{
			desc:        "pop pointer 1",
			commandType: CommandPop,
			arg1:        "pointer",
			arg2:        1,
			want: []string{
				"@SP",
				"AM=M-1",
				"D=M",
				"@4",
				"M=D",
			},
		},
		{
			desc:        "pop static 14",
			commandType: CommandPop,
			arg1:        "static",
			arg2:        14,
			want: []string{
				"@SP",
				"AM=M-1",
				"D=M",
				"@30",
				"M=D",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translator := NewTranslator(testPC, tc.commandType, tc.arg1, &tc.arg2, &testModuleName)
			got := translator.Translate()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}

func TestTranslatorLabel(t *testing.T) {
	cases := []struct {
		desc        string
		commandType CommandType
		arg1        string
		moduleName  string
		want        []string
	}{
		{
			desc:        "label Bar",
			commandType: CommandLabel,
			arg1:        "Bar",
			moduleName:  "FooModule",
			want: []string{
				"(FooModule$Bar)",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translator := NewTranslator(testPC, tc.commandType, tc.arg1, nil, &tc.moduleName)
			got := translator.Translate()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}

func TestTranslatorLabelGoto(t *testing.T) {
	cases := []struct {
		desc        string
		commandType CommandType
		arg1        string
		moduleName  string
		want        []string
	}{
		{
			desc:        "goto Bar",
			commandType: CommandGoto,
			arg1:        "Bar",
			moduleName:  "FooModule",
			want: []string{
				"@FooModule$Bar",
				"0;JMP",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translator := NewTranslator(testPC, tc.commandType, tc.arg1, nil, &tc.moduleName)
			got := translator.Translate()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}

func TestTranslatorIfGoto(t *testing.T) {
	cases := []struct {
		desc        string
		commandType CommandType
		arg1        string
		moduleName  string
		want        []string
	}{
		{
			desc:        "if-goto Bar",
			commandType: CommandIf,
			arg1:        "Bar",
			moduleName:  "FooModule",
			want: []string{
				"@SP",
				"AM=M-1",
				"D=M",
				"@FooModule$Bar",
				"D;JNE",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translator := NewTranslator(testPC, tc.commandType, tc.arg1, nil, &tc.moduleName)
			got := translator.Translate()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}

func TestTranslatorFunction(t *testing.T) {
	cases := []struct {
		desc        string
		commandType CommandType
		arg1        string
		arg2        int
		want        []string
	}{
		{
			desc:        "function Math.max 2",
			commandType: CommandFunction,
			arg1:        "Math.max",
			arg2:        2,
			want: []string{
				"(Math.max)",
				"@SP",
				"A=M",
				"M=0",
				"@SP",
				"M=M+1",
				"@SP",
				"A=M",
				"M=0",
				"@SP",
				"M=M+1",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translator := NewTranslator(testPC, tc.commandType, tc.arg1, &tc.arg2, &testModuleName)
			got := translator.Translate()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}

func TestTranslatorReturnFunction(t *testing.T) {
	cases := []struct {
		desc        string
		commandType CommandType
		arg1        string
		want        []string
	}{
		{
			desc:        "return",
			commandType: CommandReturn,
			arg1:        "return",
			want: []string{
				// FRAME=LCL
				"@LCL",
				"D=M",
				"@R13",
				"M=D",

				// RET = *(FRAME-5)
				"@R13",
				"D=M",
				"@5",
				"A=D-A",
				"D=M",
				"@R14",
				"M=D",

				// *ARG = pop()
				"@SP",
				"AM=M-1",
				"D=M",
				"@ARG",
				"A=M",
				"M=D",

				// SP = ARG+1
				"D=A",
				"@SP",
				"M=D+1",

				// THAT = *(FRAME-1)
				"@R13",
				"D=M",
				"@1",
				"A=D-A",
				"D=M",
				"@THAT",
				"M=D",

				// THIS = *(FRAME-2)
				"@R13",
				"D=M",
				"@2",
				"A=D-A",
				"D=M",
				"@THIS",
				"M=D",

				// ARG = *(FRAME-3)
				"@R13",
				"D=M",
				"@3",
				"A=D-A",
				"D=M",
				"@ARG",
				"M=D",

				// LCL = *(FRAME-4)
				"@R13",
				"D=M",
				"@4",
				"A=D-A",
				"D=M",
				"@LCL",
				"M=D",

				// goto RET
				"@R14",
				"A=M",
				"0;JMP",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translator := NewTranslator(testPC, tc.commandType, tc.arg1, nil, &testModuleName)
			got := translator.Translate()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}
