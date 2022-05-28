package ast

type AST interface {
	node()
}

type Root struct {
	Nodes []AST
}

type Function struct {
	Name AST
	Body []AST
}

type Block struct {
	Stmts []AST
}

func (*Root) node()     {}
func (*Function) node() {}
func (*Block) node()    {}

// statements
type Stmt interface {
	AST
	stmtNode()
}

type ExprStmt struct {
	Expr AST
}
type Semi struct {
	Expr AST
}

func (*ExprStmt) node() {}
func (*Semi) node()     {}

func (*ExprStmt) stmtNode() {}
func (*Semi) stmtNode()     {}

// expressions
type Expr interface {
	AST
	exprNode()
}

type Int struct {
	Value int64
}

type Ident struct {
	Name string
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
	Assign
)

func (*Int) node()   {}
func (*Ident) node() {}
func (*BinOp) node() {}

func (*Int) exprNode()   {}
func (*Ident) exprNode() {}
func (*BinOp) exprNode() {}
