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

type DirDump struct {
	Bytes int
	Body  []Node
}
