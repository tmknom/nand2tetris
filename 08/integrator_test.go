package main

import (
	"bufio"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestIntegratorIntegrate(t *testing.T) {
	cases := []struct {
		desc      string
		arg       string
		filenames []string
		destFile  string
		wantFile  string
	}{
		{
			desc:      "SimpleAdd",
			arg:       "StackArithmetic/SimpleAdd/SimpleAdd.vm",
			filenames: []string{"StackArithmetic/SimpleAdd/SimpleAdd.vm"},
			destFile:  "StackArithmetic/SimpleAdd/SimpleAdd.asm",
			wantFile:  "StackArithmetic/SimpleAdd/SimpleAdd.asm.cmp",
		},
		{
			desc:      "StackTest",
			arg:       "StackArithmetic/StackTest/StackTest.vm",
			filenames: []string{"StackArithmetic/StackTest/StackTest.vm"},
			destFile:  "StackArithmetic/StackTest/StackTest.asm",
			wantFile:  "StackArithmetic/StackTest/StackTest.asm.cmp",
		},
		{
			desc:      "BasicTest",
			arg:       "MemoryAccess/BasicTest/BasicTest.vm",
			filenames: []string{"MemoryAccess/BasicTest/BasicTest.vm"},
			destFile:  "MemoryAccess/BasicTest/BasicTest.asm",
			wantFile:  "MemoryAccess/BasicTest/BasicTest.asm.cmp",
		},
		{
			desc:      "PointerTest",
			arg:       "MemoryAccess/PointerTest/PointerTest.vm",
			filenames: []string{"MemoryAccess/PointerTest/PointerTest.vm"},
			destFile:  "MemoryAccess/PointerTest/PointerTest.asm",
			wantFile:  "MemoryAccess/PointerTest/PointerTest.asm.cmp",
		},
		{
			desc:      "StaticTest",
			arg:       "MemoryAccess/StaticTest/StaticTest.vm",
			filenames: []string{"MemoryAccess/StaticTest/StaticTest.vm"},
			destFile:  "MemoryAccess/StaticTest/StaticTest.asm",
			wantFile:  "MemoryAccess/StaticTest/StaticTest.asm.cmp",
		},
		{
			desc:      "BasicLoop",
			arg:       "ProgramFlow/BasicLoop/BasicLoop.vm",
			filenames: []string{"ProgramFlow/BasicLoop/BasicLoop.vm"},
			destFile:  "ProgramFlow/BasicLoop/BasicLoop.asm",
			wantFile:  "ProgramFlow/BasicLoop/BasicLoop.asm.cmp",
		},
		{
			desc:      "FibonacciSeries",
			arg:       "ProgramFlow/FibonacciSeries/FibonacciSeries.vm",
			filenames: []string{"ProgramFlow/FibonacciSeries/FibonacciSeries.vm"},
			destFile:  "ProgramFlow/FibonacciSeries/FibonacciSeries.asm",
			wantFile:  "ProgramFlow/FibonacciSeries/FibonacciSeries.asm.cmp",
		},
		{
			desc:      "SimpleFunction",
			arg:       "FunctionCalls/SimpleFunction/SimpleFunction.vm",
			filenames: []string{"FunctionCalls/SimpleFunction/SimpleFunction.vm"},
			destFile:  "FunctionCalls/SimpleFunction/SimpleFunction.asm",
			wantFile:  "FunctionCalls/SimpleFunction/SimpleFunction.asm.cmp",
		},
		{
			desc: "FibonacciElement",
			arg:  "FunctionCalls/FibonacciElement",
			filenames: []string{
				"FunctionCalls/FibonacciElement/Main.vm",
				"FunctionCalls/FibonacciElement/Sys.vm",
			},
			destFile: "FunctionCalls/FibonacciElement/FibonacciElement.asm",
			wantFile: "FunctionCalls/FibonacciElement/FibonacciElement.asm.cmp",
		},
		{
			desc:      "NestedCall",
			arg:       "FunctionCalls/NestedCall",
			filenames: []string{"FunctionCalls/NestedCall/Sys.vm"},
			destFile:  "FunctionCalls/NestedCall/NestedCall.asm",
			wantFile:  "FunctionCalls/NestedCall/NestedCall.asm.cmp",
		},
	}

	for _, tc := range cases {
		os.Remove(tc.destFile)
		t.Run(tc.desc, func(t *testing.T) {
			integrator := NewIntegrator(tc.filenames, tc.arg)
			integrator.Integrate()

			want := readFileQuietly(tc.wantFile)
			got := readFileQuietly(tc.destFile)

			if !reflect.DeepEqual(got, want) {
				diff, _ := exec.Command("diff", tc.wantFile, tc.destFile).Output()
				t.Errorf("failed %s: diff = \n%s", tc.desc, diff)
			} else {
				os.Remove(tc.destFile)
			}
		})
	}
}

func readFileQuietly(filename string) []string {
	file, _ := os.Open(filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}
