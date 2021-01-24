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
		desc     string
		srcFile  string
		destFile string
		wantFile string
	}{
		{
			desc:     "SimpleAdd",
			srcFile:  "StackArithmetic/SimpleAdd/SimpleAdd.vm",
			destFile: "StackArithmetic/SimpleAdd/SimpleAdd.asm",
			wantFile: "StackArithmetic/SimpleAdd/SimpleAdd.asm.cmp",
		},
		{
			desc:     "StackTest",
			srcFile:  "StackArithmetic/StackTest/StackTest.vm",
			destFile: "StackArithmetic/StackTest/StackTest.asm",
			wantFile: "StackArithmetic/StackTest/StackTest.asm.cmp",
		},
		{
			desc:     "BasicTest",
			srcFile:  "MemoryAccess/BasicTest/BasicTest.vm",
			destFile: "MemoryAccess/BasicTest/BasicTest.asm",
			wantFile: "MemoryAccess/BasicTest/BasicTest.asm.cmp",
		},
		{
			desc:     "PointerTest",
			srcFile:  "MemoryAccess/PointerTest/PointerTest.vm",
			destFile: "MemoryAccess/PointerTest/PointerTest.asm",
			wantFile: "MemoryAccess/PointerTest/PointerTest.asm.cmp",
		},
		{
			desc:     "StaticTest",
			srcFile:  "MemoryAccess/StaticTest/StaticTest.vm",
			destFile: "MemoryAccess/StaticTest/StaticTest.asm",
			wantFile: "MemoryAccess/StaticTest/StaticTest.asm.cmp",
		},
		{
			desc:     "BasicLoop",
			srcFile:  "ProgramFlow/BasicLoop/BasicLoop.vm",
			destFile: "ProgramFlow/BasicLoop/BasicLoop.asm",
			wantFile: "ProgramFlow/BasicLoop/BasicLoop.asm.cmp",
		},
		{
			desc:     "FibonacciSeries",
			srcFile:  "ProgramFlow/FibonacciSeries/FibonacciSeries.vm",
			destFile: "ProgramFlow/FibonacciSeries/FibonacciSeries.asm",
			wantFile: "ProgramFlow/FibonacciSeries/FibonacciSeries.asm.cmp",
		},
		{
			desc:     "SimpleFunction",
			srcFile:  "FunctionCalls/SimpleFunction/SimpleFunction.vm",
			destFile: "FunctionCalls/SimpleFunction/SimpleFunction.asm",
			wantFile: "FunctionCalls/SimpleFunction/SimpleFunction.asm.cmp",
		},
		{
			desc:     "FibonacciElement",
			srcFile:  "FunctionCalls/FibonacciElement/FibonacciElement.vm",
			destFile: "FunctionCalls/FibonacciElement/FibonacciElement.asm",
			wantFile: "FunctionCalls/FibonacciElement/FibonacciElement.asm.cmp",
		},
	}

	for _, tc := range cases {
		os.Remove(tc.destFile)
		t.Run(tc.desc, func(t *testing.T) {
			integrator := NewIntegrator([]string{tc.srcFile})
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
