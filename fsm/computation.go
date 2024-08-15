package fsm

import "fmt"

type Computation struct {
	Left      string
	Operator  ArithmeticSymbol
	Right     any
	ValueType VariableType
}

func (computation *Computation) ToString() string {
	return fmt.Sprintf("%s %s %v", computation.Left, computation.Operator.ASToString(), computation.Right)
}
