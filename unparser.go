package zasm

import (
	"fmt"
	"io"
)

// Unparse converts the AST back to the code.
func Unparse(n Node, out io.Writer) {
	Visit(n, &Unparser{out})
}

// Unparse converts the AST back to the code.
type Unparser struct {
	out io.Writer
}

func (u *Unparser) write(v string) {
	fmt.Fprint(u.out, v)
}

func (u *Unparser) VisitProgram(n *Program) {
	VisitMany(n.Body, u)
}

func (u *Unparser) VisitDataDef(n *DataDef) {
	for i, n := range n.Body {
		if i > 0 {
			u.write(" ")
		}
		Visit(n, u)
	}
}

func (u *Unparser) VisitDirDump(n *DirDump) {
	VisitMany(n.Body, u)
}

func (u *Unparser) VisitHexLiteral(n *HexLiteral) {
	u.write(n.Token.Text)
}

func (u *Unparser) VisitLabelDef(n *LabelDef) {
	u.write(n.Token.Text)
}
