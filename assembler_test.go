package zasm

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestAssemble(t *testing.T) {

	srcf, _ := os.ReadFile("testdata/source.asm")
	astf, _ := os.ReadFile("testdata/source.ast")
	dmpf, _ := os.ReadFile("testdata/source.dmp")
	binf, _ := os.ReadFile("testdata/source.bin")

	var prg Program
	var oks = Parse(string(srcf), &prg)

	var bin bytes.Buffer
	var dmp bytes.Buffer
	var ast bytes.Buffer

	Assemble(&prg, &dmp, &bin)

	PrintNode(&prg, &ast)

	assert(t, true, oks, "")
	assert(t, string(astf), ast.String(), "")
	assert(t, string(dmpf), dmp.String(), "")
	assert(t, string(binf), fmt.Sprintf("%X", bin.Bytes()), "")
}
