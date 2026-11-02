package zasm

import (
	"io"
	"strconv"
)

func Assemble(n Node, out io.Writer) {
	Visit(n, &Resolver{})
	Visit(n, &Assembler{out})
}

type Assembler struct {
	out io.Writer
}

func (a *Assembler) VisitProgram(n *Program) {
	VisitMany(n.Body, a)
}

func (a *Assembler) VisitLabelDef(n *LabelDef) {
}

func (a *Assembler) VisitDataDef(n *DataDef) {
	VisitMany(n.Body, a)
}

func (a *Assembler) VisitHexLiteral(n *HexLiteral) {
	for txt := n.Token.Text; len(txt) > 0; txt = txt[2:] {
		val, _ := strconv.ParseInt(txt[:2], 16, 64)
		a.out.Write([]byte{byte(val)})
	}
}
