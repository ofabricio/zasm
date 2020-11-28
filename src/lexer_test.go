package zasm

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

func TestScanEmpty(t *testing.T) {

	src := ``

	tokens := Scan(strings.NewReader(src))

	Equal(t, len(tokens), 0)
}

func TestScanA(t *testing.T) {

	src := "a"

	tokens := Scan(strings.NewReader(src))

	Equal(t, len(tokens), 1)
	Equal(t, tokens[0].Text, "a")
	Equal(t, tokens[0].Type, "WORD")
}

func TestScanSemicolon(t *testing.T) {

	src := ";"

	tokens := Scan(strings.NewReader(src))

	Equal(t, len(tokens), 0)
}

func TestScanQuote(t *testing.T) {

	src := "'"

	tokens := Scan(strings.NewReader(src))

	Equal(t, len(tokens), 1)
	Equal(t, tokens[0].Text, "'")
	Equal(t, tokens[0].Type, "INVALID")
}

func TestScanInvalidString(t *testing.T) {

	src := "'abc"

	tokens := Scan(strings.NewReader(src))

	Equal(t, len(tokens), 2)
	Equal(t, tokens[0].Text, "'")
	Equal(t, tokens[0].Type, "INVALID")
	Equal(t, tokens[1].Text, "abc")
	Equal(t, tokens[1].Type, "WORD")
}

func TestScanInvalidStringNL(t *testing.T) {

	src := "'abc\n"

	tokens := Scan(strings.NewReader(src))

	Equal(t, len(tokens), 3)
	Equal(t, tokens[0].Text, "'")
	Equal(t, tokens[0].Type, "INVALID")
	Equal(t, tokens[1].Text, "abc")
	Equal(t, tokens[1].Type, "WORD")
	Equal(t, tokens[2].Text, "\n")
	Equal(t, tokens[2].Type, "NL")
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
		txt := token.Text
		if token.Type == "NL" {
			txt = "\\n"
		}
		str := fmt.Sprintf("%s %s %d %d", txt, token.Type, token.Row, token.Col)
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
		_, fn, line, _ := runtime.Caller(1)
		t.Fatalf("\n[error] %s:%d\nExp:\n%v\nGot:\n%v\n", fn, line, exp, got)
	}
}
