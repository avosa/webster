package object

import (
	"fmt"

	"github.com/avosa/webster/lexer"
)

type Type int

const (
	UNKNOWN Type = iota
	INT
	FLOAT
	BOOLEAN
	STRING
	IDENT
	EOF
)

func (t Type) String() string {
	switch t {
	case UNKNOWN:
		return "UNKNOWN"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case BOOLEAN:
		return "BOOLEAN"
	case STRING:
		return "STRING"
	case IDENT:
		return "IDENT"
	case EOF:
		return "EOF"
	default:
		return "ILLEGAL"
	}
}

type Token struct {
	Type    Type
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("{Type:%v, Literal:%q}", t.Type, t.Literal)
}

var keywords = map[string]Type{
	"true":    BOOLEAN,
	"false":   BOOLEAN,
	"let":     IDENT,
	"var":     IDENT,
	"if":      IDENT,
	"else":    IDENT,
	"elseif":  IDENT,
	"for":     IDENT,
	"while":   IDENT,
	"fn":      IDENT,
	"return":  IDENT,
	"class":   IDENT,
	"extends": IDENT,
	"super":   IDENT,
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func NewToken(tokenType Type, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func LookupTokenType(literal string) Type {
	if _, err := lexer.ParseFloat(literal); err == nil {
		return FLOAT
	}
	if _, err := lexer.ParseInt(literal); err == nil {
		return INT
	}
	if _, ok := keywords[literal]; ok {
		return keywords[literal]
	}
	return IDENT
}
