package zasm

import (
	"os"
	"strings"
	"testing"
)

func TestAssemble(t *testing.T) {

	srcf, _ := os.ReadFile("testdata/source.asm")
	astf, _ := os.ReadFile("testdata/source.ast")
	dmpf, _ := os.ReadFile("testdata/source.dmp")

	var prg Program
	var oks = Parse(string(srcf), &prg)

	Resolve(&prg)

	// var bin []byte
	// Assemble(&prg, &bin)

	var ast strings.Builder
	PrintNode(&prg, &ast)

	var dmp strings.Builder
	Dump(&prg, &dmp)

	assert(t, true, oks, "")
	assert(t, strings.TrimSpace(string(astf)), strings.TrimSpace(ast.String()), "")
	assert(t, string(dmpf), dmp.String(), "")
}
