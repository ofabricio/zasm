package zasm

import (
	"fmt"
	"io"
	"strings"
)

func PrintNode(n Node, out io.Writer) {
	Visit(n, &NodePrinter{out: out})
}

type NodePrinter struct {
	depth int
	out   io.Writer
}

func (p *NodePrinter) Print(s string) {
	fmt.Fprintf(p.out, "%s%v\n", strings.Repeat("    ", p.depth), s)
}

func (p *NodePrinter) VisitMany(n []Node) {
	p.depth++
	VisitMany(n, p)
	p.depth--
}

func (p *NodePrinter) VisitProgram(n *Program) {
	p.Print("Program [")
	p.VisitMany(n.Body)
	p.Print("]")
}

func (p *NodePrinter) VisitLabelDef(n *LabelDef) {
	p.Print(fmt.Sprintf("Label [%s]", n.Token.Text))
}

func (p *NodePrinter) VisitDataDef(n *DataDef) {
	p.Print("DataDef [")
	p.VisitMany(n.Body)
	p.Print("]")
}

func (p *NodePrinter) VisitHexLiteral(n *HexLiteral) {
	p.Print(fmt.Sprintf("Hex [%s]", n.Token.Text))
}

func (p *NodePrinter) VisitStringLiteral(n *StringLiteral) {
	p.Print(fmt.Sprintf("String [%s]", n.Token.Text))
}

func (p *NodePrinter) VisitDirDump(n *DirDump) {
	p.Print("@Dump [")
	p.VisitMany(n.Body)
	p.Print("]")
}
