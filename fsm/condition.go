package fsm

import (
	"fmt"
	"strings"
)

type Condition struct {
	left      string
	symbol    LogicSymbol
	right     any
	valueType VariableType
}

func NewCondition(left string, symbol LogicSymbol, right any, valueType VariableType) Condition {
	return Condition{
		left:      left,
		symbol:    symbol,
		right:     right,
		valueType: valueType,
	}
}

func (condition *Condition) ToString() string {
	switch condition.valueType {
	case BOOL:
		switch condition.right {
		case "true":
			return condition.left
		case "false":
			return fmt.Sprintf("!%s", condition.left)
		}
	default:
		return fmt.Sprintf("%s %s %v", condition.left, condition.symbol.LSToString(), condition.right)
	}
	return ""
}

type Conditionals struct {
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
