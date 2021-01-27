package main

import (
	"bufio"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

func TestIntegratorIntegrate(t *testing.T) {
	cases := []struct {
		desc string
		src  []string
		dest []string
		want []string
	}{
		{
			desc: "Square",
			src: []string{
				"Square/Main.jack",
				"Square/Square.jack",
				"Square/SquareGame.jack",
			},
			dest: []string{
				"Square/MainT.xml",
				"Square/SquareT.xml",
				"Square/SquareGameT.xml",
			},
			want: []string{
				"Square/cmp/MainT.xml",
				"Square/cmp/SquareT.xml",
				"Square/cmp/SquareGameT.xml",
			},
		},
		{
			desc: "ArrayTest",
			src: []string{
				"ArrayTest/Main.jack",
			},
			dest: []string{
				"ArrayTest/MainT.xml",
			},
			want: []string{
				"ArrayTest/cmp/MainT.xml",
			},
		},
	}

	for _, tc := range cases {
		for _, dest := range tc.dest {
			os.Remove(dest)
		}

		t.Run(tc.desc, func(t *testing.T) {
			integrator := NewIntegrator(tc.src)
			integrator.Integrate()

			for i, dest := range tc.dest {
				got := readFileQuietly(dest)
				want := readFileQuietly(tc.want[i])

				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("failed: diff %s %s: (-got +want):\n%s\n", dest, tc.want[i], diff)
				} else {
					os.Remove(dest)
				}
			}
		})
	}
}

func TestIntegratorIntegrateFile(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		dest string
		want string
	}{
		{
			desc: "指定したファイルのトークナイズが実行できる",
			src:  "Square/Square.jack",
			dest: "Square/SquareT.xml",
			want: "Square/cmp/SquareT.xml",
		},
	}

	for _, tc := range cases {
		os.Remove(tc.dest)
		t.Run(tc.desc, func(t *testing.T) {
			integrator := NewIntegrator([]string{})
			integrator.integrateFile(tc.src)

			got := readFileQuietly(tc.dest)
			want := readFileQuietly(tc.want)

			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("failed: diff %s %s: (-got +want):\n%s\n", tc.dest, tc.want, diff)
			} else {
				os.Remove(tc.dest)
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
