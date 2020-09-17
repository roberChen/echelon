package renderers

import (
	"fmt"
	"github.com/roberChen/echelon"
	"github.com/roberChen/echelon/renderers/internal/console"
	"github.com/roberChen/echelon/terminal"
	"github.com/roberChen/echelon/utils"
	"io"
	"strings"
	"time"
)

// SimpleRenderer is a simple renderer with an io.Writer for output, a color for output
// color, and a map to save time  stamps. The key of time stamps is the path of scope.
type SimpleRenderer struct {
	out        io.Writer
	colors     *terminal.ColorSchema
	startTimes map[string]time.Time
}

// NewSimpleRenderer creates a simple renderer
func NewSimpleRenderer(out io.Writer, colors *terminal.ColorSchema) *SimpleRenderer {
	if colors == nil {
		colors = terminal.DefaultColorSchema()
	}
	_ = console.PrepareTerminalEnvironment()
	return &SimpleRenderer{
		out:        out,
		colors:     colors,
		startTimes: make(map[string]time.Time),
	}
}
// RenderScopeStarted function of SimpleRenderer, it will start rendering an message of entry.
func (r SimpleRenderer) RenderScopeStarted(entry *echelon.LogScopeStarted) {
	scopes := entry.GetScopes()
	level := len(scopes)
	if level == 0 {
		return
	}
	timeKey := strings.Join(scopes, "/")
	if _, ok := r.startTimes[timeKey]; ok {
		// duplicate event
		return
	}
	r.startTimes[timeKey] = time.Now()
	lastScope := scopes[level-1]
	message := terminal.GetColoredText(r.colors.NeutralColor, fmt.Sprintf("Started %s", quotedIfNeeded(lastScope)))
	r.renderEntry(message)
}

// RenderScopeFinished will render a finished entry, which will print to task result of an entry.
func (r SimpleRenderer) RenderScopeFinished(entry *echelon.LogScopeFinished) {
	scopes := entry.GetScopes()
	level := len(scopes)
	if level == 0 {
		return
	}
	now := time.Now()
	startTime := now
	if t, ok := r.startTimes[strings.Join(scopes, "/")]; ok {
		startTime = t
	}
	duration := now.Sub(startTime)
	formatedDuration := utils.FormatDuration(duration, true)
	lastScope := scopes[level-1]
	if entry.Success() {
		message := fmt.Sprintf("%s succeeded in %s!", quotedIfNeeded(lastScope), formatedDuration)
		coloredMessage := terminal.GetColoredText(r.colors.SuccessColor, message)
		r.renderEntry(coloredMessage)
	} else {
		message := fmt.Sprintf("%s failed in %s!", quotedIfNeeded(lastScope), formatedDuration)
		coloredMessage := terminal.GetColoredText(r.colors.NeutralColor, message)
		r.renderEntry(coloredMessage)
	}
}

// RenderMessage will render message from entry for simple renderer, it sends message of 
// entry to renderEntry of renderer.
func (r SimpleRenderer) RenderMessage(entry *echelon.LogEntryMessage) {
	r.renderEntry(entry.GetMessage())
}

// RenderProcess function of SimpleRenderer, it will do nothing, simple renderer doesn't 
// support process rendering
func (r SimpleRenderer) RenderProcess(entry *echelon.LogProcessMessage) {}

// renderEntry will render message of simple renderer, it directly output the message to io.Writer of SimpleRenderer
func (r SimpleRenderer) renderEntry(message string) {
	_, _ = r.out.Write([]byte(message + "\n"))
}

// ScopeHasStarted returns whether the scope specified by path 'scpoes' has started. A finished scope is still 
// started.
func (r SimpleRenderer) ScopeHasStarted(scopes []string) bool {
	level := len(scopes)
	if level == 0 {
		return true
	}
	timeKey := strings.Join(scopes, "/")
	_, result := r.startTimes[timeKey]
	return result
}

// quotedIfNeeded will quotes string with ' if no ' or " appears in string
func quotedIfNeeded(s string) string {
	if strings.ContainsAny(s, "'\"") {
		return s
	}
	return "'" + s + "'"
}
