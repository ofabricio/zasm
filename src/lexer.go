package main

import (
	"io"
	"io/ioutil"
	"unicode"
	"unicode/utf8"
)

func Scan(r io.Reader) []*Token {
	lex := &Lexer{}
	return lex.Scan(r)
}

func (t *Lexer) Scan(r io.Reader) []*Token {

	data, _ := ioutil.ReadAll(r)

	var tokens []*Token

	row, col := 1, 1
	for len(data) > 0 {

		size, typ := t.tokenize(data)

		chars := utf8.RuneCount(data[:size])

		col += chars

		if typ == "NL" {
			row++
			col = 1
			typ = ""
		}

		if typ != "" {
			t := &Token{Text: string(data[:size]), Type: typ, Row: row, Col: col - chars}
			tokens = append(tokens, t)
		}

		data = data[size:]
	}

	return tokens
}

func (t *Lexer) tokenize(data []byte) (int, string) {
	r, size := utf8.DecodeRune(data)
	if r == '\n' {
		return size, "NL"
	}
	if unicode.IsSpace(r) {
		return size, ""
	}
	if r == ';' {
		return scan(data, Until('\n')), ""
	}
	if r == '\'' {
		return size + scan(data[size:], Until('\'')) + size, "STRING"
	}
	if IsAlphaNum_(r) {
		return scan(data, IsAlphaNum_), "WORD"
	}
	if IsTwinSymbol(r) {
		return scan(data, While(r)), "SYMBOL"
	}
	if IsSymbol(r) {
		return size, "SYMBOL"
	}
	return size, "UNKNOWN"
}

type Lexer struct {
}

type Token struct {
	Text string
	Type string
	Row  int
	Col  int
}

func scan(data []byte, cond MatchFunc) (adv int) {
	for {
		r, size := utf8.DecodeRune(data)
		if cond(r) {
			adv += size
		} else {
			break
		}
		data = data[size:]
	}
	return
}

func IsAlphaNum_(r rune) bool {
	return unicode.In(r, unicode.Letter, unicode.Number) || r == '_'
}

func IsSymbol(r rune) bool {
	return unicode.In(r, unicode.Symbol, unicode.Punct)
}

func IsTwinSymbol(r rune) bool {
	return r == '=' || r == '<' || r == '>' || r == '$'
}

func While(r rune) MatchFunc {
	return func(ru rune) bool {
		return ru == r
	}
}

func Until(r rune) MatchFunc {
	return func(ru rune) bool {
		return ru != r
	}
}

type MatchFunc func(rune) bool
