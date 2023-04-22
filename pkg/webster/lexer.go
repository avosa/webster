package webster

import "strconv"

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenEOF
)

type Token struct {
	Type  TokenType
	Value float64
}

type Lexer struct {
	input   string
	pos     int
	current byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.pos >= len(l.input) {
		l.current = 0
	} else {
		l.current = l.input[l.pos]
	}
	l.pos++
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.current {
	case '+':
		tok = Token{Type: TokenPlus}
	case '-':
		tok = Token{Type: TokenMinus}
	case '*':
		tok = Token{Type: TokenMultiply}
	case '/':
		tok = Token{Type: TokenDivide}
	case 0:
		tok = Token{Type: TokenEOF}
	default:
		if isDigit(l.current) {
			value := l.readNumber()
			tok = Token{Type: TokenNumber, Value: value}
			return tok
		}
		tok = Token{Type: TokenEOF}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.current == ' ' || l.current == '\t' || l.current == '\n' || l.current == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() float64 {
	startPos := l.pos - 1

	for isDigit(l.current) {
		l.readChar()
	}

	value, _ := strconv.ParseFloat(l.input[startPos:l.pos-1], 64)
	return value
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
