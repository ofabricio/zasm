package zasm

type Node any

type Program struct {
	Body []Node
}

type DataDef struct {
	Body []Node
}

type LabelDef struct {
	Token  Token
	Offset int
}

type HexLiteral struct {
	Token Token
}

type StringLiteral struct {
	Token Token
	Value string
}

type DirDump struct {
	Bytes int
	Body  []Node
}

type DirAST struct {
	Body []Node
}

func NodeIs[T any](n Node) bool {
	_, ok := n.(T)
	return ok
}
