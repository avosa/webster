package evaluator

import (
	"fmt"
	"go/token"

	"github.com/avosa/webster/ast"
	"github.com/avosa/webster/object"
)

var (
	// NULL represents the null object.
	NULL = &Null{}
	// TRUE represents the true boolean object.
	TRUE = &Boolean{Value: true}
	// FALSE represents the false boolean object.
	FALSE = &Boolean{Value: false}
)

// Eval evaluates an AST node and returns an Object.
func Eval(node ast.Node, env *Environment) Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &Float{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.StringLiteral:
		return &String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.DictionaryLiteral:
		pairs := make(map[HashKey]Pair)
		for keyNode, valueNode := range node.Pairs {
			key := Eval(keyNode, env)
			if isError(key) {
				return key
			}
			hashKey, ok := key.(Hashable)
			if !ok {
				return newError("unusable as dictionary key: %s", key.Type())
			}
			value := Eval(valueNode, env)
			if isError(value) {
				return value
			}
			pairs[hashKey.HashKey()] = Pair{Key: key, Value: value}
		}
		return &object.Dictionary{Pairs: pairs}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.AssignExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		if node.Token.Type == token.ASSIGN {
			return evalAssignExpression(node.Left, right, env)
		}
		return evalIncompatibleAssign(node.Left, right)
	default:
		return newError("unknown expression type: %T", node)
	}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.DICTIONARY_OBJ:
		return evalDictionaryIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	strObj := str.(*object.String)
	idx := index.(*object.Integer).Value
	if idx < 0 || idx >= int64(len(strObj.Value)) {
		return NULL
	}
	return &object.String{Value: string(strObj.Value[idx])}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObj.Elements[idx]
}

func evalDictionaryIndexExpression(dic, key object.Object) object.Object {
	dicObj := dic.(*object.Dictionary)
	hashKey, ok := key.(Hashable)
	if !ok {
		return newError("unusable as dictionary key: %s", key.Type())
	}
	pair, ok := dicObj.Pairs[hashKey.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalAssignExpression(left ast.Expression, right object.Object, env *Environment) object.Object {
	switch left := left.(type) {
	case *ast.Identifier:
		return evalAssignIdentifier(left, right, env)
	case *ast.IndexExpression:
		return evalAssignIndexExpression(left, right, env)
	default:
		return newError("invalid assignment target: %s", left.String())
	}
}

func evalIncompatibleAssign(left ast.Expression, right Object) *Error {
	return NewError(ErrorTypeInvalidAssignment,
		fmt.Sprintf("cannot assign %s to %s", right.Type(), left.String()))
}

func evalIndexExpression(indexable, index ast.Expression, env *Environment) Object {
	switch {
	case indexable.Type() == ObjectTypeArray && index.Type() == ObjectTypeInteger:
		return evalArrayIndexExpression(indexable, index, env)
	case indexable.Type() == ObjectTypeDictionary:
		return evalDictionaryIndexExpression(indexable, index, env)
	default:
		return NewError(ErrorTypeInvalidIndexExpression,
			fmt.Sprintf("index operator not supported: %s[%s]", indexable.Type(), index.Type()))
	}
}

func evalArrayIndexExpression(array, index ast.Expression, env *Environment) Object {
	arr := array.(*Array)
	idx := index.(*Integer).Value
	max := int64(len(arr.Elements) - 1)
	if idx < 0 || idx > max {
		return Null
	}
	return arr.Elements[idx]
}

func evalDictionaryIndexExpression(dict, key ast.Expression, env *Environment) Object {
	dictObj := dict.(*Dictionary)
	hashKey, ok := key.(Hashable)
	if !ok {
		return NewError(ErrorTypeInvalidHashKey,
			fmt.Sprintf("unusable as hash key: %s", key.Type()))
	}
	pair, ok := dictObj.Pairs[hashKey.HashKey()]
	if !ok {
		return Null
	}
	return pair.Value
}

func evalIfExpression(ie *ast.IfExpression, env *Environment) Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return Null
	}
}

func evalWhileExpression(we *ast.WhileExpression, env *Environment) Object {
	var result Object
	condition := Eval(we.Condition, env)
	if isError(condition) {
		return condition
	}
	for isTruthy(condition) {
		result = Eval(we.Body, env)
		if isError(result) {
			return result
		}
		condition = Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}
	}
	return result
}

func evalForExpression(fe *ast.ForExpression, env *Environment) Object {
	var result Object
	iterator := Eval(fe.Iterable, env)
	if isError(iterator) {
		return iterator
	}
	switch iterable := iterator.(type) {
	case *Array:
		result = evalForArrayExpression(iterable, fe, env)
	case *Dictionary:
		result = evalForDictionaryExpression(iterable, fe, env)
	default:
		return NewError(ErrorTypeInvalidIterable,
			fmt.Sprintf("'%s' cannot be used as iterable", iterable.Type()))
	}
	return result
}

func evalForArrayExpression(array *Array, fe *ast.ForExpression, env *Environment) Object {
	var result Object
	for i, element := range array.Elements {
		newEnv := env.Copy()
		newEnv.Set(fe.Variable.Value, element)
		if i == len(array.Elements)-1 {
			newEnv.Set(fe.IsLast.Value, TRUE)
		} else {
			newEnv.Set(fe.IsLast.Value, FALSE)
		}
		result = Eval(fe.Body, newEnv)
		if isError(result) {
			return result
		}
	}
	return result
}

func evalForDictionaryExpression(dict *Dictionary, fe *ast.ForExpression, env *Environment) Object {
	var result Object

	for _, pair := range dict.Pairs {
		newEnv := extendForEnvironment(fe.Key.String(), pair.Key, fe.Value.String(), pair.Value, env)
		if newEnv == nil {
			return newError("for loop environment error")
		}
		result = Eval(fe.Body, newEnv)
		if isError(result) {
			return result
		}
	}

	return result
}

func extendForEnvironment(keyName string, keyValue Object, valueName string, valueValue Object, env *Environment) *Environment {
	newEnv := NewEnclosedEnvironment(env)

	err := newEnv.Set(keyName, keyValue)
	if err != nil {
		return nil
	}

	err = newEnv.Set(valueName, valueValue)
	if err != nil {
		return nil
	}

	return newEnv
}
