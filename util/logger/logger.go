package logger

import (
	"fmt"
	"os"
	"time"
)

type Level struct {
	order uint8
	name  string
	color string
}

const (
	reset  = "\u001B[0m"
	red    = "\u001B[31m"
	green  = "\u001B[32m"
	yellow = "\u001B[33m"
	blue   = "\u001B[34m"
)

var (
	levelOff   = &Level{0, "OFF", ""}
	levelError = &Level{1, "ERROR", red}
	levelWarn  = &Level{2, "WARN", yellow}
	levelInfo  = &Level{3, "INFO", green}
	levelDebug = &Level{4, "DEBUG", blue}
)

var out func(level *Level, s string) error = LogOutToConsole

var LogOutToDateFile func(level *Level, s string) error = func(level *Level, s string) error {
	file, err := os.OpenFile(time.Now().Format(time.DateOnly)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("logger out: %w", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("[%-5s] %s", level.name, s))
	if err != nil {
		return fmt.Errorf("logger out: %w", err)
	}
	return nil
}

var LogOutToConsole func(level *Level, s string) error = func(level *Level, s string) error {
	fmt.Printf("%s[%-5s] %s%s", level.color, level.name, s, reset)
	return nil
}

var logLevelOrder uint8 = levelInfo.order

func SetLogOut(outFunc func(level *Level, s string) error) {
	out = outFunc
}

func SetLogLevel(level string) {
	switch level {
	case "off", "OFF":
		setLogLevel(levelOff)
	case "error", "ERROR", "err", "ERR":
		setLogLevel(levelError)
	case "warn", "WARN", "warning", "WARNING":
		setLogLevel(levelWarn)
	case "info", "INFO":
		setLogLevel(levelInfo)
	case "debug", "DEBUG":
		setLogLevel(levelDebug)
	}
}

func setLogLevel(level *Level) {
	logLevelOrder = level.order
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
	if logLevelOrder >= levelDebug.order {
		out(levelDebug, fmt.Sprintf("[%s] %s\n", logger.name, message))
	}
}

func (logger *Logger) Debugf(format string, args ...any) {
	if logLevelOrder >= levelDebug.order {
		out(levelDebug, fmt.Sprintf("[%s] %s\n", logger.name, fmt.Sprintf(format, args...)))
	}
}

func (logger *Logger) Info(message string) {
	if logLevelOrder >= levelInfo.order {
		out(levelInfo, fmt.Sprintf("[%s] %s\n", logger.name, message))
	}
}

func (logger *Logger) Infof(format string, args ...any) {
	if logLevelOrder >= levelInfo.order {
		out(levelInfo, fmt.Sprintf("[%s] %s\n", logger.name, fmt.Sprintf(format, args...)))
	}
}

func (logger *Logger) Warn(message string) {
	if logLevelOrder >= levelWarn.order {
		out(levelWarn, fmt.Sprintf("[%s] %s\n", logger.name, message))
	}
}

func (logger *Logger) Warnf(format string, args ...any) {
	if logLevelOrder >= levelWarn.order {
		out(levelWarn, fmt.Sprintf("[%s] %s\n", logger.name, fmt.Sprintf(format, args...)))
	}
}

func (logger *Logger) Error(message string) {
	if logLevelOrder >= levelError.order {
		out(levelError, fmt.Sprintf("[%s] %s\n", logger.name, message))
	}
}

func (logger *Logger) Errorf(format string, args ...any) {
	if logLevelOrder >= levelError.order {
		out(levelError, fmt.Sprintf("[%s] %s\n", logger.name, fmt.Sprintf(format, args...)))
	}
}
