package zasm

import (
	"bytes"
	"fmt"
	"io"
)

func Dump(n *Program, out io.Writer) {
	c := &CodeWidth{}
	Visit(n, c)
	Visit(n, &CodeDumper{width: c.width, out: out})
}

type CodeDumper struct {
	width int
	out   io.Writer
	pos   int
}

func (d *CodeDumper) VisitProgram(n *Program) {
	d.width += 5
	VisitMany(n.Body, d)
}

func (d *CodeDumper) VisitLabelDef(n *LabelDef) {
	d.writeCode(" ")
	d.writeCode(n.Token.Text)
	d.writePad()
	fmt.Fprintf(d.out, " OFFSET %02X", n.Offset)
	d.writeNL()
}

func (d *CodeDumper) VisitDataDef(n *DataDef) {
	d.writeCode("   ")
	for _, n := range n.Body {
		d.writeCode(" ")
		Visit(n, d)
	}
	d.writePad()
	var b bytes.Buffer
	Assemble(n, &b)
	for _, b := range b.Bytes() {
		fmt.Fprintf(d.out, " %02X", b)
	}
	d.writeNL()
}

func (d *CodeDumper) VisitHexLiteral(n *HexLiteral) {
	d.writeCode(n.Token.Text)
}

func (d *CodeDumper) writeCode(v string) {
	fmt.Fprint(d.out, v)
	d.pos += len(v)
}

func (d *CodeDumper) writePad() {
	for range d.width - d.pos {
		fmt.Fprint(d.out, " ")
	}
	fmt.Fprint(d.out, ";")
}

func (d *CodeDumper) writeNL() {
	d.pos = 0
	fmt.Fprint(d.out, "\n")
}

type CodeWidth struct {
	width int
}

func (d *CodeWidth) VisitProgram(n *Program) {
	VisitMany(n.Body, d)
}

func (d *CodeWidth) VisitLabelDef(n *LabelDef) {
	d.width = max(d.width, len(n.Token.Text))
}

func (d *CodeWidth) VisitDataDef(n *DataDef) {
	sum := len(n.Body) - 1
	for _, n := range n.Body {
		w := CodeWidth{}
		Visit(n, &w)
		sum += w.width
	}
	d.width = max(d.width, sum)
}

func (d *CodeWidth) VisitHexLiteral(n *HexLiteral) {
	d.width = len(n.Token.Text)
}
