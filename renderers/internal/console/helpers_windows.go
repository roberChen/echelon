package console

import (
	"golang.org/x/sys/windows"
	"os"
)

// PrepareTerminalEnvironment for windows platform.
func PrepareTerminalEnvironment() error {
	// enable handling ASCII codes
	err := addConsoleMode(windows.Stdout, windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	if err != nil {
		return err
	}
	return addConsoleMode(windows.Stderr, windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}

func addConsoleMode(handle windows.Handle, flags uint32) error {
	var mode uint32

	err := windows.GetConsoleMode(handle, &mode)
	if err != nil {
		return err
	}
	return windows.SetConsoleMode(handle, mode|flags)
}

// TerminalHeight will return terminal height of windows platform, currently this
// function will return -1
func TerminalHeight(file *os.File) int {
	// TODO: figure out how to find out console height on Windows
	return -1
}

// TerminalWidth will return terminal width of windows platform, currently this
// function will return -1
func TerminalWidth(file *os.File) int {
	// TODO: figure out how to find out console width on Windows
	return -1
}
