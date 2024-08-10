package fsm

type VariableType uint

const (
	FLOAT  VariableType = 0
	INT    VariableType = 1
	BOOL   VariableType = 2
	STRING VariableType = 3
)

type Variables struct {
	values map[string]any
	types  map[string]VariableType
}

func NewVariables() Variables {
	return Variables{
		values: map[string]any{},
		types:  map[string]VariableType{},
	}
}

func (variables *Variables) Set(key string, value any) {
	variables.values[key] = value
}

func (variables *Variables) SetType(key string, value VariableType) {
	variables.types[key] = value
}

func (variables *Variables) Get(key string) any {
	return variables.values[key]
}

func (variables *Variables) GetType(key string) VariableType {
	return variables.types[key]
}

func GetAndCast[T any](variables *Variables, key string) T {
	return variables.Get(key).(T)
}
