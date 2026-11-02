package zasm

import (
	"regexp"
	"strings"
)

func Parse(src string, out *Program) bool {
	p := Parser{src: src}
	return p.Parse(out)
}

func (p *Parser) Parse(out *Program) bool {
	return p.Program(out)
}

func (p *Parser) Program(out *Program) bool {
	var stmts []Node
	var stmt Node
	for p.VS(true) && (p.Statement(&stmt) || p.Comment()) {
		if stmt != nil {
			stmts = append(stmts, stmt)
			stmt = nil
		}
	}
	*out = Program{Body: stmts}
	return true
}

func (p *Parser) Statement(out *Node) bool {
	return p.DataDef(out) || p.Label(out)
}

func (p *Parser) Label(out *Node) bool {
	var t Token
	if p.RegexOut(LBL, &t) {
		*out = &LabelDef{Token: t}
		return true
	}
	return false
}

func (p *Parser) DataDef(out *Node) bool {
	var v Node
	var vs []Node
	for p.HS(true) && p.HexLiteral(&v) {
		vs = append(vs, v)
	}
	if len(vs) > 0 {
		*out = &DataDef{Body: vs}
		return true
	}
	return false
}

func (p *Parser) HexLiteral(out *Node) bool {
	var t Token
	if p.RegexOut(HEX, &t) {
		*out = &HexLiteral{Token: t}
		return true
	}
	return false
}

func (p *Parser) RegexOut(r *regexp.Regexp, out *Token) bool {
	return p.Out(p.Mark(), p.Regex(r), out)
}

func (p *Parser) Comment() bool {
	return p.Regex(CMT)
}

func (p *Parser) HS(opt bool) bool {
	return p.Regex(HSP) || opt
}

func (p *Parser) VS(opt bool) bool {
	return p.Regex(VSP) || opt
}

func (p *Parser) Out(m Marker, cond bool, out *Token) bool {
	if cond {
		*out = p.Token(m)
	}
	return cond
}

func (p *Parser) Match(v string) bool {
	return strings.HasPrefix(p.Tail(), v) && p.Advance(v)
}

func (p *Parser) Regex(r *regexp.Regexp) bool {
	return p.Advance(r.FindString(p.Tail()))
}

func (p *Parser) Advance(v string) bool {
	p.Pos += len(v)
	return len(v) > 0
}

func (p *Parser) Tail() string {
	return p.src[p.Pos:]
}

func (p *Parser) Mark() Marker {
	return p.Marker
}

func (p *Parser) Token(m Marker) Token {
	return Token{Marker: m, Text: p.src[m.Pos:p.Pos]}
}

type Parser struct {
	Marker
	src string
}

type Token struct {
	Marker
	Text string
}

type Marker struct {
	Pos  int
	Line int
	Col  int
}

var HEX = regexp.MustCompile(`^[0-9A-F]+\b`)
var VSP = regexp.MustCompile(`^\s+`)
var HSP = regexp.MustCompile(`^[ \t]+`)
var CMT = regexp.MustCompile(`^;[^\n]*`)
var LBL = regexp.MustCompile(`^(\w+|>)`)
