package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ofabricio/zasm"
)

func main() {

	var src string
	var out string

	flag.StringVar(&src, "src", "", "the source code file")
	flag.StringVar(&out, "out", "", "the binary output file")
	flag.Parse()

	data, err := os.ReadFile(src)
	if err != nil {
		panic(fmt.Errorf("source file: %w", err))
	}

	var prg zasm.Program

	if err := zasm.Parse(string(data), &prg); err != nil {
		panic("failed to parse source code")
	}

	tgt := io.Discard

	if out != "" {
		f, err := os.Create(out)
		if err != nil {
			panic(fmt.Errorf("failed to create output file: %w", err))
		}
		defer f.Close()
		tgt = f
	}

	zasm.Assemble(&prg, os.Stdout, tgt)
}
