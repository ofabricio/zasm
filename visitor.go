package zasm

type Visitor[T any] interface {
	VisitProgram(*Program, T)
	VisitLabelDef(*LabelDef, T)
	VisitDataDef(*DataDef, T)
	VisitHexLiteral(*HexLiteral, T)
}

func VisitMany[T any](n []Node, v Visitor[T], out T) {
	for _, n := range n {
		Visit(n, v, out)
	}
}

func Visit[T any](n Node, v Visitor[T], out T) {
	switch n := n.(type) {
	case *Program:
		v.VisitProgram(n, out)
	case *LabelDef:
		v.VisitLabelDef(n, out)
	case *DataDef:
		v.VisitDataDef(n, out)
	case *HexLiteral:
		v.VisitHexLiteral(n, out)
	}
}
