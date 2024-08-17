package fsm

import (
	"fmt"
	"strings"
)

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

type Conditionals struct {
	IGenerate
	Conditions []Condition
}

func (conditionals Conditionals) Generate() string {
	switch len(conditionals.Conditions) {
	case 0:
		return "nil"
	default:
		conditionalStrings := make([]string, len(conditionals.Conditions))
		for i := 0; i < len(conditionalStrings); i++ {
			conditionalStrings[i] = (conditionals.Conditions)[i].ToString()
		}
		return fmt.Sprintf("func() bool { return %s }", strings.Join(conditionalStrings, " && "))
	}
}
