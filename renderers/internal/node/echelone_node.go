package node

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/roberChen/echelon/renderers/config"
	"github.com/roberChen/echelon/terminal"
	"github.com/roberChen/echelon/utils"
	"golang.org/x/text/width"
)

const defaultVisibleLines = 5

// EchelonNode is a log node for interactive renderer, it is designed for coroutine safe object
//
// It has title with specific title color. it's max visible lines can be specified.
// A node has children nodes.
type EchelonNode struct {
	lock                    sync.RWMutex
	done                    sync.WaitGroup
	status                  string
	title                   string
	titleColor              int
	description             []string
	visibleDescriptionLines int
	config                  *config.InteractiveRendererConfig
	startTime               time.Time
	endTime                 time.Time
	children                []*EchelonNode

	// bar setting
	Pbar *Bar
	// terminal width
	width int
}

// StartNewEchelonNode will create new EchelonNode with title and configuration, and start it.
func StartNewEchelonNode(title string, width int, total int64, config *config.InteractiveRendererConfig) *EchelonNode {
	result := NewEchelonNode(title, width, config)
	result.Start(total)
	return result
}

// NewEchelonNode will create new EchelonNode, and set startTime the function calling time
func NewEchelonNode(title string, width int, config *config.InteractiveRendererConfig) *EchelonNode {
	zeroTime := time.Time{}
	result := &EchelonNode{
		// the default status is pause status
		status:     "⏸",
		title:      title,
		titleColor: config.Colors.NeutralColor,
		// description is the texts will be diplayed to output
		description:             make([]string, 0),
		visibleDescriptionLines: defaultVisibleLines,
		config:                  config,
		startTime:               zeroTime,
		endTime:                 zeroTime,
		children:                make([]*EchelonNode, 0),

		width: width,
	}
	result.done.Add(1)
	return result
}

// GetChildren will returns all childrens of node, it's a coroutine safe function
func (node *EchelonNode) GetChildren() []*EchelonNode {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return node.children
}

// UpdateTitle will update title of node, it's a coroutine safe function
func (node *EchelonNode) UpdateTitle(text string) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.title = text
}

// UpdateConfig will update configuration of node, it's a coroutine safe function
func (node *EchelonNode) UpdateConfig(config *config.InteractiveRendererConfig) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.config = config
}

// ClearAllChildren will remove all children nodes of current node, it's a coroutine
// safe function
func (node *EchelonNode) ClearAllChildren() {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.children = make([]*EchelonNode, 0)
}

// ClearDescription will clean description of node, it will set description
// as a empty string list. it's a coroutine safe function
func (node *EchelonNode) ClearDescription() {
	node.SetDescription(make([]string, 0))
}

// SetDescription will set description of node, it's a coroutine safe function.
func (node *EchelonNode) SetDescription(description []string) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.description = description
}

// SetVisibleDescriptionLines will set max line number allowed to display for node.
// It's a coroutine safe function
func (node *EchelonNode) SetVisibleDescriptionLines(count int) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.visibleDescriptionLines = count
}

// DescriptionLength returns the length of description of node. It's a coroutine safe function
func (node *EchelonNode) DescriptionLength() int {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return len(node.description)
}

// Render function will output the rendered text of a node, with it's sub nodes.
//
// If the sub nodes lines are greater than max limitation, it will use '...' to
// replace some head lines.
func (node *EchelonNode) Render() []string {
	return node.render(0)
}

// returns strings with start indent space
func (node *EchelonNode) render(indent int) []string {
	title := node.fancyTitle(indent)
	newindent := indent + 2 // two spaces by default
	props, _ := width.LookupString(title)
	if props.Kind() == width.EastAsianWide || props.Kind() == width.EastAsianFullwidth {
		newindent++ // three spaces since title start with a wide emoji, thus indent expands 3 in total
	}
	result := []string{title}
	result = append(result, node.renderChildren(newindent)...)
	node.lock.RLock()
	defer node.lock.RUnlock()
	// add indent for descriptions
	sindent := strings.Repeat(" ", newindent)
	if len(node.description) > node.visibleDescriptionLines && node.visibleDescriptionLines >= 0 {
		result = append(result, sindent+"...")
		for _, line := range node.description[(len(node.description) - node.visibleDescriptionLines):] {
			result = append(result, sindent+line)
		}
	} else {
		for _, line := range node.description {
			result = append(result, sindent+line)
		}
	}

	return result
}

// renderChildren will render all childs nodes with specific indent and return the line strings.
func (node *EchelonNode) renderChildren(indent int) []string {
	node.lock.RLock()
	defer node.lock.RUnlock()
	var result []string
	for _, child := range node.children {
		result = append(result, child.render(indent)...)
	}
	return result
}

