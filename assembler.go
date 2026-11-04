package zasm

import (
	"bytes"
	"io"
	"strconv"
)

func Assemble(n Node, dmp, bin io.Writer) {
	Visit(n, &Resolver{})
	Visit(n, &Assembler{bin: bin, dmp: dmp})
}

type Assembler struct {
	bin io.Writer
	dmp io.Writer
	off int
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
		a.write([]byte{byte(val)})
	}
}

func (a *Assembler) VisitStringLiteral(n *StringLiteral) {
	a.write([]byte(n.Value))
}

func (a *Assembler) VisitDirDump(n *DirDump) {
	if a.dmp != nil {
		ad := AssemblerDumper{a: a, d: DumpWriter{out: a.dmp}}
		ad.VisitDirDump(n)
	}
}

func (a *Assembler) VisitDirAST(n *DirAST) {
	if a.dmp != nil {
		Visit(n, &NodePrinter{out: a.dmp})
	}
	tmp := a.dmp
	a.dmp = nil
	VisitMany(n.Body, a)
	a.dmp = tmp
}

func (a *Assembler) write(b []byte) {
	a.bin.Write(b)
	a.off += len(b)
}

type AssemblerDumper struct {
	a *Assembler
	d DumpWriter
}

func (d *AssemblerDumper) VisitDirDump(n *DirDump) {

	b := BytesCounter{}
	VisitMany(n.Body, &b)

	siz := d.d.siz
	d.d.SetOffset(d.a.off)
	d.d.SetBytes(n.Bytes)
	d.d.WriteNL()
	d.d.WriteTitle(b.count)
	for i, nn := range n.Body {
		if i == 0 && !NodeIs[*DirDump](nn) {
			d.d.WriteNL() // Add a NL when first node is not a @dump.
		}
		if i >= 1 && !NodeIs[*DirDump](nn) && NodeIs[*DirDump](n.Body[i-1]) {
			d.d.WriteNL() // Add a NL when exiting a @dump and the current node is not a @dump.
		}
		Visit(nn, d)
	}
	d.d.SetBytes(siz)
}

func (d *AssemblerDumper) VisitProgram(n *Program) {
	d.a.VisitProgram(n)
}

func (d *AssemblerDumper) VisitDataDef(n *DataDef) {
	var buf bytes.Buffer
	bin := d.a.bin
	d.a.bin = io.MultiWriter(bin, &buf)
	d.a.VisitDataDef(n)
	d.a.bin = bin
	d.d.WriteInstr(n, buf.Bytes())
}

func (d *AssemblerDumper) VisitHexLiteral(n *HexLiteral) {
	d.a.VisitHexLiteral(n)
}

func (d *AssemblerDumper) VisitStringLiteral(n *StringLiteral) {
	d.a.VisitStringLiteral(n)
}

func (d *AssemblerDumper) VisitLabelDef(n *LabelDef) {
	d.d.WriteLabel(n)
	d.a.VisitLabelDef(n)
}

func (d *AssemblerDumper) VisitDirAST(n *DirAST) {
	d.a.VisitDirAST(n)
}
