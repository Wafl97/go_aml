package types

type Option[T any] struct { // DEPRECATED
	value    T
	hasValue bool
}

func None[T any]() Option[T] {
	return Option[T]{
		hasValue: false,
	}
}

func Some[T any](value T) Option[T] {
	return Option[T]{
		value:    value,
		hasValue: true,
	}
}

func (option Option[T]) Get() T {
	return option.value
}

func (option Option[T]) GetOrElse(other T) T {
	if option.hasValue {
		return option.value
	}
	return other
}

func (option Option[T]) IsSome() bool {
	return option.hasValue
}

func (option Option[T]) IsNone() bool {
	return !option.hasValue
}

func (option Option[T]) HasValue(callable func(T)) Option[T] {
	if option.hasValue {
		callable(option.value)
	}
	return option
}

func (option Option[T]) Else(callable func()) Option[T] {
	if !option.hasValue {
		callable()
	}
	return option
}

func (option Option[T]) OrElse(defaultValue T, callable func(T)) Option[T] {
	if option.hasValue {
		callable(option.value)
	} else {
		callable(defaultValue)
	}
	return option
}
