package glox

import (
	"fmt"
)

type parser struct {
	tokens []Token
	pos    int
	errors []error
}

func Parse(tokens []Token) (Expr, []error) {
	p := parser{tokens: tokens, pos: 0}
	return p.expression(), p.errors
}

func (p *parser) expression() Expr {
	return p.equality()
}

func (p *parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.peek(-1).TokenType
		right := p.comparison()
		expr = BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.peek(-1).TokenType
		right := p.term()
		expr = BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.peek(-1).TokenType
		right := p.factor()
		expr = BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.peek(-1).TokenType
		right := p.unary()
		expr = BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.peek(-1).TokenType
		right := p.unary()
		return UnaryExpr{Operator: operator, Expr: right}
	}

	return p.primary()
}
func (p *parser) primary() Expr {

	if p.match(FALSE) {
		return LiteralExpr{false}
	}
	if p.match(TRUE) {
		return LiteralExpr{true}
	}
	if p.match(NIL) {
		return LiteralExpr{nil}
	}

	if p.match(NUMBER, STRING) {
		return LiteralExpr{p.peek(-1).Literal}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expected ')' after expression.")
		return GroupingExpr{Expr: expr}
	}

	p.error(p.peek(0), "Expected Expression")
	return nil
}

func (p *parser) consume(t TokenType, message string) Token {
	if p.check(t) {
		return p.advance()
	}

	p.error(p.peek(0), message)
	return Token{}
}

func (p *parser) error(t Token, message string) {
	p.errors = append(p.errors, fmt.Errorf("Error at line %v around %v: %s\n", t.Line, t.Lexme, message))
	p.syncronize()
}

func (p *parser) syncronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.peek(-1).TokenType == SEMICOLON {
			return
		}

		t := p.peek(0).TokenType
		if t == CLASS || t == FOR || t == FUN || t == IF || t == PRINT || t == RETURN || t == VAR || t == WHILE {
			return
		}

		p.advance()
	}
}

// Utils
func (p *parser) match(tokenTypes ...TokenType) bool {
	for _, t := range tokenTypes {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) advance() Token {
	if !p.isAtEnd() {
		p.pos++
	}
	return p.peek(-1)
}

func (p *parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek(0).TokenType == t
}

func (p *parser) peek(offset int) Token {
	return p.tokens[p.pos+offset]
}

func (p *parser) isAtEnd() bool {
	return p.peek(0).TokenType == EOF
}
