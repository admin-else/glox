package glox

//go:generate go run ../generateast/generateast.go

type Expr interface {
	Accept(visitor ExprVisitor) (any, error)
} 

type BinaryExpr struct { 
	Left Expr
	Operator Token
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
	Operator Token
	Expr Expr
}

func (e UnaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnaryExpr(e)
}

type VariableExpr struct { 
	Name Token
}

func (e VariableExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVariableExpr(e)
}

type AssignExpr struct { 
	Name Token
	Value Expr
}

func (e AssignExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitAssignExpr(e)
}

type LogicalExpr struct { 
	Left Expr
	Operator Token
	Right Expr
}

func (e LogicalExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLogicalExpr(e)
}

type ExprVisitor interface { 
	VisitBinaryExpr(expr BinaryExpr) (any, error)
	VisitGroupingExpr(expr GroupingExpr) (any, error)
	VisitLiteralExpr(expr LiteralExpr) (any, error)
	VisitUnaryExpr(expr UnaryExpr) (any, error)
	VisitVariableExpr(expr VariableExpr) (any, error)
	VisitAssignExpr(expr AssignExpr) (any, error)
	VisitLogicalExpr(expr LogicalExpr) (any, error)
}

type Stmt interface {
	Accept(visitor StmtVisitor) (any, error)
} 

type ExprStmt struct { 
	Expr Expr
}

func (e ExprStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitExprStmt(e)
}

type PrintStmt struct { 
	Expr Expr
}

func (e PrintStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitPrintStmt(e)
}

type VarDecl struct { 
	Name Token
	Initializer Expr
}

func (e VarDecl) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitVarDecl(e)
}

type Block struct { 
	Stmts []Stmt
}

func (e Block) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitBlock(e)
}

type IfStmt struct { 
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (e IfStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitIfStmt(e)
}

type StmtVisitor interface { 
	VisitExprStmt(expr ExprStmt) (any, error)
	VisitPrintStmt(expr PrintStmt) (any, error)
	VisitVarDecl(expr VarDecl) (any, error)
	VisitBlock(expr Block) (any, error)
	VisitIfStmt(expr IfStmt) (any, error)
}

