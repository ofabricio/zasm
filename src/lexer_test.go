package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

func TestScanEmpty(t *testing.T) {

	src := ``

	tokens := Scan(strings.NewReader(src))

	Equal(t, len(tokens), 0)
}

func TestScan90(t *testing.T) {

	src := `90`

	tokens := Scan(strings.NewReader(src))

	Equal(t, len(tokens), 1)
	Equal(t, tokens[0].Text, "90")
	Equal(t, tokens[0].Type, "WORD")
}

func TestScan(t *testing.T) {

	// Given.

	outputGoldenFile := "test/lexer.output.golden"

	src, _ := ioutil.ReadFile("test/lexer.input.golden")
	exp, _ := ioutil.ReadFile(outputGoldenFile)

	exps := strings.Split(strings.TrimSpace(string(exp)), "\n")

	// When.

	toks := Scan(bytes.NewReader(src))

	// Then.

	if len(toks) != len(exps) {
		t.Fatalf("not enough tokens. Exp: %d, Got: %d", len(exps), len(toks))
	}

	all := ""
	got := ""
	for i, token := range toks {
		str := fmt.Sprintf("%s %s %d %d", token.Text, token.Type, token.Row, token.Col)
		if str != exps[i] {
			got += fmt.Sprintf("\nError at line %d: Exp: '%s', Got: '%s'", i+1, exps[i], str)
		}
		all += str + "\n"
	}
	if *update {
		ioutil.WriteFile(outputGoldenFile, []byte(all), 0644)
	}
	if got != "" {
		t.Fatalf(got)
	}
}

func Equal(t *testing.T, got, exp interface{}) {
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\nExp:\n%s\nGot:\n%s\n", exp, got)
	}
}
