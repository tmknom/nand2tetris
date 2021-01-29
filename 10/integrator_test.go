package main

import (
	"bufio"
	"github.com/google/go-cmp/cmp"
	"os"
	"strings"
	"testing"
)

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
			desc: "Square",
			src: []string{
				"Fixture/Square/Main.jack",
				"Fixture/Square/Square.jack",
				"Fixture/Square/SquareGame.jack",
			},
			destTokenizedXML: []string{
				"Fixture/Square/MainT.xml",
				"Fixture/Square/SquareT.xml",
				"Fixture/Square/SquareGameT.xml",
			},
			destParsedXML: []string{
				"Fixture/Square/Main.xml",
				"Fixture/Square/Square.xml",
				"Fixture/Square/SquareGame.xml",
			},
			wantTokenizedXML: []string{
				"Fixture/Square/cmp/MainT.xml",
				"Fixture/Square/cmp/SquareT.xml",
				"Fixture/Square/cmp/SquareGameT.xml",
			},
			wantParsedXML: []string{
				"Fixture/Square/cmp/Main.xml",
				"Fixture/Square/cmp/Square.xml",
				"Fixture/Square/cmp/SquareGame.xml",
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

			for _, dest := range tc.destParsedXML {
				//got := readFileQuietly(dest)
				//want := readFileQuietly(tc.wantParsedXML[i])
				//
				//if diff := cmp.Diff(got, want); diff != "" {
				//	t.Errorf("failed: diff %s %s: (-got +want):\n%s\n", dest, tc.wantTokenizedXML[i], diff)
				//} else {
				//	os.Remove(dest)
				//}
				os.Remove(dest)
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
			src:              "Fixture/Square/Square.jack",
			destTokenizedXML: "Fixture/Square/SquareT.xml",
			destParsedXML:    "Fixture/Square/Square.xml",
			wantTokenizedXML: "Fixture/Square/cmp/SquareT.xml",
			wantParsedXML:    "Fixture/Square/cmp/Square.xml",
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

			//gotParsedXML := readFileQuietly(tc.destParsedXML)
			//wantParsedXML := readFileQuietly(tc.wantParsedXML)
			//
			//if diff := cmp.Diff(gotParsedXML, wantParsedXML); diff != "" {
			//	t.Errorf("failed: diff %s %s: (-got +want):\n%s\n", tc.destParsedXML, tc.wantParsedXML, diff)
			//} else {
			//	os.Remove(tc.destParsedXML)
			//}
			os.Remove(tc.destParsedXML)
		})
	}
}

// TODO パース処理が一通り実装できるまで、途中状態のXMLファイルをリグレッションテストできるようにしておく
func TestIntegratorProvisional(t *testing.T) {
	cases := []struct {
		desc             string
		src              string
		destTokenizedXML string
		destParsedXML    string
		wantParsedXML    string
	}{
		{
			desc:             "途中状態のXMLファイルをリグレッションテスト",
			src:              "Fixture/ExpressionLessSquare/Main.jack",
			destTokenizedXML: "Fixture/ExpressionLessSquare/MainT.xml",
			destParsedXML:    "Fixture/ExpressionLessSquare/Main.xml",
			wantParsedXML:    "Fixture/ExpressionLessSquare/provisional/Main.xml",
		},
	}

	for _, tc := range cases {
		os.Remove(tc.destTokenizedXML)
		os.Remove(tc.destParsedXML)

		t.Run(tc.desc, func(t *testing.T) {
			integrator := NewIntegrator([]string{})
			integrator.integrateFile(tc.src)

			gotParsedXML := readFileQuietly(tc.destParsedXML)
			wantParsedXML := readFileQuietly(tc.wantParsedXML)

			if diff := cmp.Diff(gotParsedXML, wantParsedXML); diff != "" {
				t.Errorf("failed: diff %s %s: (-got +want):\n%s\n", tc.destParsedXML, tc.wantParsedXML, diff)
			} else {
				os.Remove(tc.destParsedXML)
			}
			os.Remove(tc.destTokenizedXML)
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
