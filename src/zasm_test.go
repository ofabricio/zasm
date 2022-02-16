package zasm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	. "github.com/ofabricio/calm"
	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {

	in, _ := os.ReadFile("test/input.golden")
	ou, _ := os.ReadFile("test/output.parser.golden")

	p := &Parser{c: New(string(in))}
	ast, ok := p.Parse()

	assert.True(t, ok)
	assert.Equal(t, PrintParserOutput(ast)+"\n", string(ou))

	g := Generate(ast)

	assert.Equal(t, "89C0B80700000090B808000000B802000000B81C000000B81E00000089C0B802000000", fmt.Sprintf("%X", g))
}

func PrintParserOutput(a *Ast) string {
	if a == nil {
		return ""
	}
	var s []string
	for _, v := range a.Args {
		s = append(s, PrintParserOutput(v))
	}
	switch a.Type {
	case "program":
		return strings.Join(s, "\n")
	}
	if s == nil {
		return fmt.Sprintf("{ %s %s | O %d | B %d | S %d }",
			a.Type, a.Name.Text, a.Offs, a.Bits, a.Size)
	}
	return fmt.Sprintf("{ %s %s | O %d | B %d | S %d | A [%s] }",
		a.Type, a.Name.Text, a.Offs, a.Bits, a.Size, strings.Join(s, ", "))
}

func (p *Parser) Parse() (*Ast, bool) {
	ast := &Ast{Name: Token{Text: "Program"}, Type: "program"}
	if !p.Program(ast).Run(p.c) {
		return nil, false
	}
	Walk(&offsetVisitor{}, ast)
	Walk(&labelsVisitor{}, ast)
	return ast, true
}

func (p *Parser) Program(a *Ast) MatcherFunc {
	err := Until(Eq("\n")).On(func(t Token) {
		fmt.Printf("%s\n^ Unknown token in L%d C%d\n", t.Text, t.Row, t.Col)
	})
	return And(p.statement(a).ZeroToMany(), err.Not())
}

func (p *Parser) statement(a *Ast) MatcherFunc {
	var stmt Ast
	return Or(
		wz.False(),
		p.comment(),
		p.directive(&stmt).On(AddArgs(a, &stmt)),
		p.instMov(&stmt).On(AddArgs(a, &stmt)),
		p.datadef(&stmt).On(AddArgs(a, &stmt)),
		p.labeldef(&stmt).On(AddArgs(a, &stmt)),
	)
}

func (p *Parser) directive(a *Ast) MatcherFunc {
	return Or(
		p.directiveBits(a),
		p.directiveAlign(a),
		p.directivePrint(a),
	)
}

func (p *Parser) directiveBits(a *Ast) MatcherFunc {
	return S("@64").On(EmitNode(a, "dir"))
}

func (p *Parser) directiveAlign(a *Ast) MatcherFunc {
	var arg Ast
	return And(S("@align").On(EmitNode(a, "dir")), ws, S("4").On(EmitNode(&arg, "cst"))).On(AddArgs(a, &arg))
}

func (p *Parser) directivePrint(a *Ast) MatcherFunc {
	body := func(c *Code) bool {
		return p.statement(a).ZeroToMany().Run(c)
	}
	return And(S("@print").On(EmitNode(a, "dir")), wz, S("{"), body, wz, S("}"))
}

func (p *Parser) comment() MatcherFunc {
	return And(S(";"), Until(S("\n")))
}

func (p *Parser) datadef(n *Ast) MatcherFunc {
	var a Ast
	return And(S("db").On(EmitNode(n, "def")), ws, F(unicode.IsDigit).OneToMany().On(EmitNode(&a, "cst"))).On(AddArgs(n, &a))
}

func (p *Parser) instMov(n *Ast) MatcherFunc {
	var a, b Ast
	args := Or(
		wz.False(),
		And(p.reg32(&a), ws, p.reg32(&b)).Undo(),
		And(p.reg32(&a), ws, p.label(&b)).Undo(),
	).On(AddArgs(n, &a, &b))
	return And(S("mov").On(EmitNode(n, "ins")), args)
}

func (p *Parser) label(a *Ast) MatcherFunc {
	return Or(
		F(unicode.IsLetter).OneToMany().On(EmitNode(a, "lab")),
		S(">").OneToMany().On(EmitNode(a, "ref")),
		S("<").OneToMany().On(EmitNode(a, "ref")),
	)
}

func (p *Parser) labeldef(a *Ast) MatcherFunc {
	return Or(
		F(unicode.IsLetter).OneToMany().On(EmitNode(a, "labdef")),
		S(">").On(EmitNode(a, "refdef")),
	)
}

