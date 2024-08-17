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
