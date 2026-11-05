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
	Visit(n, &DumperDecorator{a: a, d: DumpWriter{out: a.dmp}})
}

func (a *Assembler) VisitDirAST(n *DirAST) {
	Visit(n, &NodePrinter{out: a.dmp})
	VisitMany(n.Body, a)
}

func (a *Assembler) write(b []byte) {
	a.bin.Write(b)
	a.off += len(b)
}

type DumperDecorator struct {
	a *Assembler
	d DumpWriter
}

func (d *DumperDecorator) VisitDirDump(n *DirDump) {

	b := BytesCounter{}
	VisitMany(n.Body, &b)

	siz := d.d.siz
	d.d.SetOffset(d.a.off)
	d.d.SetBytes(n.Bytes)
	d.d.WriteNL()
	d.d.WriteTitle(b.count)
	d.d.WriteNL()
	VisitMany(n.Body, d)
	d.d.SetBytes(siz)
}

func (d *DumperDecorator) VisitProgram(n *Program) {
	d.a.VisitProgram(n)
}

func (d *DumperDecorator) VisitDataDef(n *DataDef) {
	var buf bytes.Buffer
	bin := d.a.bin
	d.a.bin = io.MultiWriter(bin, &buf)
	d.a.VisitDataDef(n)
	d.a.bin = bin
	d.d.WriteInstr(n, buf.Bytes())
}

func (d *DumperDecorator) VisitHexLiteral(n *HexLiteral) {
	d.a.VisitHexLiteral(n)
}

func (d *DumperDecorator) VisitStringLiteral(n *StringLiteral) {
	d.a.VisitStringLiteral(n)
}

func (d *DumperDecorator) VisitLabelDef(n *LabelDef) {
	d.d.WriteLabel(n)
	d.a.VisitLabelDef(n)
}

func (d *DumperDecorator) VisitDirAST(n *DirAST) {
	d.a.VisitDirAST(n)
}
