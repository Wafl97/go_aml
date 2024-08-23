package fsm

import (
	"fmt"
	"github.com/Wafl97/go_aml/util/functions"
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

func (condition *Condition) Evaluate(variables *Variables) bool {
	originalValue, _ := variables.Get(condition.left)
	switch condition.symbol {
	case EQUAL:
		return originalValue == condition.right
	case NOT_EQUAL:
		return originalValue != condition.right
	case GRATER_THAN:
		return grater(originalValue, condition.right, condition.valueType)
	case GRATER_THAN_OR_EQUAL:
		return grater_equal(originalValue, condition.right, condition.valueType)
	case LESS_THAN:
		return less(originalValue, condition.right, condition.valueType)
	case LESS_THAN_OR_EQUAL:
		return less_equal(originalValue, condition.right, condition.valueType)
	}
	return false
}

func grater(a, b any, v VariableType) bool {
	switch v {
	case FLOAT:
		af := a.(float32)
		bf := b.(float32)
		return af > bf
	case INT:
		ai := a.(int32)
		bi := a.(int32)
		return ai > bi
	case BOOL:
		return false
	case STRING:
		return false
	}
	return false
}

func grater_equal(a, b any, v VariableType) bool {
	switch v {
	case FLOAT:
		af := a.(float32)
		bf := b.(float32)
		return af >= bf
	case INT:
		ai := a.(int32)
		bi := a.(int32)
		return ai >= bi
	case BOOL:
		return false
	case STRING:
		return false
	}
	return false
}

func less(a, b any, v VariableType) bool {
	switch v {
	case FLOAT:
		af := a.(float32)
		bf := b.(float32)
		return af < bf
	case INT:
		ai := a.(int32)
		bi := a.(int32)
		return ai < bi
	case BOOL:
		return false
	case STRING:
		return false
	}
	return false
}

func less_equal(a, b any, v VariableType) bool {
	switch v {
	case FLOAT:
		af := a.(float32)
		bf := b.(float32)
		return af <= bf
	case INT:
		ai := a.(int32)
		bi := a.(int32)
		return ai <= bi
	case BOOL:
		return false
	case STRING:
		return false
	}
	return false
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
		return fmt.Sprintf("%s %s %v", condition.left, condition.symbol.String(), condition.right)
	}
	return ""
}

type Conditionals struct {
	conditions []Condition
	function   functions.Predicate[*Variables]
}

func NewConditionals(conditions []Condition, function functions.Predicate[*Variables]) Conditionals {
	if conditions != nil && function != nil {
		return Conditionals{conditions, function}
	}
	if function != nil {
		return Conditionals{function: function}
	}
	return Conditionals{
		conditions: conditions,
		function: func(v *Variables) bool {
			for _, c := range conditions {
				if !c.Evaluate(v) {
					return false
				}
			}
			return true
		},
	}
}

func (conditionals *Conditionals) Generate() string {
	switch len(conditionals.conditions) {
	case 0:
		return "nil"
	default:
		conditionalStrings := make([]string, len(conditionals.conditions))
		for i := 0; i < len(conditionalStrings); i++ {
			conditionalStrings[i] = (conditionals.conditions)[i].ToString()
		}
		return fmt.Sprintf("func() bool { return %s }", strings.Join(conditionalStrings, " && "))
	}
}

func (conditionals *Conditionals) Evaluate(variables *Variables) bool {
	return conditionals.function(variables)
}
