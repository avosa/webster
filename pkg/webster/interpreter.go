package webster

import "log"

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

func (i *Interpreter) Expr() float64 {
	i.token = i.lexer.NextToken()

	result := i.Term()

	for i.token.Type == TokenPlus || i.token.Type == TokenMinus {
		if i.token.Type == TokenPlus {
			i.Consume(TokenPlus)
			result += i.Term()
		} else if i.token.Type == TokenMinus {
			i.Consume(TokenMinus)
			result -= i.Term()
		}
	}

	return result
}

func (i *Interpreter) Term() float64 {
	result := i.Factor()

	for i.token.Type == TokenMultiply || i.token.Type == TokenDivide {
		if i.token.Type == TokenMultiply {
			i.Consume(TokenMultiply)
			result *= i.Factor()
		} else if i.token.Type == TokenDivide {
			i.Consume(TokenDivide)
			result /= i.Factor()
		}
	}

	return result
}

func (i *Interpreter) Factor() float64 {
	result := i.token.Value
	i.Consume(TokenNumber)
	return result
}
