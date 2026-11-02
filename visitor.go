package zasm

type Visitor interface {
	VisitProgram(*Program)
	VisitLabelDef(*LabelDef)
	VisitDataDef(*DataDef)
	VisitHexLiteral(*HexLiteral)
}

func Visit(n Node, v Visitor) {
	switch n := n.(type) {
	case *Program:
		v.VisitProgram(n)
	case *LabelDef:
		v.VisitLabelDef(n)
	case *DataDef:
		v.VisitDataDef(n)
	case *HexLiteral:
		v.VisitHexLiteral(n)
	}
}

func VisitMany(n []Node, v Visitor) {
	for _, n := range n {
		Visit(n, v)
	}
}
