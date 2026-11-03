package zasm

type Resolver struct {
	offset int
}

func (r *Resolver) VisitProgram(n *Program) {
	VisitMany(n.Body, r)
}

func (r *Resolver) VisitLabelDef(n *LabelDef) {
	n.Offset = r.offset
}

func (r *Resolver) VisitDataDef(n *DataDef) {
	VisitMany(n.Body, r)
}

func (r *Resolver) VisitHexLiteral(n *HexLiteral) {
	r.offset += len(n.Token.Text) / 2
}

func (r *Resolver) VisitDirDump(n *DirDump) {
	VisitMany(n.Body, r)
}
