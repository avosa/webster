package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"math"
)

type Type string

const (
	INTEGER_OBJ      = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	ARRAY_OBJ        = "ARRAY"
	DICTIONARY_OBJ   = "DICTIONARY"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Float struct {
	Value float64
}

func (f *Float) Type() Type {
	return FLOAT_OBJ
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%g", f.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() Type {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type {
	return RETURN_VALUE_OBJ
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() Type {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return "Error: " + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Env        *Environment
	Body       *ast.BlockStatement
	Name       string
}

func (f *Function) Type() Type {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("func ")
	if f.Name != "" {
		out.WriteString(f.Name)
	}
	out.WriteString("(")
	out.WriteString(joinStrings(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() Type {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

type Array struct {
	Elements []Object
}

func (ao *Array) Type() Type {
	return ARRAY_OBJ
}

func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range ao.Elements {
		elements = append(elements, el.Inspect())
	}
	out.WriteString("[")
	out.WriteString(joinStrings(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type DictionaryPair struct {
	Key   Hashable
	Value Object
}

type Dictionary struct {
	Pairs map[HashKey]DictionaryPair
}

func (d *Dictionary) Type() Type {
	return DICTIONARY_OBJ
}

func (d *Dictionary) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range d.Pairs {
		pairs = append(pairs, pair.Key.String()+": "+pair.Value.Inspect())
	}
	out.WriteString("{")
	out.WriteString(joinStrings(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  Type
	Value uint64
}

type IntegerHashKey struct {
	Value int64
}

func (ik *IntegerHashKey) HashKey() HashKey {
	return HashKey{Type: INTEGER_OBJ, Value: uint64(ik.Value)}
}

type FloatHashKey struct {
	Value float64
}

func (fk *FloatHashKey) HashKey() HashKey {
	bits := math.Float64bits(fk.Value)
	return HashKey{Type: FLOAT_OBJ, Value: bits}
}

type BooleanHashKey struct {
	Value bool
}

func (bk *BooleanHashKey) HashKey() HashKey {
	var value uint64
	if bk.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: BOOLEAN_OBJ, Value: value}
}

type StringHashKey struct {
	Value string
}

func (sk *StringHashKey) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(sk.Value))
	return HashKey{Type: STRING_OBJ, Value: h.Sum64()}
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}
