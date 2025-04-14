package glox

import (
	"fmt"
)

type RunTimeError struct {
	t Token
	m string
}

func (e RunTimeError) Error() string {
	return fmt.Sprintf("Error :%v around '%v': %v", e.t.Line, e.t.Lexme, e.m)
}

type interpreter struct {
	lol string
}

// VisitVarDecl implements StmtVisitor.
func (i *interpreter) VisitVarDecl(expr VarDecl) (any, error) {
	panic("unimplemented")
}

func (i *interpreter) VisitExprStmt(expr ExprStmt) (any, error) {
	_, err := i.evaluate(expr.Expr)
	return nil, err
}

func (i *interpreter) VisitPrintStmt(expr PrintStmt) (any, error) {
	v, err := i.evaluate(expr.Expr)
	fmt.Println(v)
	return nil, err
}

func (i *interpreter) VisitBinaryExpr(expr BinaryExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Type {
	case EQUAL_EQUAL:
		return left == right, nil
	case BANG_EQUAL:
		return left != right, nil
	}
	switch left.(type) {
	case float64:
		left, _ := left.(float64)
		right_num, ok := right.(float64)
		if ok == false {
			return nil, RunTimeError{m: "binary expr with a left num must have a right num", t: expr.Operator}
		}
		switch expr.Operator.Type {
		case PLUS:
			return left + right_num, nil
		case MINUS:
			return left - right_num, nil
		case SLASH:
			if right_num == 0 {
				return nil, RunTimeError{m: "cannot divide by 0", t: expr.Operator}
			}
			return left / right_num, nil
		case STAR:
			return left * right_num, nil
		case GREATER:
			return left > right_num, nil
		case GREATER_EQUAL:
			return left >= right_num, nil
		case LESS:
			return left < right_num, nil
		case LESS_EQUAL:
			return left <= right_num, nil
		default:
			panic("Unreachable")
		}
	case string:
		left, _ := left.(string)
		right, ok := right.(string)
		if !ok {
			return nil, RunTimeError{m: "binary epxr with a left string must have a right string", t: expr.Operator}
		}
		if expr.Operator.Type != PLUS {
			return nil, RunTimeError{m: "binary expr on strings only support +", t: expr.Operator}
		}
		return left + right, nil
	default:
		return nil, RunTimeError{m: "binary expr does not accept this type", t: expr.Operator}
	}
}

func (i *interpreter) VisitGroupingExpr(expr GroupingExpr) (any, error) {
	return i.evaluate(expr.Expr)
}

func (i *interpreter) VisitLiteralExpr(expr LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *interpreter) VisitUnaryExpr(expr UnaryExpr) (any, error) {
	right, err := i.evaluate(expr.Expr)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case BANG:
		return i.isTruthy(right), nil
	case MINUS:
		v, ok := right.(float64)
		if !ok {
			return nil, RunTimeError{m: "negation can only be done on numbers", t: expr.Operator}
		}
		return -v, nil
	default:
		panic("Unreachable")
	}
}

func (i *interpreter) isTruthy(value any) bool {
	switch value.(type) {
	case nil:
		return false
	case bool:
		return value.(bool)
	default:
		return true
	}
}

func (i *interpreter) evaluate(expr Expr) (any, error) {
	return expr.Accept(i)
}

func (i *interpreter) execute(stmt Stmt) (any, error) {
	return stmt.Accept(i)
}

func Interpret(e []Stmt) error {
	i := &interpreter{}
	for _, v := range e {
		_, err := i.execute(v)
		if err != nil {
			return err
		}
	}
	return nil
}
