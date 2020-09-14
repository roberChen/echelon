package renderers

import (
	"bufio"
	"github.com/roberChen/echelon"
	"github.com/roberChen/echelon/renderers/config"
	"github.com/roberChen/echelon/renderers/internal/console"
	"github.com/roberChen/echelon/renderers/internal/node"
	"github.com/roberChen/echelon/terminal"
	"os"
	"sync"
	"time"
)

// resetAutoWrap, \u001B(16) == \033(10) == \x1b(16) == ESC key.
//
// ESC[={value}l : Resets the mode by using the same values that 
// Set Mode uses, except for 7, which disables line wrapping. The 
// last character in this escape sequence is a lowercase L.
const resetAutoWrap = "\u001B[?7l"
const defaultFrameBufSize = 38400 // 80 by 120 of 4 bytes UTF-8 characters

// InteractiveRenderer is a interactive rendere which is a dynamic one. 
//
// It conains a bufio.Writer for output, a root node and terminal height
type InteractiveRenderer struct {
	out               *bufio.Writer
	rootNode          *node.EchelonNode
	config            *config.InteractiveRendererConfig
	currentFrameLines []string
	drawLock          sync.Mutex
	terminalHeight    int
}

// NewInteractiveRenderer creates a new InteractiveRenderer
func NewInteractiveRenderer(out *os.File, rendererConfig *config.InteractiveRendererConfig) *InteractiveRenderer {
	if rendererConfig == nil {
		rendererConfig = config.NewDefaultRenderingConfig()
	}
	return &InteractiveRenderer{
		out:            bufio.NewWriterSize(out, defaultFrameBufSize),
		rootNode:       node.NewEchelonNode("root", rendererConfig),
		config:         rendererConfig,
		terminalHeight: console.TerminalHeight(out),
	}
}

// findScopeNode will return node with path 'scopes' in InteractiveRenderer, if the
// node does not exist, it will create one
func findScopedNode(scopes []string, r *InteractiveRenderer) *node.EchelonNode {
	result := r.rootNode
	for _, scope := range scopes {
		result = result.FindOrCreateChild(scope)
	}
	return result
}

// RenderScopeStarted starts render the node specified by the entry
func (r *InteractiveRenderer) RenderScopeStarted(entry *echelon.LogScopeStarted) {
	findScopedNode(entry.GetScopes(), r).Start()
}

// RenderScopeFinished will render an finished node specified by entry. 
//
// If the node is succeeded, all sub nodes (which must be succeeded as well) will hides. 
// If the node is failed, the node will keep showing at output with FailureColor(red)
func (r *InteractiveRenderer) RenderScopeFinished(entry *echelon.LogScopeFinished) {
	n := findScopedNode(entry.GetScopes(), r)
	if entry.Success() {
		if n != r.rootNode {
			n.ClearAllChildren()
			n.ClearDescription()
		}
		n.CompleteWithColor(r.config.SuccessStatus, r.config.Colors.SuccessColor)
	} else {
		n.SetVisibleDescriptionLines(r.config.DescriptionLinesWhenFailed)
		n.CompleteWithColor(r.config.FailureStatus, r.config.Colors.FailureColor)
	}
}

// RenderMessage will render message of node specified by entry, it will add the messages of
// entry to the node.
func (r *InteractiveRenderer) RenderMessage(entry *echelon.LogEntryMessage) {
	findScopedNode(entry.GetScopes(), r).AppendDescription(entry.GetMessage() + "\n")
}

// StartDrawing will start drawing Interactiverenderer until the root node has completed
//
// Each two frame has a time gap which can be configured in InteractiveRenderer.config.RefreshRate
func (r *InteractiveRenderer) StartDrawing() {
	_ = console.PrepareTerminalEnvironment()
	// don't wrap lines since it breaks incremental redraws
	_, _ = r.out.WriteString(resetAutoWrap)
	for !r.rootNode.HasCompleted() {
		r.DrawFrame()
		time.Sleep(r.config.RefreshRate)
	}
}

// StopDrawing will stop the InteractiveRenderer, it will complete the root node and draw final frame
func (r *InteractiveRenderer) StopDrawing() {
	r.rootNode.Complete()
	// one last redraw
	r.DrawFrame()
}

// DrawFrame will first generate full output lines and put it to terminal.
//
// This function is coroutine safe.
func (r *InteractiveRenderer) DrawFrame() {
	r.drawLock.Lock()
	defer r.drawLock.Unlock()
	var newFrameLines []string
	for _, n := range r.rootNode.GetChildren() {
		newFrameLines = append(newFrameLines, n.Render()...)
	}
	if r.terminalHeight > 0 {
		terminal.CalculateIncrementalUpdateMaxLines(r.out, r.currentFrameLines, newFrameLines, r.terminalHeight)
	} else {
		terminal.CalculateIncrementalUpdate(r.out, r.currentFrameLines, newFrameLines)
	}
	r.currentFrameLines = newFrameLines
}
