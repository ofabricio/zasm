package zasm

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test(t *testing.T) {

	src, _ := os.ReadFile("testdata/source.asm")
	ast, _ := os.ReadFile("testdata/source.ast")

	var out Program
	var oks = Parse(string(src), &out)

	var buf strings.Builder
	PrintNode(&out, &buf)

	assert(t, true, oks, "")
	assert(t, strings.TrimSpace(string(ast)), strings.TrimSpace(buf.String()), "")
}

func assert(t *testing.T, exp, got any, msg string) {
	t.Helper()
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("\nMsg:\n%s\nGot:\n%v\nExp:\n%v\n", msg, got, exp)
	}
}
