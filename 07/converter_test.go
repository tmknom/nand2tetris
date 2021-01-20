package main

import (
	"fmt"
	"reflect"
	"testing"
)

const testPC = 100

func TestConverterConvertArithmetic(t *testing.T) {
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
				"M=D+M",
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
				"@EQ",
				"D;JEQ",
				"@NEQ",
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
				"@EQ",
				"D;JLT",
				"@NEQ",
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
				"@EQ",
				"D;JGT",
				"@NEQ",
				"D;JLE",
				"@SP",
				"M=M+1",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			converter := NewConverter(testPC, CommandArithmetic, tc.arg1, nil)
			got := converter.convertArithmetic()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("failed:\ngot = %s,\nwant = %s", prettySlice(got), prettySlice(tc.want))
			}
		})
	}
}

func TestConverterConvertPush(t *testing.T) {
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
			got := converter.convertPush()
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
