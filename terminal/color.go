package terminal

import "fmt"

// ColorSchema contains success/failure/neutral color
type ColorSchema struct {
	SuccessColor int
	FailureColor int
	NeutralColor int
}

// ResetSequence reset ANSI sequence.
const ResetSequence = "\033[0m"

const (
	// BlackColor color
	BlackColor = iota
	// RedColor color
	RedColor
	// GreenColor color
	GreenColor
	// YellowColor color
	YellowColor
	// BlueColor color
	BlueColor
	// MagentaColor color
	MagentaColor
	// CyanColor color
	CyanColor
	// WhiteColor color
	WhiteColor
)

// DefaultColorSchema will returns a color schema
//
// By default, success color is green, failure color is red, and neutral color is yellow
func DefaultColorSchema() *ColorSchema {
	return &ColorSchema{
		SuccessColor: GreenColor,
		FailureColor: RedColor,
		NeutralColor: YellowColor,
	}
}

// GetColoredText will colorize text with color
func GetColoredText(color int, text string) string {
	return fmt.Sprintf("%s%s%s", GetColorSequence(color), text, ResetSequence)
}

// GetColorSequence will returns color control escape sequence
func GetColorSequence(code int) string {
	if code < 0 {
		return ResetSequence
	}
	return fmt.Sprintf("\033[3%dm", code)
}
