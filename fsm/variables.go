package fsm

type Variables struct {
	values map[string]any
}

func NewVariables() Variables {
	return Variables{
		values: map[string]any{},
	}
}

func (variables *Variables) Set(key string, value any) {
	variables.values[key] = value
}

func (variables *Variables) Get(key string) any {
	return variables.values[key]
}

func GetAndCast[T any](variables *Variables, key string) T {
	return variables.Get(key).(T)
}
