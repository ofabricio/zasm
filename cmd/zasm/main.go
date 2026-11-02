package main

import (
	"os"

	"github.com/ofabricio/zasm"
)

func main() {
	src := os.Args[1]

	data, err := os.ReadFile(src)
	if err != nil {
		panic(err)
	}

	var prg zasm.Program
	var oks = zasm.Parse(string(data), &prg)

	if !oks {
		panic("failed to parse " + src)
	}

	zasm.Resolve(&prg)
	zasm.Dump(&prg, os.Stdout)
}
