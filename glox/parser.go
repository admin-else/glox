package glox

import (
	"fmt"
)

type parser struct {
	tokens      []Token
	pos         int
	errors      []error
	syncronized bool
}

func ParseCode(code string) (stmts []Stmt, errs []error) {
	tokens, err := Scan(code)
	if err != nil {
		errs = append(errs, err)
		return
	}

	return Parse(tokens)
}

func Parse(tokens []Token) ([]Stmt, []error) {
	p := parser{tokens: tokens, pos: 0, syncronized: true}
	stmts := []Stmt{}
	for !p.isAtEnd() {
		stmt := p.decleration()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}
	return stmts, p.errors
}

func (p *parser) decleration() Stmt {
	var ret Stmt
	if p.match(VAR) {
		ret = p.varDecl()
	} else {
		ret = p.statement()
	}
	if !p.syncronized {
		p.syncronize()
		return nil
	}
	return ret
}

func (p *parser) varDecl() Stmt {
	name := p.consume(IDENTIFIER, "Expected variable name")

	var init Expr
	if p.match(EQUAL) {
		init = p.expression()
	}

	p.consume(SEMICOLON, "Expected ';' after variable expr")
	return VarDecl{Name: name, Initializer: init}
}

func (p *parser) statement() Stmt {
	if p.match(IF) {
		return p.ifStmt()
	}
	if p.match(PRINT) {
		return p.printStmt()
	}
	if p.match(LEFT_BRACE) {
		return p.block()
	}

	return p.exprStmt()
}

func (p *parser) block() Block {
	stmts := []Stmt{}

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.decleration())
	}
	p.consume(RIGHT_BRACE, "Expect '}' after blokc")
	return Block{stmts}
}

func (p *parser) exprStmt() ExprStmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expected semicolon")
	return ExprStmt{Expr: expr}
}

func (p *parser) printStmt() PrintStmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return PrintStmt{Expr: expr}
}

func (p *parser) ifStmt() IfStmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if contition")

	then := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return IfStmt{Condition: condition, ThenBranch: then, ElseBranch: elseBranch}
}

func (p *parser) expression() Expr {
	return p.assignment()
}

func (p *parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		equals := p.peek(-1)
		value := p.assignment()

		varr, ok := expr.(VariableExpr)
		if !ok {
			p.error(equals, "Invalid assignment target")
			return nil
		}
		return AssignExpr{varr.Name, value}
	}

	return expr
}

func (p *parser) or() Expr {
	expr := p.and()
	for p.match(OR) {
		op := p.peek(-1)
		right := p.and()
		expr = LogicalExpr{expr, op, right}
	}

	return expr
}

func (p *parser) and() Expr {
	expr := p.equality()
	for p.match(AND) {
		op := p.peek(-1)
		right := p.equality()
		expr = LogicalExpr{expr, op, right}
	}

	return expr
}

func (p *parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.peek(-1)
		right := p.comparison()
		expr = BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.peek(-1)
		right := p.term()
		expr = BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.peek(-1)
		right := p.factor()
		expr = BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.peek(-1)
		right := p.unary()
		expr = BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.peek(-1)
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

	if p.match(IDENTIFIER) {
		return VariableExpr{p.peek(-1)}
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
	p.errors = append(p.errors, fmt.Errorf("Error at line %v around %v: %s", t.Line, t.Lexme, message))
	p.syncronized = false
}

func (p *parser) syncronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.peek(-1).Type == SEMICOLON {
			return
		}

		t := p.peek(0).Type
		if t == CLASS || t == FOR || t == FUN || t == IF || t == PRINT || t == RETURN || t == VAR || t == WHILE {
			return
		}

		p.advance()
	}
	p.syncronized = true
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
	return p.peek(0).Type == t
}

func (p *parser) peek(offset int) Token {
	return p.tokens[p.pos+offset]
}

func (p *parser) isAtEnd() bool {
	return p.peek(0).Type == EOF
}
