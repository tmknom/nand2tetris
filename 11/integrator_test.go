package main

import (
	"./parsing"
	"./symbol"
	"bufio"
	"github.com/google/go-cmp/cmp"
	"os"
	"strings"
	"testing"
)

func SetupTestForIntegrator() {
	// デバッグフラグの無効化
	parsing.DebugCode = false
	symbol.DebugSymbolTables = false

	// シンボルテーブルの初期化
	symbol.GlobalSymbolTables.Reset("Testing")
	symbol.GlobalSymbolTables.ResetSubroutine("TestRun")

	// ID生成器の初期化
	symbol.GlobalIdGenerator.Reset()
}

func TestIntegratorGenerate(t *testing.T) {
	cases := []struct {
		desc string
		src  []string
		dest []string
		want []string
	}{
		{
			desc: "Seven",
			src: []string{
				"Fixture/Seven/Main.jack",
			},
			dest: []string{
				"Fixture/Seven/Main.vm",
			},
			want: []string{
				"Fixture/Seven/cmp/Main.vm",
			},
		},
		{
			desc: "ConvertToBin",
			src: []string{
				"Fixture/ConvertToBin/Main.jack",
			},
			dest: []string{
				"Fixture/ConvertToBin/Main.vm",
			},
			want: []string{
				"Fixture/ConvertToBin/cmp/Main.vm",
			},
		},
		{
			desc: "Square",
			src: []string{
				"Fixture/Square/Main.jack",
				"Fixture/Square/Square.jack",
				"Fixture/Square/SquareGame.jack",
			},
			dest: []string{
				"Fixture/Square/Main.vm",
				"Fixture/Square/Square.vm",
				"Fixture/Square/SquareGame.vm",
			},
			want: []string{
				"Fixture/Square/cmp/Main.vm",
				"Fixture/Square/cmp/Square.vm",
				"Fixture/Square/cmp/SquareGame.vm",
			},
		},
	}

	for _, tc := range cases {
		for _, file := range tc.dest {
			os.Remove(file)
		}

		t.Run(tc.desc, func(t *testing.T) {
			// いろいろ初期化
			SetupTestForIntegrator()

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

func TestIntegratorIntegrate(t *testing.T) {
	cases := []struct {
		desc             string
		src              []string
		destTokenizedXML []string
		destParsedXML    []string
		wantTokenizedXML []string
		wantParsedXML    []string
	}{
		{
			desc: "ExpressionLessSquare",
			src: []string{
				"Fixture/ExpressionLessSquare/Main.jack",
				"Fixture/ExpressionLessSquare/Square.jack",
				"Fixture/ExpressionLessSquare/SquareGame.jack",
			},
			destTokenizedXML: []string{
				"Fixture/ExpressionLessSquare/MainT.xml",
				"Fixture/ExpressionLessSquare/SquareT.xml",
				"Fixture/ExpressionLessSquare/SquareGameT.xml",
			},
			destParsedXML: []string{
				"Fixture/ExpressionLessSquare/Main.xml",
				"Fixture/ExpressionLessSquare/Square.xml",
				"Fixture/ExpressionLessSquare/SquareGame.xml",
			},
			wantTokenizedXML: []string{
				"Fixture/ExpressionLessSquare/cmp/MainT.xml",
				"Fixture/ExpressionLessSquare/cmp/SquareT.xml",
				"Fixture/ExpressionLessSquare/cmp/SquareGameT.xml",
			},
			wantParsedXML: []string{
				"Fixture/ExpressionLessSquare/cmp/Main.xml",
				"Fixture/ExpressionLessSquare/cmp/Square.xml",
				"Fixture/ExpressionLessSquare/cmp/SquareGame.xml",
			},
		},
		{
			desc: "ArrayTest",
			src: []string{
				"Fixture/ArrayTest/Main.jack",
			},
			destTokenizedXML: []string{
				"Fixture/ArrayTest/MainT.xml",
			},
			destParsedXML: []string{
				"Fixture/ArrayTest/Main.xml",
			},
			wantTokenizedXML: []string{
				"Fixture/ArrayTest/cmp/MainT.xml",
			},
			wantParsedXML: []string{
				"Fixture/ArrayTest/cmp/Main.xml",
			},
		},
		{
			desc: "SquareVersion10",
			src: []string{
				"Fixture/SquareVersion10/Main.jack",
				"Fixture/SquareVersion10/Square.jack",
				"Fixture/SquareVersion10/SquareGame.jack",
			},
			destTokenizedXML: []string{
				"Fixture/SquareVersion10/MainT.xml",
				"Fixture/SquareVersion10/SquareT.xml",
				"Fixture/SquareVersion10/SquareGameT.xml",
			},
			destParsedXML: []string{
				"Fixture/SquareVersion10/Main.xml",
				"Fixture/SquareVersion10/Square.xml",
				"Fixture/SquareVersion10/SquareGame.xml",
			},
			wantTokenizedXML: []string{
				"Fixture/SquareVersion10/cmp/MainT.xml",
				"Fixture/SquareVersion10/cmp/SquareT.xml",
				"Fixture/SquareVersion10/cmp/SquareGameT.xml",
			},
			wantParsedXML: []string{
				"Fixture/SquareVersion10/cmp/Main.xml",
				"Fixture/SquareVersion10/cmp/Square.xml",
				"Fixture/SquareVersion10/cmp/SquareGame.xml",
			},
		},
	}

	for _, tc := range cases {
		for _, file := range tc.destTokenizedXML {
			os.Remove(file)
		}
		for _, file := range tc.destParsedXML {
			os.Remove(file)
		}

		t.Run(tc.desc, func(t *testing.T) {
			integrator := NewIntegrator(tc.src)
			integrator.Integrate()

			for i, dest := range tc.destTokenizedXML {
				got := readFileQuietly(dest)
				want := readFileQuietly(tc.wantTokenizedXML[i])

				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("failed: diff %s %s: (-got +want):\n%s\n", dest, tc.wantTokenizedXML[i], diff)
				} else {
					os.Remove(dest)
				}
			}

			for i, dest := range tc.destParsedXML {
				got := readFileQuietly(dest)
				want := readFileQuietly(tc.wantParsedXML[i])

				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("failed: diff %s %s: (-got +want):\n%s\n", dest, tc.wantTokenizedXML[i], diff)
				} else {
					os.Remove(dest)
				}
			}
		})
	}
}

func TestIntegratorIntegrateFile(t *testing.T) {
	cases := []struct {
		desc             string
		src              string
		destTokenizedXML string
		destParsedXML    string
		wantTokenizedXML string
		wantParsedXML    string
	}{
		{
			desc:             "指定したファイルのトークナイズが実行できる",
			src:              "Fixture/ExpressionLessSquare/Square.jack",
			destTokenizedXML: "Fixture/ExpressionLessSquare/SquareT.xml",
			destParsedXML:    "Fixture/ExpressionLessSquare/Square.xml",
			wantTokenizedXML: "Fixture/ExpressionLessSquare/cmp/SquareT.xml",
			wantParsedXML:    "Fixture/ExpressionLessSquare/cmp/Square.xml",
		},
	}

	for _, tc := range cases {
		os.Remove(tc.destTokenizedXML)
		os.Remove(tc.destParsedXML)

		t.Run(tc.desc, func(t *testing.T) {
			integrator := NewIntegrator([]string{})
			integrator.integrateFile(tc.src)

			gotTokenizedXML := readFileQuietly(tc.destTokenizedXML)
			wantTokenizedXML := readFileQuietly(tc.wantTokenizedXML)

			if diff := cmp.Diff(gotTokenizedXML, wantTokenizedXML); diff != "" {
				t.Errorf("failed: diff %s %s: (-got +want):\n%s\n", tc.destTokenizedXML, tc.wantTokenizedXML, diff)
			} else {
				os.Remove(tc.destTokenizedXML)
			}

			gotParsedXML := readFileQuietly(tc.destParsedXML)
			wantParsedXML := readFileQuietly(tc.wantParsedXML)

			if diff := cmp.Diff(gotParsedXML, wantParsedXML); diff != "" {
				t.Errorf("failed: diff %s %s: (-got +want):\n%s\n", tc.destParsedXML, tc.wantParsedXML, diff)
			} else {
				os.Remove(tc.destParsedXML)
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
		line := strings.TrimSpace(scanner.Text())
		lines = append(lines, line)
	}
	return lines
}
