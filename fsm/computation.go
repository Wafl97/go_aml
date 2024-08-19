package fsm

import (
	"fmt"
	"strings"
)

type Computation struct {
	Left      string
	Operator  ArithmeticSymbol
	Right     any
	ValueType VariableType
}

func (computation *Computation) ToString() string {
	return fmt.Sprintf("%s %s %v", computation.Left, computation.Operator.ASToString(), computation.Right)
}

func (computation *Computation) Compute(variables *Variables) {
	switch computation.Operator {
	case ASSIGN:
		variables.Set(computation.Left, computation.Right)
	case ADD_ASSIGN:
		variables.Set(computation.Left, add(variables.Get(computation.Left), computation.Right, computation.ValueType))
	case SUB_ASSIGN:
		variables.Set(computation.Left, sub(variables.Get(computation.Left), computation.Right, computation.ValueType))
	case MUL_ASSIGN:
		variables.Set(computation.Left, mul(variables.Get(computation.Left), computation.Right, computation.ValueType))
	case DIV_ASSIGN:
		variables.Set(computation.Left, div(variables.Get(computation.Left), computation.Right, computation.ValueType))
	}
}

func add(a, b any, v VariableType) any {
	switch v {
	case FLOAT:
		af := a.(float32)
		bf := b.(float32)
		return af + bf
	case INT:
		ai := a.(int32)
		bi := a.(int32)
		return ai + bi
	case BOOL:
		return nil
	case STRING:
		as := a.(string)
		bs := b.(string)
		return as + bs
	}
	return nil
}

func sub(a, b any, v VariableType) any {
	switch v {
	case FLOAT:
		af := a.(float32)
		bf := b.(float32)
		return af - bf
	case INT:
		ai := a.(int32)
		bi := a.(int32)
		return ai - bi
	case BOOL:
		return nil
	case STRING:
		return nil
	}
	return nil
}

func mul(a, b any, v VariableType) any {
	switch v {
	case FLOAT:
		af := a.(float32)
		bf := b.(float32)
		return af * bf
	case INT:
		ai := a.(int32)
		bi := a.(int32)
		return ai * bi
	case BOOL:
		return nil
	case STRING:
		return nil
	}
	return nil
}

func div(a, b any, v VariableType) any {
	switch v {
	case FLOAT:
		af := a.(float32)
		bf := b.(float32)
		return af / bf
	case INT:
		ai := a.(int32)
		bi := a.(int32)
		return ai / bi
	case BOOL:
		return nil
	case STRING:
		return nil
	}
	return nil
}

type Computational struct {
	FuncSignature string
	Computations  []*Computation
}

func (computational Computational) Generate() string {
	switch len(computational.Computations) {
	case 0:
		return "nil"
	default:
		computationStrings := make([]string, len(computational.Computations))
		for i := 0; i < len(computationStrings); i++ {
			computationStrings[i] = (computational.Computations)[i].ToString()
		}
		return fmt.Sprintf("%s { %s }", computational.FuncSignature, strings.Join(computationStrings, "; "))
	}
}
