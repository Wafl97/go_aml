package fsm

import "fmt"

type Condition struct {
	Left      string
	Symbol    LogicSymbol
	Right     any
	ValueType VariableType
}

func (condition *Condition) ToString() string {
	switch condition.ValueType {
	case BOOL:
		switch condition.Right {
		case "true":
			return condition.Left
		case "false":
			return fmt.Sprintf("!%s", condition.Left)
		}
	default:
		return fmt.Sprintf("%s %s %v", condition.Left, condition.Symbol.LSToString(), condition.Right)
	}
	return ""
}