func (p *Parser) reg32(a *Ast) MatcherFunc {
	return Or(
		S("eax"),
	).On(EmitNode(a, "r32"))
}

func AddArgs(n *Ast, args ...*Ast) func(Token) {
	return func(Token) {
		for _, a := range args {
			aa := *a
			n.Args = append(n.Args, &aa)
		}
	}
}

func EmitNode(a *Ast, typ string) func(Token) {
	return func(t Token) {
		*a = Ast{Name: t, Type: typ}
	}
}

var ws = Or(S(" "), S("\t")).OneToMany()
var wz = F(unicode.IsSpace).ZeroToMany()

type Parser struct {
	c *Code
}

func genMov(i *Ast) []byte {
	op1 := i.Args[0].Name.Text
	op2 := i.Args[1].Name.Text
	// Opcode | Instruction    | Op/En | 64-Bit Mode | Compat/Leg Mode | Description
	// 89 /r  | MOV r/m32,r32  | MR    | Valid       | Valid           | Move r32 to r/m32
	return opcode_r32_r32(0x89, op1, op2)
}

func opcode_r32_r32(opcode byte, op1, op2 string) []byte {
	mod := registers[op1].Mod << 6
	reg := registers[op2].REG << 3
	rm := registers[op1].RM
	return []byte{opcode, mod + reg + rm}
}

var registers = map[string]RegInfo{
	"eax": {Mod: 3, REG: 0, RM: 0, Bits: 32},
}

type RegInfo struct {
	Mod  byte
	REG  byte
	RM   byte
	Bits byte
}

// registers = {
//     'eax': {'Mod': 3, 'REG': 0, 'RM': 0, 'bits': 32},
//     'ecx': {'Mod': 3, 'REG': 1, 'RM': 1, 'bits': 32},
//     'edx': {'Mod': 3, 'REG': 2, 'RM': 2, 'bits': 32},
//     'ebx': {'Mod': 3, 'REG': 3, 'RM': 3, 'bits': 32},
//     'esp': {'Mod': 3, 'REG': 4, 'RM': 4, 'bits': 32},
//     'ebp': {'Mod': 3, 'REG': 5, 'RM': 5, 'bits': 32},
//     'esi': {'Mod': 3, 'REG': 6, 'RM': 6, 'bits': 32},
//     'edi': {'Mod': 3, 'REG': 7, 'RM': 7, 'bits': 32},
// }

type Ast struct {
	Type string
	Name Token
	Bits int
	Size int
	Offs int
	Args []*Ast
}

func Visit(n *Ast, fn func(*Ast) bool) {
	if n == nil {
		return
	}
	if fn(n) {
		for _, v := range n.Args {
			Visit(v, fn)
		}
	}
}

func Generate(ast *Ast) []byte {
	gen := &codegenVisitor{}
	Walk(gen, ast)
	return gen.bytecode
}

type labelsVisitor struct {
	labdefs  map[string]int
	refdefs  map[int]int
	usage    []*Ast
	refCount int
}

func (t *labelsVisitor) Visit(n *Ast) Visitor {

	if t.labdefs == nil {
		t.labdefs = map[string]int{}
		t.refdefs = map[int]int{}
	}

	if n == nil {
		return nil
	}

	switch n.Type {
	case "program":
		WalkAll(t, n.Args)

		for _, u := range t.usage {
			if u.Type == "ref" {
				// TODO: check if ref is not defined and
				// return an error "ref label not found".
				u.Offs = t.refdefs[u.Offs]
			}
			if u.Type == "lab" {
				// TODO: check if u.Name.Text is not defined
				// and return an error "label not found".
				u.Offs = t.labdefs[u.Name.Text]
			}
		}
		return nil
	case "labdef":
		t.labdefs[n.Name.Text] = n.Offs
	case "refdef":
		t.refdefs[len(t.refdefs)] = n.Offs
		t.refCount = len(t.refdefs)
	case "lab":
		t.usage = append(t.usage, n)
	case "ref":
		len := utf8.RuneCountInString(n.Name.Text)
		n.Offs = t.refCount - (len + 0) // backward
		if strings.HasPrefix(n.Name.Text, ">") {
			n.Offs = t.refCount + (len - 1) // forward
		}
		// n.Offs above receives the refCount that will be used
		// later to fill in n.Offs with the actual offset.
		t.usage = append(t.usage, n)
	}
	return t
}

type offsetVisitor struct {
	offset int
}

