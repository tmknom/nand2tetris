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
	}

	for _, tc := range cases {
		os.Remove(tc.destFile)
		t.Run(tc.desc, func(t *testing.T) {
			integrator := NewIntegrator(tc.srcFile)
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
