package zasm

import (
	"fmt"
	"io"
	"strings"
)

func PrintNode(n Node, out io.Writer) {
	Visit(n, &NodePrinter{}, out)
}

type NodePrinter struct {
	depth int
}

func (p *NodePrinter) Print(s string, out io.Writer) {
	fmt.Fprintf(out, "%s%v\n", strings.Repeat("    ", p.depth), s)
}

func (p *NodePrinter) VisitMany(n []Node, out io.Writer) {
	p.depth++
	VisitMany(n, p, out)
	p.depth--
}

func (p *NodePrinter) VisitProgram(n *Program, out io.Writer) {
	p.Print("Program [", out)
	p.VisitMany(n.Body, out)
	p.Print("]", out)
}

func (p *NodePrinter) VisitLabelDef(n *LabelDef, out io.Writer) {
	p.Print(fmt.Sprintf("Label [%s]", n.Token.Text), out)
}

func (p *NodePrinter) VisitDataDef(n *DataDef, out io.Writer) {
	p.Print("DataDef [", out)
	p.VisitMany(n.Body, out)
	p.Print("]", out)
}

func (p *NodePrinter) VisitHexLiteral(n *HexLiteral, out io.Writer) {
	p.Print(fmt.Sprintf("Hex [%s]", n.Token.Text), out)
}
