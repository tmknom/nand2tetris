package main

import (
	"bufio"
	"github.com/google/go-cmp/cmp"
	"os"
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
				"Square/Main.jack",
				"Square/Square.jack",
				"Square/SquareGame.jack",
			},
			destTokenizedXML: []string{
				"Square/MainT.xml",
				"Square/SquareT.xml",
				"Square/SquareGameT.xml",
			},
			destParsedXML: []string{
				"Square/Main.xml",
				"Square/Square.xml",
				"Square/SquareGame.xml",
			},
			wantTokenizedXML: []string{
				"Square/cmp/MainT.xml",
				"Square/cmp/SquareT.xml",
				"Square/cmp/SquareGameT.xml",
			},
			wantParsedXML: []string{
				"Square/cmp/Main.xml",
				"Square/cmp/Square.xml",
				"Square/cmp/SquareGame.xml",
			},
		},
		{
			desc: "ArrayTest",
			src: []string{
				"ArrayTest/Main.jack",
			},
			destTokenizedXML: []string{
				"ArrayTest/MainT.xml",
			},
			destParsedXML: []string{
				"ArrayTest/Main.xml",
			},
			wantTokenizedXML: []string{
				"ArrayTest/cmp/MainT.xml",
			},
			wantParsedXML: []string{
				"ArrayTest/cmp/Main.xml",
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
			src:              "Square/Square.jack",
			destTokenizedXML: "Square/SquareT.xml",
			destParsedXML:    "Square/Square.xml",
			wantTokenizedXML: "Square/cmp/SquareT.xml",
			wantParsedXML:    "Square/cmp/Square.xml",
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
			src:              "Square/Main.jack",
			destTokenizedXML: "Square/MainT.xml",
			destParsedXML:    "Square/Main.xml",
			wantParsedXML:    "Square/provisional/Main.xml",
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
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}
