//go:build windows

package cfg

import (
	"os"
	"syscall"
)

func init() {
	stdout := syscall.Handle(os.Stdout.Fd())
	var mode uint32
	syscall.GetConsoleMode(stdout, &mode)
	mode |= 0x0004
	syscall.MustLoadDLL("kernel32").MustFindProc("SetConsoleMode").Call(uintptr(stdout), uintptr(mode))
}
