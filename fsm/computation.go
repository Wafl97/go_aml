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

func (builder *FiniteStateMachineBuilder) NewComputations(function func(variables Variables)) []Computation {
	computations := make([]Computation, 0)
	builder.variables.listener = VarChangeListener{
		callback: func(a any, o ArithmeticSymbol, v VariableType) {
			builder.logger.Debugf("VAR CHANGE %v %v %s", a, o.String(), v.String())
		},
	}
	function(builder.variables)

	return computations
}

func (computation *Computation) ToString() string {
	return fmt.Sprintf("%s %s %v", computation.Left, computation.Operator.String(), computation.Right)
}

func (computation *Computation) Compute(variables *Variables) {
	originalValue, _ := variables.Get(computation.Left)
	switch computation.Operator {
	case ASSIGN:
		variables.Set(computation.Left, computation.ValueType, computation.Right)
	case ADD_ASSIGN:
		variables.Set(computation.Left, computation.ValueType, add(originalValue, computation.Right, computation.ValueType))
	case SUB_ASSIGN:
		variables.Set(computation.Left, computation.ValueType, sub(originalValue, computation.Right, computation.ValueType))
	case MUL_ASSIGN:
		variables.Set(computation.Left, computation.ValueType, mul(originalValue, computation.Right, computation.ValueType))
	case DIV_ASSIGN:
		variables.Set(computation.Left, computation.ValueType, div(originalValue, computation.Right, computation.ValueType))
	}
}

func add(a, b any, v VariableType) any {
	switch v {
	case FLOAT:
		af := a.(float64)
		bf := b.(float64)
		return af + bf
	case INT:
		ai := a.(int)
		bi := a.(int)
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
	Computations  []Computation
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
