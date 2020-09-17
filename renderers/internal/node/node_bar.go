package node

import (
	"fmt"
	"strings"
	"sync"

	"github.com/roberChen/echelon/terminal"
)

const (
	// DefaultStyle is the default bar style
	DefaultStyle = "╢▌▌░╟"
	// SimpleStyle is simple style of bar, similar to the bar of wget
	SimpleStyle = "[=>-]"
)

// Bar is a structure which defines the horizontal progress bar
type Bar struct {
	/*
			'1st rune' stands for left boundary rune

		    '2nd rune' stands for fill rune

		    '3rd rune' stands for tip rune

		    '4th rune' stands for space rune

		    '5th rune' stands for right boundary rune

			"╢▌▌░╟" by default
	*/
	lbound, fill, tip, space, rbound rune

	mu sync.Mutex

	// total size of task
	total int64
	// precent progress of task
	now int64
	// percentage of progress, 0-100
	percentage int
}

// SetProgress set the progress of bar, it's a coroutine safe function
func (b *Bar) SetProgress(i int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.now = i
	b.percentage = int(100 * i / b.total)
}

// AddProgress add the progress of bar, it's a coroutine safe function
func (b *Bar) AddProgress(i int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.now += i
	if b.now > b.total {
		b.now = b.total
	}
	b.percentage = int(100 * b.now / b.total)
}

// IsFinished returns whether bar has done
func (b *Bar) IsFinished() bool {
	return b.percentage == 100
}

// SetPercentage sets the percentage of bar
func (b *Bar) SetPercentage(i int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if i > 100 {
		b.percentage = 100
	} else {
		b.percentage = i
	}
	b.now = b.total * int64(b.percentage)
}

// AddPercentage adds percentage for bar, it's a coroutine safe function
func (b *Bar) AddPercentage(i int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.percentage += i
	if b.percentage > 100 {
		b.percentage = 100
	}
	b.now = b.total * int64(b.percentage)
}

// String returns the render of bar
func (b *Bar) String(width int) string {
	if width <= 2 {
		return ""
	}
	width -= 2
	finished := int(width * b.percentage / 100)
	remains := width - 1 - finished
	// finished, no tip needed
	if finished == width {
		remains++
	}
	if finished <0 || remains < 0 {
		fmt.Printf("ERROR: negative repeat: finished:%d\tremains:%d\n", finished, remains)
	}
	var out string
	if !b.IsFinished() {
		out = fmt.Sprintf("%c%s%c%s%c", b.lbound, strings.Repeat(string(b.fill), finished),
			b.tip, strings.Repeat(string(b.space), remains), b.rbound)
	} else {
		out =  terminal.GetColoredText(terminal.GreenColor," Done")
	}

	return out
}

// SetStyle will set the style of bar, it's a coroutine safe function
func (b *Bar) SetStyle(style []rune) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if style == nil {
		style = []rune(DefaultStyle)
	}
	if len(style) != 5 {
		return fmt.Errorf("Invalid style")
	}
	b.lbound = style[0]
	b.fill = style[1]
	b.tip = style[2]
	b.space = style[3]
	b.rbound = style[4]
	return nil
}

// NewBar creates new bar, the style must have five utf8 char
func NewBar(total int64, style []rune) *Bar {
	if style == nil || len(style) != 5 {
		style = []rune(DefaultStyle)
	}
	return &Bar{
		lbound: style[0],
		fill:   style[1],
		tip:    style[2],
		space:  style[3],
		rbound: style[4],
		total:  total,
	}
}
