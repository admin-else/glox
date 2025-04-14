package glox

import (
	"fmt"
)

type interpreter struct {
	lol string
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
	switch expr.Operator {
	case EQUAL_EQUAL:
		return left == right, nil
	case BANG_EQUAL:
		return left != right, nil
	}
	switch left.(type) {
	case float64:
		left, _ := left.(float64)
		right, ok := right.(float64)
		if ok == false {
			return nil, fmt.Errorf("binary expr with a left num must have a right num")
		}
		switch expr.Operator {
		case PLUS:
			return left + right, nil
		case MINUS:
			return left - right, nil
		case SLASH:
			return left / right, nil
		case STAR:
			return left * right, nil
		case GREATER:
			return left > right, nil
		case GREATER_EQUAL:
			return left >= right, nil
		case LESS:
			return left < right, nil
		case LESS_EQUAL:
			return left <= right, nil
		default:
			panic("Unreachable")
		}
	case string:
		left, _ := left.(string)
		right, ok := right.(string)
		if !ok {
			return nil, fmt.Errorf("binary epxr with a left string must have a right string")
		}
		if expr.Operator != PLUS {
			return nil, fmt.Errorf("binary expr on strings only support +")
		}
		return left + right, nil
	default:
		return nil, fmt.Errorf("binary expr does not accept this type")
	}
}

func (i *interpreter) VisitGroupingExpr(expr GroupingExpr) (any, error) {
	return i.evaluate(expr)
}

func (i *interpreter) VisitLiteralExpr(expr LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *interpreter) VisitUnaryExpr(expr UnaryExpr) (any, error) {
	right, err := i.evaluate(expr.Expr)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	case BANG:
		return i.isTruthy(right), nil
	case MINUS:
		v, ok := right.(float64)
		if !ok {
			return nil, fmt.Errorf("negation can only be done on numbers")
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

func Interpret(e Expr) {
	i := &interpreter{}
	e.Accept(i)
}
