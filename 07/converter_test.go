package main

import (
	"fmt"
	"reflect"
	"testing"
)

const testPC = 100

func TestConverterArithmetic(t *testing.T) {
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
				"@R15",
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
				"D;JNE",
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
				"@R15",
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
				"D;JGE",
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
				"@R15",
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
				"D;JLE",
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
			converter := NewConverter(testPC, CommandArithmetic, tc.arg1, nil)
			got := converter.Convert()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("failed:\ngot = %s,\nwant = %s", prettySlice(got), prettySlice(tc.want))
			}
		})
	}
}

func TestConverterPush(t *testing.T) {
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
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			converter := NewConverter(testPC, tc.commandType, tc.arg1, &tc.arg2)
			got := converter.Convert()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("failed:\ngot = %s,\nwant = %s", prettySlice(got), prettySlice(tc.want))
			}
		})
	}
}

func TestConverterPop(t *testing.T) {
	cases := []struct {
		desc        string
		commandType CommandType
		arg1        string
		arg2        int
		want        []string
	}{
		{
			desc:        "pop local 20",
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
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			converter := NewConverter(testPC, tc.commandType, tc.arg1, &tc.arg2)
			got := converter.Convert()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("failed:\ngot = %s,\nwant = %s", prettySlice(got), prettySlice(tc.want))
			}
		})
	}
}

func prettySlice(list []string) string {
	contents := []string{}
	for i, element := range list {
		pretty := fmt.Sprintf("  %d: %s,\n", i, element)
		contents = append(contents, pretty)
	}
	return fmt.Sprintf("\n%s", contents)
}
