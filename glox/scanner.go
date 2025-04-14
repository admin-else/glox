package glox

import (
	"fmt"
	"strconv"
)

type TokenType int32

const (
	LEFT_PAREN    TokenType = iota + 1 // (
	RIGHT_PAREN                        // )
	LEFT_BRACE                         // {
	RIGHT_BRACE                        // }
	COMMA                              // ,
	DOT                                // .
	MINUS                              // -
	PLUS                               // +
	SEMICOLON                          // ;
	SLASH                              // /
	STAR                               // *
	BANG                               // !
	BANG_EQUAL                         // !=
	EQUAL                              // =
	EQUAL_EQUAL                        // ==
	GREATER                            // >
	GREATER_EQUAL                      // >=
	LESS                               // <
	LESS_EQUAL                         // <=
	IDENTIFIER                         // [a-zA-Z][a-zA-Z0-9]*
	STRING                             // "(.*)"
	NUMBER                             // (\d+|\d\.\d)
	AND                                // and
	CLASS                              // class
	ELSE                               // else
	FALSE                              // false
	FUN                                // fun
	FOR                                // for
	IF                                 // if
	NIL                                // nil
	OR                                 // or
	PRINT                              // print
	RETURN                             // return
	SUPER                              // super
	THIS                               // this
	TRUE                               // true
	VAR                                // var
	WHILE                              // while
	EOF
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Token struct {
	Type    TokenType
	Lexme   string
	Literal any
	Line    int
}

type scanner struct {
	code   string
	start  int
	pos    int
	line   int
	tokens []Token
}

func Scan(code string) ([]Token, error) {
	scanner := scanner{code: code, pos: 0, line: 1}
	err := scanner.scan()
	if err != nil {
		return nil, err
	}
	return scanner.tokens, nil
}

func (s *scanner) scan() error {
	for !s.isAtEnd() {
		c := s.advance()
		switch c {
		case '(':
			s.addToken(LEFT_PAREN, nil)
		case ')':
			s.addToken(RIGHT_PAREN, nil)
		case '{':
			s.addToken(LEFT_BRACE, nil)
		case '}':
			s.addToken(RIGHT_BRACE, nil)
		case ',':
			s.addToken(COMMA, nil)
		case '.':
			s.addToken(DOT, nil)
		case '-':
			s.addToken(MINUS, nil)
		case '+':
			s.addToken(PLUS, nil)
		case ';':
			s.addToken(SEMICOLON, nil)
		case '*':
			s.addToken(STAR, nil)
		case '!':
			if s.match('=') {
				s.addToken(BANG_EQUAL, nil)
			} else {
				s.addToken(BANG, nil)
			}
		case '=':
			if s.match('=') {
				s.addToken(EQUAL_EQUAL, nil)
			} else {
				s.addToken(EQUAL, nil)
			}
		case '<':
			if s.match('=') {
				s.addToken(LESS_EQUAL, nil)
			} else {
				s.addToken(LESS, nil)
			}
		case '>':
			if s.match('=') {
				s.addToken(GREATER_EQUAL, nil)
			} else {
				s.addToken(GREATER, nil)
			}
		case '/':
			if s.match('/') {
				for s.peek(0) != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else if s.peek(0) == '*' {
				for s.peek(0) != '*' || s.peek(1) != '/' {
					if s.advance() == '\n' {
						s.line += 1
					}
				}
				s.advance()
				s.advance()
			} else {
				s.addToken(SLASH, nil)
			}
		case '"':
			for s.peek(0) != '"' && !s.isAtEnd() {
				if s.peek(0) == '\n' {
					s.line++
				}
				s.advance()
			}

			if s.isAtEnd() {
				return fmt.Errorf("unterminated string at line %d", s.line)
			}

			s.advance()
			s.addToken(STRING, s.code[s.start+1:s.pos-1])
		case ' ':
		case '\r':
		case '\t':
		case '\n':
			s.line++
		default:
			if isDigit(c) {
				for isDigit(s.peek(0)) {
					s.advance()
				}

				if s.peek(0) == '.' && isDigit(s.peek(1)) {
					s.advance()
					for isDigit(s.peek(0)) {
						s.advance()
					}
				}

				num, err := strconv.ParseFloat(s.code[s.start:s.pos], 64)
				if err != nil {
					return fmt.Errorf("error while parsing number: %v", err)
				}
				s.addToken(NUMBER, num)
			} else if isApha(c) {
				for isAlpaNumeric(s.peek(0)) {
					s.advance()
				}
				text := s.code[s.start:s.pos]
				token := keywords[text]
				if token == 0 {
					token = IDENTIFIER
				}
				s.addToken(token, nil)
			} else {
				return fmt.Errorf("unexpected character '%c' at line %d", c, s.line)
			}
		}
		s.start = s.pos
	}
	s.addToken(EOF, nil)
	return nil
}

func (s *scanner) addToken(tokenType TokenType, literal any) {
	s.tokens = append(s.tokens, Token{Type: tokenType, Literal: literal, Line: s.line, Lexme: s.code[s.start:s.pos]})
}

func (s *scanner) advance() byte {
	c := s.code[s.pos]
	s.pos++
	return c
}

func (s *scanner) peek(offset int) byte {
	if s.isAtEnd() {
		return 0
	}
	return s.code[s.pos+offset]
}

func (s *scanner) isAtEnd() bool {
	return len(s.code) <= s.pos
}

func (s *scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.code[s.pos] != expected {
		return false
	}
	s.pos++
	return true
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isApha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlpaNumeric(c byte) bool {
	return isApha(c) || isDigit(c)
}