func (t *offsetVisitor) Visit(n *Ast) Visitor {
	if n == nil {
		return nil
	}
	switch n.Type {
	case "ins":
		if n.Name.Text == "mov" {
			if n.Args[0].Type == "r32" && n.Args[1].Type == "r32" {
				n.Size = 2
			}
			if n.Args[0].Type == "r32" && n.Args[1].Type == "lab" {
				n.Size = 5
			}
			if n.Args[0].Type == "r32" && n.Args[1].Type == "ref" {
				n.Size = 5
			}
			n.Offs = t.offset
		}
	case "dir":
		if n.Name.Text == "@align" {
			v, _ := strconv.Atoi(n.Args[0].Name.Text)
			n.Size = v - (t.offset % v)
			if n.Size == v {
				n.Size = 0
			}
			n.Offs = t.offset
		}
		if n.Name.Text == "@print" {
			n.Offs = t.offset
			for _, a := range n.Args {
				Walk(t, a)
			}
			n.Size = t.offset - n.Offs
			return nil
		}
	case "def":
		if n.Name.Text == "db" {
			n.Size = len(n.Args)
			n.Offs = t.offset
		}
	case "r32", "lab", "ref", "cst":
		n.Offs = 0
	default:
		n.Offs = t.offset
	}
	t.offset += n.Size
	return t
}

type Visitor interface {
	Visit(*Ast) Visitor
}

func Walk(v Visitor, node *Ast) {
	if v = v.Visit(node); v == nil {
		return
	}
	WalkAll(v, node.Args)
	v.Visit(nil)
}

func WalkAll(v Visitor, nodes []*Ast) {
	for _, n := range nodes {
		Walk(v, n)
	}
}

type codegenVisitor struct {
	bytecode []byte
}

func (t *codegenVisitor) Visit(n *Ast) Visitor {
	if n == nil {
		return nil
	}
	switch n.Type {
	case "ins":
		if n.Args[0].Type == "r32" && n.Args[1].Type == "r32" {
			t.bytecode = append(t.bytecode, genMov(n)...)
		}
		if n.Args[0].Type == "r32" && (n.Args[1].Type == "lab" || n.Args[1].Type == "ref") {
			b := make([]byte, 4)
			binary.LittleEndian.PutUint32(b, uint32(n.Args[1].Offs))
			t.bytecode = append(t.bytecode, 0xB8)
			t.bytecode = append(t.bytecode, b...)
		}
	case "dir":
		if n.Name.Text == "@align" {
			t.bytecode = append(t.bytecode, bytes.Repeat([]byte{0x90}, n.Size)...)
		}
	case "def":
		if n.Name.Text == "db" {
			v, _ := strconv.ParseInt(n.Args[0].Name.Text, 16, 0)
			t.bytecode = append(t.bytecode, bytes.Repeat([]byte{byte(v)}, n.Size)...)
		}
	}
	return t
}

type printDirectiveVisitor struct {
	buf io.Writer
	dep int
}

func (v *printDirectiveVisitor) Visit(n *Ast) Visitor {
	switch n.Type {
	case "program":
		v.nl()
		for _, a := range n.Args {
			if v.isFatPrint(a) {
				v.nl()
				Walk(v, a)
			}
		}
		v.nl()
		return nil
	case "dir":
		if n.Name.Text == "@print" && len(n.Args) != 0 {
			v.dep++
			v.pad()
			v.print("[%d bytes]", n.Size)
			v.nl()
			v.pad()
			v.nl()
			var prev *Ast
			for i, a := range n.Args {
				if i > 0 && (v.isFatPrint(a) || v.isFatPrint(prev)) {
					v.pad()
					v.nl()
				}
				prev = a
				Walk(v, a)
			}
			v.dep--
			return nil
		}
	case "ins", "def":
		v.pad()
		v.print("% X", Generate(n))
		v.print("    %s", n.Name.Text)
		for _, a := range n.Args {
			v.print(" %s", a.Name.Text)
		}
		v.nl()
		return nil
	}
	return nil
}

// Fat Print means that it is a print directive
// with at least one instruction in it's body.
func (v *printDirectiveVisitor) isFatPrint(a *Ast) bool {
	if a.Type == "dir" && a.Name.Text == "@print" {
		for _, arg := range a.Args {
			if arg.Type != "dir" || arg.Name.Text != "@print" {
				return true
			}
		}
	}
	return false
}

func (v *printDirectiveVisitor) print(format string, token ...interface{}) {
	fmt.Fprintf(v.buf, format, token...)
}

func (v *printDirectiveVisitor) pad() {
	pad := "    "
	if v.dep > 0 {
		pad = "|   "
	}
	dep := v.dep - 1
	if dep < 0 {
		dep = 0
	}
	fmt.Fprint(v.buf, "    "+strings.Repeat(pad, dep))
}

func (v *printDirectiveVisitor) nl() {
	fmt.Fprint(v.buf, "\n")
}