// fancyTitle will decorate title with specific indent. It renders the node it self.
// It's a coroutine safe function
//
// It will decorate with prefix which shows the status of node(stop, running, succeeded, failed),
// the colored title and spent time.
//
// If the node has children, it won't show the decimal of time. The structure will be like:
//
// indent space+prefix+title+spendtime
func (node *EchelonNode) fancyTitle(indent int) string {
	duration := utils.FormatDuration(node.ExecutionDuration(), len(node.children) == 0)
	isRunning := node.IsRunning()

	node.lock.RLock()
	defer node.lock.RUnlock()
	prefix := node.status
	if isRunning {
		prefix = node.config.CurrentProgressIndicatorFrame()
	}
	coloredTitle := node.title
	if node.titleColor >= 0 {
		coloredTitle = terminal.GetColoredText(node.titleColor, node.title)
	}
	out := strings.Repeat(" ", indent) + fmt.Sprintf("%s %s %s", prefix, coloredTitle, duration)
	// progress bar rendering
	if node.Pbar != nil {
		outwidth := runewidth.StringWidth(out)
		out = out + node.Pbar.String(node.width - outwidth)
	}
	return out
}

// ExecutionDuration will returns the spent time of node, it's a coroutine safe function
//
// If the node is still in progress, it returns the passed time since the start of a node,
// else it returns the total time of a node spent.
func (node *EchelonNode) ExecutionDuration() time.Duration {
	node.lock.RLock()
	defer node.lock.RUnlock()
	if !node.startTime.IsZero() && node.endTime.IsZero() {
		return time.Since(node.startTime)
	}
	return node.endTime.Sub(node.startTime)
}

// HasStarted returns wheter a node has started, a finished node is also started. This
// is a coroutine safe function
func (node *EchelonNode) HasStarted() bool {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return !node.startTime.IsZero()
}

// HasCompleted returns wheter a node  has completed, it's a coroutine safe function
func (node *EchelonNode) HasCompleted() bool {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return !node.endTime.IsZero()
}

// IsRunning returns wheter a node is running, it's a coroutine safe function
func (node *EchelonNode) IsRunning() bool {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return !node.startTime.IsZero() && node.endTime.IsZero()
}

// StartNewChild will create a child node with current node configuration for node
func (node *EchelonNode) StartNewChild(childName string) *EchelonNode {
	child := StartNewEchelonNode(childName, node.width, 0, node.config)
	node.AddNewChild(child)
	return child
}

// FindOrCreateChild will return the last child node with specific node title
// if the node doesn't exist, it will create one. It's a coroutine safe function
func (node *EchelonNode) FindOrCreateChild(childTitle string) *EchelonNode {
	node.lock.Lock()
	defer node.lock.Unlock()
	// look from the end since this is a common pattern to get the last child
	for i := len(node.children) - 1; i >= 0; i-- {
		child := node.children[i]
		if child.title == childTitle {
			return child
		}
	}
	child := NewEchelonNode(childTitle, node.width,node.config)
	node.children = append(node.children, child)
	return child
}

// AddNewChild will add child node for node. It's a coroutine safe function
func (node *EchelonNode) AddNewChild(child *EchelonNode) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.children = append(node.children, child)
}

// Start will start node, it sets the start time. It's a coroutine safe function
func (node *EchelonNode) Start(total int64) {
	node.lock.Lock()
	defer node.lock.Unlock()
	if node.startTime.IsZero() {
		node.startTime = time.Now()
	}
	if total != 0 {
		node.Pbar = NewBar(total, nil )
	}
}

// CompleteWithColor will stop a node with specific status and color. It's a coroutine
// safe function
func (node *EchelonNode) CompleteWithColor(status string, titleColor int) {
	if !node.endTime.IsZero() {
		return
	}
	node.lock.Lock()
	defer node.lock.Unlock()
	node.endTime = time.Now()
	if node.startTime.IsZero() {
		node.startTime = node.endTime
	}
	node.status = status
	node.titleColor = titleColor
	node.done.Done()
}

// Complete will stop a node. It's a coroutine safe function
func (node *EchelonNode) Complete() {
	if !node.endTime.IsZero() {
		return
	}
	node.lock.Lock()
	defer node.lock.Unlock()
	node.endTime = time.Now()
	if node.startTime.IsZero() {
		node.startTime = node.endTime
	}
	node.done.Done()
}

// SetTitleColor will set color of node title, it's a coroutine safe function
func (node *EchelonNode) SetTitleColor(ansiColor int) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.titleColor = ansiColor
}

// SetStatus will set status off node, it's a coroutine safe function
func (node *EchelonNode) SetStatus(text string) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.status = text
}

// WaitCompletion will wait the node to complete, it won't wait child nodes
func (node *EchelonNode) WaitCompletion() {
	node.done.Wait()
}

// AppendDescription will add text which might be multilines to node description, it
// won't start new line at the end of description. It's a coroutine safe function.
func (node *EchelonNode) AppendDescription(text string) {
	if node.HasCompleted() {
		return
	}
	node.lock.Lock()
	defer node.lock.Unlock()
	linesToAppend := strings.Split(text, "\n")
	if len(linesToAppend) == 0 {
		return
	}
	if len(node.description) == 0 {
		node.description = linesToAppend
		return
	}
	// append first new line to the last one
	node.description[len(node.description)-1] = node.description[len(node.description)-1] + linesToAppend[0]
	if len(linesToAppend) > 1 {
		node.description = append(node.description, linesToAppend[1:]...)
	}
}
