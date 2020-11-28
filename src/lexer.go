package zasm

import (
	"io"

	"github.com/ofabricio/scanner"
)

// Scan is a handy function.
func Scan(r io.Reader) []*Token {
	return NewLexer(r).Scan()
}

// NewLexer creates a new lexer/scanner.
func NewLexer(r io.Reader) *Lexer {
	s := scanner.NewScanner(r)
	s.Space("^[ \t\r]+")
	return &Lexer{Scanner: s}
}

// Scan scans for tokens.
func (t *Lexer) Scan() (tokens []*Token) {
	for t.More() {
		typ := t.tokenize()
		if typ == "SPACE" || typ == "COMMENT" {
			continue
		}
		tok := &Token{Text: t.Text(), Type: typ, Row: t.Row(), Col: t.Col()}
		tokens = append(tokens, tok)
	}
	return
}

func (t *Lexer) tokenize() string {

	if t.Match("^\n") {
		return "NL"
	}

	if t.Match("^;.*") {
		return "COMMENT"
	}

	if t.String("'") {
		return "STRING"
	}

	if t.Match("^[\\w_]+") {
		return "WORD"
	}

	if t.Match(`^[()[\]{}|&?*+\-<>$,.@#=]`) {
		return "SYMBOL"
	}

	t.Match(".")
	return "INVALID"
}

// Lexer is a scanner.
type Lexer struct {
	*scanner.Scanner
}

// Token is a token.
type Token struct {
	Text string
	Type string
	Row  int
	Col  int
}
