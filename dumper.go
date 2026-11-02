package zasm

import (
	"fmt"
	"io"
	"strings"
)

func Dump(n *Program, out io.Writer) {
	var width int
	Visit(n, &CodeWidth{}, &width)
	Visit(n, &CodeDumper{width}, out)
}

type CodeDumper struct {
	width int
}

func (d *CodeDumper) VisitProgram(n *Program, out io.Writer) {
	VisitMany(n.Body, d, out)
}

func (d *CodeDumper) VisitLabelDef(n *LabelDef, out io.Writer) {
	d.write(out, " "+n.Token.Text, fmt.Sprintf(" OFFSET 0x%02X", n.Offset))
}

func (d *CodeDumper) VisitDataDef(n *DataDef, out io.Writer) {
	var code strings.Builder
	code.WriteString("   ")
	for _, n := range n.Body {
		code.WriteString(" ")
		Visit[io.Writer](n, d, &code)
	}
	var b []byte
	Assemble(n, &b)
	var cmmt strings.Builder
	for _, b := range b {
		fmt.Fprintf(&cmmt, " %02X", b)
	}
	d.write(out, code.String(), cmmt.String())
}

func (d *CodeDumper) VisitHexLiteral(n *HexLiteral, out io.Writer) {
	fmt.Fprintf(out, "%s", n.Token.Text)
}

func (d *CodeDumper) write(out io.Writer, code, cmmt string) {
	pad := strings.Repeat(" ", d.width+5-len(code))
	fmt.Fprintf(out, "%s%s;%s\n", code, pad, cmmt)
}

type CodeWidth struct{}

func (d *CodeWidth) VisitProgram(n *Program, out *int) {
	VisitMany(n.Body, d, out)
}

func (d *CodeWidth) VisitLabelDef(n *LabelDef, out *int) {
	*out = max(*out, len(n.Token.Text))
}

func (d *CodeWidth) VisitDataDef(n *DataDef, out *int) {
	sum := len(n.Body) - 1
	for _, n := range n.Body {
		w := 0
		Visit(n, d, &w)
		sum += w
	}
	*out = max(*out, sum)
}

func (d *CodeWidth) VisitHexLiteral(n *HexLiteral, out *int) {
	*out = len(n.Token.Text)
}
