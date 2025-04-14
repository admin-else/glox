package glox

//go:generate go run ../generateast/generateast.go
type Expr interface {
	Accept(visitor ExprVisitor) (any, error)
}

type BinaryExpr struct { 
	Left Expr
	Operator TokenType
	Right Expr
}

func (e BinaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinaryExpr(e)
}

type GroupingExpr struct { 
	Expr Expr
}

func (e GroupingExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGroupingExpr(e)
}

type LiteralExpr struct { 
	Value any
}

func (e LiteralExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteralExpr(e)
}

type UnaryExpr struct { 
	Operator TokenType
	Expr Expr
}

func (e UnaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnaryExpr(e)
}

type ExprVisitor interface { 
	VisitBinaryExpr(expr BinaryExpr) (any, error)
	VisitGroupingExpr(expr GroupingExpr) (any, error)
	VisitLiteralExpr(expr LiteralExpr) (any, error)
	VisitUnaryExpr(expr UnaryExpr) (any, error)
}
