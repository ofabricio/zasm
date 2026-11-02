package zasm

func Resolve(n *Program) {
	Visit(n, &Resolver{}, nil)
}

type Resolver struct {
	offset int
}

func (r *Resolver) VisitProgram(n *Program, out *int) {
	VisitMany(n.Body, r, out)
}

func (r *Resolver) VisitLabelDef(n *LabelDef, out *int) {
	n.Offset = r.offset
}

func (r *Resolver) VisitDataDef(n *DataDef, out *int) {
	VisitMany(n.Body, r, out)
}

func (r *Resolver) VisitHexLiteral(n *HexLiteral, out *int) {
	r.offset += len(n.Token.Text) / 2
}
