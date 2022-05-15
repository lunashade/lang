package ast

type AST interface{}

type Root struct {
	Nodes []AST
}

type Int struct {
	Value int
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
