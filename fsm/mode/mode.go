package mode

type Mode int

const (
	CONTINUE   Mode = 0
	DEADLOCK   Mode = 1
	TERMINATED Mode = 2
	CRASH      Mode = 3
)
