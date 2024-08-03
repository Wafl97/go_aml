package logger

import "fmt"

type Level uint8

const (
	reset  = "\u001B[0m"
	red    = "\u001B[31m"
	green  = "\u001B[32m"
	yellow = "\u001B[33m"
	blue   = "\u001B[34m"
)

const (
	OFF   Level = 0
	ERROR Level = 1
	WARN  Level = 2
	INFO  Level = 3
	DEBUG Level = 4
)

var logLevel Level = INFO

func SetLogLevelByString(level string) {
	switch level {
	case "off":
		SetLogLevel(OFF)
	case "error":
		SetLogLevel(ERROR)
	case "warn":
		SetLogLevel(WARN)
	case "info":
		SetLogLevel(INFO)
	case "debug":
		SetLogLevel(DEBUG)
	}
}

func SetLogLevel(level Level) {
	logLevel = level
}

type Logger struct {
	name string
}

func New(name string) Logger {
	return Logger{
		name,
	}
}

func (logger *Logger) Debug(message string) {
	if logLevel >= DEBUG {
		fmt.Printf("%s[DEBUG] [%s] %s%s\n", blue, logger.name, message, reset)
	}
}

func (logger *Logger) Debugf(format string, args ...any) {
	if logLevel >= DEBUG {
		fmt.Printf("%s[DEBUG] [%s] %s%s\n", blue, logger.name, fmt.Sprintf(format, args...), reset)
	}
}

func (logger *Logger) Info(message string) {
	if logLevel >= INFO {
		fmt.Printf("%s[INFO ] [%s] %s%s\n", green, logger.name, message, reset)
	}
}

func (logger *Logger) Infof(format string, args ...any) {
	if logLevel >= INFO {
		fmt.Printf("%s[INFO ] [%s] %s%s\n", green, logger.name, fmt.Sprintf(format, args...), reset)
	}
}

func (logger *Logger) Warn(message string) {
	if logLevel >= WARN {
		fmt.Printf("%s[WARN ] [%s] %s%s\n", yellow, logger.name, message, reset)
	}
}

func (logger *Logger) Warnf(format string, args ...any) {
	if logLevel >= WARN {
		fmt.Printf("%s[WARN ] [%s] %s%s\n", yellow, logger.name, fmt.Sprintf(format, args...), reset)
	}
}

func (logger *Logger) Error(message string) {
	if logLevel >= ERROR {
		fmt.Printf("%s[ERROR] [%s] %s%s\n", red, logger.name, message, reset)
	}
}

func (logger *Logger) Errorf(format string, args ...any) {
	if logLevel >= ERROR {
		fmt.Printf("%s[ERROR] [%s] %s%s\n", red, logger.name, fmt.Sprintf(format, args...), reset)
	}
}
