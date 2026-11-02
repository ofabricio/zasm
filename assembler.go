package zasm

import "strconv"

func Assemble(n Node, out *[]byte) {
	Visit(n, &Assembler{}, out)
}

type Assembler struct{}

func (a *Assembler) VisitProgram(n *Program, out *[]byte) {
	VisitMany(n.Body, a, out)
}

func (a *Assembler) VisitLabelDef(n *LabelDef, out *[]byte) {
}

func (a *Assembler) VisitDataDef(n *DataDef, out *[]byte) {
	VisitMany(n.Body, a, out)
}

func (a *Assembler) VisitHexLiteral(n *HexLiteral, out *[]byte) {
	for txt := n.Token.Text; len(txt) > 0; txt = txt[2:] {
		val, _ := strconv.ParseInt(txt[:2], 16, 64)
		*out = append(*out, byte(val))
	}
}
