package webster

import (
	"log"
)

type Interpreter struct {
	lexer *Lexer
	token Token
}

func NewInterpreter(input string) *Interpreter {
	lexer := NewLexer(input)
	return &Interpreter{lexer: lexer}
}

func (i *Interpreter) Consume(t TokenType) {
	if i.token.Type == t {
		i.token = i.lexer.NextToken()
	} else {
		log.Fatalf("Unexpected token: %v", i.token)
	}
}

func (i *Interpreter) Expr() interface{} {
	i.token = i.lexer.NextToken()

	result := i.Term()

	for i.token.Type == TokenPlus || i.token.Type == TokenMinus {
		if i.token.Type == TokenPlus {
			i.Consume(TokenPlus)
			result = i.add(result, i.Term())
		} else if i.token.Type == TokenMinus {
			i.Consume(TokenMinus)
			result = i.subtract(result, i.Term())
		}
	}

	return result
}

func (i *Interpreter) Term() interface{} {
	result := i.Factor()

	for i.token.Type == TokenMultiply || i.token.Type == TokenDivide {
		if i.token.Type == TokenMultiply {
			i.Consume(TokenMultiply)
			result = i.multiply(result, i.Factor())
		} else if i.token.Type == TokenDivide {
			i.Consume(TokenDivide)
			result = i.divide(result, i.Factor())
		}
	}

	return result
}

func (i *Interpreter) Factor() interface{} {
	var result interface{}

	switch i.token.Type {
	case TokenNumber:
		result = i.token.Value
		i.Consume(TokenNumber)
	case TokenString:
		result = i.token.Value.(string)
		i.Consume(TokenString)
	case TokenBool:
		result = i.token.Value.(bool)
		i.Consume(TokenBool)
	case TokenArray:
		result = i.token.Value.([]interface{})
		i.Consume(TokenArray)
	default:
		log.Fatalf("Unexpected token: %v", i.token)
	}

	return result
}

func (i *Interpreter) add(a interface{}, b interface{}) interface{} {
	switch a := a.(type) {
	case float64:
		if b, ok := b.(float64); ok {
			return a + b
		}
	case string:
		if b, ok := b.(string); ok {
			return a + b
		}
	case []interface{}:
		if b, ok := b.([]interface{}); ok {
			return append(a, b...)
		}
	default:
		log.Fatalf("Incompatible types for addition: %T and %T", a, b)
	}
	return nil
}

func (i *Interpreter) subtract(a interface{}, b interface{}) interface{} {
	if a, ok := a.(float64); ok {
		if b, ok := b.(float64); ok {
			return a - b
		}
	}
	log.Fatalf("Incompatible types for subtraction: %T and %T", a, b)
	return nil
}

func (i *Interpreter) multiply(a interface{}, b interface{}) interface{} {
	switch a := a.(type) {
	case float64:
		if b, ok := b.(float64); ok {
			return a * b
		}
	default:
		log.Fatalf("Incompatible types for multiplication: %T and %T", a, b)
	}
	return nil
}

func (i *Interpreter) divide(a interface{}, b interface{}) interface{} {
	if a, ok := a.(float64); ok {
		if b, ok := b.(float64); ok {
			return a / b
		}
	}
	log.Fatalf("Incompatible types for division: %T and %T", a, b)
	return nil
}
