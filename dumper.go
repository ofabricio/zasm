package zasm

import (
	"fmt"
	"io"
)

type DumpWriter struct {
	out io.Writer
	off int
	siz int
}

func (d *DumpWriter) SetOffset(v int) {
	d.off = v
}

func (d *DumpWriter) SetBytes(v int) {
	d.siz = v
}

func (d *DumpWriter) WriteTitle(n int) {
	fmt.Fprintf(d.out, "    --- %d bytes ---\n", n)
}

func (d *DumpWriter) WriteLabel(n *LabelDef) {
	fmt.Fprintf(d.out, "    %s\n", n.Token.Text)
}

func (d *DumpWriter) WriteInstr(n Node, bin []byte) {
	pad := 4
	code := false
	for i, b := range bin {
		if (i)%d.siz == 0 {
			d.dumpOffset()
		}
		fmt.Fprintf(d.out, " %02X", b)
		if (i+1)%d.siz == 0 {
			if !code {
				d.dumpCode(n, pad)
				code = true
			}
			d.WriteNL()
		}
		d.off++
	}
	if !code {
		d.dumpCode(n, d.siz*3-len(bin)*3+pad)
	}
	if len(bin)%d.siz != 0 {
		d.WriteNL()
	}
}

func (d *DumpWriter) WriteNL() {
	fmt.Fprint(d.out, "\n")
}

func (d *DumpWriter) dumpOffset() {
	d.writeTabs()
	fmt.Fprintf(d.out, "%08X ", d.off)
}

func (d *DumpWriter) dumpCode(n Node, pad int) {
	d.writePad(pad)
	Unparse(n, d.out)
}

func (d *DumpWriter) writeTabs() {
	d.writePad(4)
}

func (d *DumpWriter) writePad(n int) {
	for range n {
		fmt.Fprint(d.out, " ")
	}
}

type BytesCounter struct {
	count int
}

func (c *BytesCounter) VisitProgram(n *Program) {
	VisitMany(n.Body, c)
}

func (c *BytesCounter) VisitLabelDef(n *LabelDef) {
}

func (c *BytesCounter) VisitDataDef(n *DataDef) {
	VisitMany(n.Body, c)
}

func (c *BytesCounter) VisitHexLiteral(n *HexLiteral) {
	c.count += len(n.Token.Text) / 2
}

func (c *BytesCounter) VisitStringLiteral(n *StringLiteral) {
	c.count += len(n.Value)
}

func (c *BytesCounter) VisitDirDump(n *DirDump) {
	VisitMany(n.Body, c)
}

func (c *BytesCounter) VisitDirAST(n *DirAST) {
	VisitMany(n.Body, c)
}
