// +build !windows

package console

import (
	"golang.org/x/sys/unix"
	"os"
)

// PrepareTerminalEnvironment will prepare for terminal environment, for unix platform,
// nothing needs to be done.
func PrepareTerminalEnvironment() error {
	// no need on unix
	return nil
}

// TerminalHeight returns the height of current terminal height
func TerminalHeight(file *os.File) int {
	ws, err := unix.IoctlGetWinsize(int(file.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return -1
	}

	return int(ws.Row)
}

// TerminalWidth returns the width of current terminal height
func TerminalWidth(file *os.File) int {
	ws, err := unix.IoctlGetWinsize(int(file.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return -1
	}

	return int(ws.Col)
}