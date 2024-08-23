package fsm

import "fmt"

type VariableType uint

func (vt VariableType) String() string {
	switch vt {
	case FLOAT:
		return "float"
	case INT:
		return "int"
	case BOOL:
		return "bool"
	case STRING:
		return "string"
	}
	return ""
}

type Listener interface {
	Accept(a any, o ArithmeticSymbol, v VariableType)
}

type VarChangeListener struct {
	Listener
	callback func(a any, o ArithmeticSymbol, v VariableType)
}

func (vcl VarChangeListener) Accept(a any, o ArithmeticSymbol, v VariableType) {
	vcl.callback(a, o, v)
}

const (
	FLOAT  VariableType = 0
	INT    VariableType = 1
	BOOL   VariableType = 2
	STRING VariableType = 3
)

type Variables struct {
	ints     map[string]int64
	floats   map[string]float64
	bools    map[string]bool
	strings  map[string]string
	types    map[string]VariableType
	listener Listener
}

func NewVariables() Variables {
	return Variables{
		ints:    map[string]int64{},
		floats:  map[string]float64{},
		bools:   map[string]bool{},
		strings: map[string]string{},
		types:   map[string]VariableType{},
	}
}

func (variables *Variables) Update(key string, operation ArithmeticSymbol, value any) {
	originalValue, hasValue := variables.Get(key)
	if !hasValue {
		fmt.Printf("FAILED %s\n", key)
	}
	switch operation {
	case ASSIGN:
		variables.Set(key, variables.types[key], value)
	case ADD_ASSIGN:
		variables.Set(key, variables.types[key], add(originalValue, value, variables.types[key]))
	case SUB_ASSIGN:
		variables.Set(key, variables.types[key], sub(originalValue, value, variables.types[key]))
	case MUL_ASSIGN:
		variables.Set(key, variables.types[key], mul(originalValue, value, variables.types[key]))
	case DIV_ASSIGN:
		variables.Set(key, variables.types[key], div(originalValue, value, variables.types[key]))
	}

	if variables.listener != nil {
		variables.listener.Accept(value, operation, variables.types[key])
	}
}

func (variables *Variables) Set(key string, varType VariableType, value any) {
	variables.types[key] = varType
	switch varType {
	case INT:
		variables.ints[key] = value.(int64)
	case FLOAT:
		variables.floats[key] = value.(float64)
	case BOOL:
		variables.bools[key] = value.(bool)
	case STRING:
		variables.strings[key] = value.(string)
	default:
		fmt.Println("FAILED TO SET VALUE")
	}
}

func (variables *Variables) SetType(key string, value VariableType) {
	fmt.Printf("%s TYPE = %s", key, value.String())
	variables.types[key] = value
}

func (variables *Variables) Get(key string) (any, bool) {

	t := variables.types[key]
	switch t {
	case INT:
		v, b := variables.ints[key]
		return v, b
	case FLOAT:
		v, b := variables.floats[key]
		return v, b
	case BOOL:
		v, b := variables.bools[key]
		return v, b
	case STRING:
		v, b := variables.strings[key]
		return v, b
	}
	return nil, false
}

func (variables *Variables) GetType(key string) VariableType {
	return variables.types[key]
}
