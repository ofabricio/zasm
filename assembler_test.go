package zasm

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestAll(t *testing.T) {

	tt := []string{
		"testdata/@ast.test",
		"testdata/@dump.test",
	}

	for _, tc := range tt {
		t.Run(tc, func(t *testing.T) {

			srcf, dmpf, binf := SplitFile(tc)

			var prg Program
			var oks = Parse(srcf, &prg)

			var bin bytes.Buffer
			var dmp bytes.Buffer
			Assemble(&prg, &dmp, &bin)

			assert(t, true, oks, tc)
			assert(t, string(dmpf), dmp.String(), tc)
			assert(t, string(binf), fmt.Sprintf("%X", bin.Bytes()), tc)
		})
	}
}

func SplitFile(file string) (string, string, string) {
	tst, _ := os.ReadFile(file)
	spt := strings.Split(string(tst), "---------- SPLIT ----------\n")
	for range 3 - len(spt) {
		spt = append(spt, "")
	}
	return spt[0], spt[1], strings.TrimSpace(spt[2])
}

func assert(t *testing.T, exp, got any, msg string) {
	t.Helper()
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("\nMsg:\n%s\nGot:\n%v\nExp:\n%v\n", msg, got, exp)
	}
}
