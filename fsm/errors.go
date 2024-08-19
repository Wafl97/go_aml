package fsm

import (
	"errors"
	"fmt"
)

type (
	DeadlockError struct {
		cause string
	}
	CrashError struct {
		cause string
	}
)

var (
	TerminateError error = errors.New("fsm_terminate: model has terminated")
)

func NewDeadlockError(cause string) DeadlockError {
	return DeadlockError{cause: cause}
}

func NewDeadlockErrorf(format string, args ...any) DeadlockError {
	return DeadlockError{cause: fmt.Sprintf(format, args...)}
}

func NewCrashError(cause string) CrashError {
	return CrashError{cause: cause}
}

func NewCrashErrorf(format string, args ...any) CrashError {
	return CrashError{cause: fmt.Sprintf(format, args...)}
}

func (e DeadlockError) Error() string {
	return fmt.Sprintf("fsm_deadlock: %s", e.cause)
}

func (e CrashError) Error() string {
	return fmt.Sprintf("fsm_crash: %s", e.cause)
}
