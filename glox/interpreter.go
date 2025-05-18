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

type Interpreter struct {
	*Enviorment
}

func (i *Interpreter) VisitLogicalExpr(expr LogicalExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}
	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitIfStmt(expr IfStmt) (_ any, err error) {
	ret, err := i.evaluate(expr.Condition)
	if err != nil {
		return
	}
	if i.isTruthy(ret) {
		_, err = i.execute(expr.ThenBranch)
	} else {
		_, err = i.execute(expr.ElseBranch)
	}
	return
}

func (i *Interpreter) VisitBlock(expr Block) (_ any, err error) {
	err = i.executeBlock(expr.Stmts, &Enviorment{i.Enviorment, map[string]any{}})
	return
}

func (i *Interpreter) executeBlock(stmts []Stmt, env *Enviorment) (err error) {
	prev := i.Enviorment
	i.Enviorment = env
	for _, stmt := range stmts {
		_, err = i.execute(stmt)
		if err != nil {
			return err
		}
	}
	i.Enviorment = prev
	return
}

func (i *Interpreter) VisitAssignExpr(expr AssignExpr) (value any, err error) {
	value, err = i.evaluate(expr.Value)
	if err != nil {
		return
	}
	err = i.Assign(expr.Name, value)
	return
}

func (i *Interpreter) VisitVariableExpr(expr VariableExpr) (any, error) {
	return i.Get(expr.Name)
}

func (i *Interpreter) VisitVarDecl(expr VarDecl) (_ any, err error) {
	var value any
	if expr.Initializer != nil {
		value, err = i.evaluate(expr.Initializer)
		if err != nil {
			return
		}
	}
	i.Enviorment.Put(expr.Name.Lexme, value)
	return
}

func (i *Interpreter) VisitExprStmt(expr ExprStmt) (any, error) {
	_, err := i.evaluate(expr.Expr)
	return nil, err
}

func (i *Interpreter) VisitPrintStmt(expr PrintStmt) (_ any, err error) {
	v, err := i.evaluate(expr.Expr)
	if err != nil {
		return
	}
	fmt.Println(v)
	return
}

func (i *Interpreter) VisitBinaryExpr(expr BinaryExpr) (any, error) {
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
		left := left.(float64)
		right_num, ok := right.(float64)
		if !ok {
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

func (i *Interpreter) VisitGroupingExpr(expr GroupingExpr) (any, error) {
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) VisitLiteralExpr(expr LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(expr UnaryExpr) (any, error) {
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

func (i *Interpreter) isTruthy(value any) bool {
	switch value := value.(type) {
	case nil:
		return false
	case bool:
		return value
	default:
		return true
	}
}

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt Stmt) (any, error) {
	return stmt.Accept(i)
}

func Interpret(e []Stmt) error {
	i := &Interpreter{
		&Enviorment{
			values:    map[string]any{},
			enclosing: nil,
		},
	}
	for _, v := range e {
		_, err := i.execute(v)
		if err != nil {
			return err
		}
	}
	return nil
}
