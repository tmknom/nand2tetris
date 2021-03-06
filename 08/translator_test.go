package main

import (
	"github.com/google/go-cmp/cmp"
	"strconv"

	"testing"
)

const (
	testPC  = 100
	testRaw = "dummy raw"
)

var testModuleName = "TestModule" // 定数だとアドレス参照できなかったのでvarで定義

func TestTranslatorsTranslateAll(t *testing.T) {
	cases := []struct {
		desc    string
		command *Command
		want    []string
	}{
		{
			desc: "初期化コードのみ",
			want: []string{
				"@256",
				"D=A",
				"@SP",
				"M=D",
				"@300",
				"D=A",
				"@LCL",
				"M=D",
				"@400",
				"D=A",
				"@ARG",
				"M=D",
				"@3000",
				"D=A",
				"@THIS",
				"M=D",
				"@3010",
				"D=A",
				"@THAT",
				"M=D",
				"@END",
				"0;JMP",
				"(TRUE)",
				"  D=-1",
				"  @R14",
				"  A=M",
				"  0;JMP",
				"(FALSE)",
				"  D=0",
				"  @R14",
				"  A=M",
				"  0;JMP",
				"(END)",
			},
		},
		{
			desc: "notコマンドを含む",
			command: &Command{
				raw:         "not",
				commandType: CommandArithmetic,
				arg1:        "not",
			},
			want: []string{
				"@256",
				"D=A",
				"@SP",
				"M=D",
				"@300",
				"D=A",
				"@LCL",
				"M=D",
				"@400",
				"D=A",
				"@ARG",
				"M=D",
				"@3000",
				"D=A",
				"@THIS",
				"M=D",
				"@3010",
				"D=A",
				"@THAT",
				"M=D",
				"@SP",
				"AM=M-1",
				"M=!M",
				"@SP",
				"M=M+1",
				"@END",
				"0;JMP",
				"(TRUE)",
				"  D=-1",
				"  @R14",
				"  A=M",
				"  0;JMP",
				"(FALSE)",
				"  D=0",
				"  @R14",
				"  A=M",
				"  0;JMP",
				"(END)",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translators := NewTranslators()
			if tc.command != nil {
				translators.Add(tc.command)
			}

			got := translators.TranslateAll()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}

func TestTranslatorsCalculatePC(t *testing.T) {
	cases := []struct {
		desc      string
		assembler []string
		pc        int
		want      int
	}{
		{
			desc: "ラベル定義がない",
			assembler: []string{
				"@SP",
				"AM=M-1",
				"D=M",
				"@SP",
				"AM=M-1",
				"M=M+D",
				"@SP",
				"M=M+1",
			},
			pc:   100,
			want: 108,
		},
		{
			desc: "ラベル定義のみ",
			assembler: []string{
				"(FooModule$Bar)",
			},
			pc:   100,
			want: 100,
		},
		{
			desc: "ラベル定義がある",
			assembler: []string{
				"(Math.max)",
				"@SP",
				"A=M",
				"M=0",
				"@SP",
				"M=M+1",
			},
			pc:   100,
			want: 105,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translators := NewTranslators()
			translators.pc = tc.pc
			translators.calculatePC(tc.assembler)

			got := translators.pc
			if got != tc.want {
				t.Errorf("failed assembler length: got = %d, want = %d", got, tc.want)
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
			translator := NewTranslator(testPC, testRaw, CommandArithmetic, tc.arg1, nil, &testModuleName)
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
			translator := NewTranslator(tc.pc, testRaw, CommandArithmetic, "eq", nil, &testModuleName)
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
				"@TestModule.14",
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
			translator := NewTranslator(testPC, testRaw, tc.commandType, tc.arg1, &tc.arg2, &testModuleName)
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
				"@R13",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@R13",
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
				"@R13",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@R13",
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
				"@R13",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@R13",
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
				"@R13",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@R13",
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
				"@TestModule.14",
				"M=D",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translator := NewTranslator(testPC, testRaw, tc.commandType, tc.arg1, &tc.arg2, &testModuleName)
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
			translator := NewTranslator(testPC, testRaw, tc.commandType, tc.arg1, nil, &tc.moduleName)
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
			translator := NewTranslator(testPC, testRaw, tc.commandType, tc.arg1, nil, &tc.moduleName)
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
			translator := NewTranslator(testPC, testRaw, tc.commandType, tc.arg1, nil, &tc.moduleName)
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
			translator := NewTranslator(testPC, testRaw, tc.commandType, tc.arg1, &tc.arg2, &testModuleName)
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
			translator := NewTranslator(testPC, testRaw, tc.commandType, tc.arg1, nil, &testModuleName)
			got := translator.Translate()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}

func TestTranslatorCall(t *testing.T) {
	cases := []struct {
		desc        string
		commandType CommandType
		arg1        string
		arg2        int
		pc          int
		want        []string
	}{
		{
			desc:        "call Math.max 0",
			commandType: CommandCall,
			arg1:        "Math.max",
			arg2:        0,
			pc:          100,
			want: []string{
				// push return-address
				"@RETURN-ADDRESS$TestModule$Math.max$100",
				"D=A",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",

				// push LCL
				"@LCL",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",

				// push ARG
				"@ARG",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",

				// push THIS
				"@THIS",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",

				// push THAT
				"@THAT",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",

				// @ARG = SP-n-5
				"@0",
				"D=A",
				"@5",
				"D=D+A",
				"@SP",
				"D=M-D",
				"@ARG",
				"M=D",

				// @LCL=SP
				"@SP",
				"D=M",
				"@LCL",
				"M=D",

				// goto f
				"@Math.max",
				"0;JMP",

				// (return-address)
				"(RETURN-ADDRESS$TestModule$Math.max$100)",
			},
		},
		{
			desc:        "call Math.min 2",
			commandType: CommandCall,
			arg1:        "Math.min",
			arg2:        2,
			pc:          397,
			want: []string{
				// push return-address
				"@RETURN-ADDRESS$TestModule$Math.min$397",
				"D=A",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",

				// push LCL
				"@LCL",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",

				// push ARG
				"@ARG",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",

				// push THIS
				"@THIS",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",

				// push THAT
				"@THAT",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",

				// @ARG = SP-n-5
				"@2",
				"D=A",
				"@5",
				"D=D+A",
				"@SP",
				"D=M-D",
				"@ARG",
				"M=D",

				// @LCL=SP
				"@SP",
				"D=M",
				"@LCL",
				"M=D",

				// goto f
				"@Math.min",
				"0;JMP",

				// (return-address)
				"(RETURN-ADDRESS$TestModule$Math.min$397)",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			translator := NewTranslator(tc.pc, testRaw, tc.commandType, tc.arg1, &tc.arg2, &testModuleName)
			got := translator.Translate()

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("failed %s: diff (-got +want):\n%s", tc.desc, diff)
			}
		})
	}
}
