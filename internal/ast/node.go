package ast

type AST interface{}

type Int struct{ Value int }
type Add struct{ LHS, RHS AST }
type Sub struct{ LHS, RHS AST }
type Mul struct{ LHS, RHS AST }
type Div struct{ LHS, RHS AST }
