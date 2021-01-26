package main

import (
	"testing"
)

func TestIntegratorIntegrate(t *testing.T) {
	cases := []struct {
		desc      string
		filenames []string
	}{
		{
			desc: "処理が最後まで完了すること（ファイル比較が可能になるまでの暫定）",
			filenames: []string{
				"Square/Main.jack",
				"Square/Square.jack",
				"Square/SquareGame.jack",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			integrator := NewIntegrator(tc.filenames)
			err := integrator.Integrate()
			if err != nil {
				t.Errorf("failed integrate: %+v", err)
			}
		})
	}
}
