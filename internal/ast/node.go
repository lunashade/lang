package ast

type AST interface{}

type Root struct {
	Nodes []AST
}

type Expr struct {
	Node AST
}

type Int struct {
	Value int64
}

type BinOp struct {
	Kind     BinOpKind
	LHS, RHS AST
}
type BinOpKind int

const (
	Add BinOpKind = iota + 1
	Sub
	Mul
	Div
)
