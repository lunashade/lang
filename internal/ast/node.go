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
	Type     BinOpType
	LHS, RHS AST
}
type BinOpType int

const (
	Add BinOpType = iota + 1
	Sub
	Mul
	Div
)
