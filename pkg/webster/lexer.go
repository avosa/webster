package webster

import (
	"strconv"
)

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenString
	TokenBool
	TokenArray
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenEOF
)

type Token struct {
	Type  TokenType
	Value interface{}
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
		if l.peekChar() == '/' {
			l.skipSingleLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipMultiLineComment()
			return l.NextToken()
		} else {
			tok = Token{Type: TokenDivide}
		}
	case 0:
		tok = Token{Type: TokenEOF}
	case '"':
		tok = Token{Type: TokenString, Value: l.readString()}
		return tok
	case '[':
		tok = Token{Type: TokenArray, Value: l.readArray()}
		return tok
	default:
		if isDigit(l.current) {
			value := l.readNumber()
			tok = Token{Type: TokenNumber, Value: value}
			return tok
		} else if l.current == 't' || l.current == 'f' {
			value := l.readBoolean()
			tok = Token{Type: TokenBool, Value: value}
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

func (l *Lexer) skipSingleLineComment() {
	for l.current != '\n' && l.current != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipMultiLineComment() {
	for l.current != 0 {
		l.readChar()
		if l.current == '*' && l.peekChar() == '/' {
			l.readChar()
			l.readChar()
			break
		}
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

func (l *Lexer) readString() string {
	startPos := l.pos

	for {
		l.readChar()
		if l.current == '"' || l.current == 0 {
			break
		}
	}

	return l.input[startPos : l.pos-1]
}

func (l *Lexer) readArray() []interface{} {
	var arr []interface{}

	for l.current != ']' && l.current != 0 {
		if isDigit(l.current) {
			value := l.readNumber()
			arr = append(arr, value)
		} else if l.current == '"' {
			value := l.readString()
			arr = append(arr, value)
		} else if l.current == 't' || l.current == 'f' {
			value := l.readBoolean()
			arr = append(arr, value)
		}
		l.readChar()
	}

	return arr
}

func (l *Lexer) readBoolean() bool {
	startPos := l.pos - 1

	for l.current != ' ' && l.current != '\t' && l.current != '\n' && l.current != '\r' && l.current != 0 {
		l.readChar()
	}

	if l.input[startPos:l.pos-1] == "true" {
		return true
	}
	return false
}

func (l *Lexer) peekChar() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
