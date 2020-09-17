package config

import (
	"github.com/roberChen/echelon/terminal"
	"runtime"
	"time"
)

// InteractiveRendererConfig is a structure which defines config of interactive renderer
type InteractiveRendererConfig struct {
	Colors                         *terminal.ColorSchema
	RefreshRate                    time.Duration
	ProgressIndicatorFrames        []string
	ProgressIndicatorCycleDuration time.Duration
	SuccessStatus                  string
	FailureStatus                  string
	DescriptionLinesWhenFailed     int
}

// NewDefaultRenderingConfig returns default config for current system
func NewDefaultRenderingConfig() *InteractiveRendererConfig {
	if runtime.GOOS == "windows" {
		return NewDefaultWindowsRenderingConfig()
	}
	return NewDefaultUnixRenderingConfig()
}

// NewDefaultUnixRenderingConfig returns default config for *nix
func NewDefaultUnixRenderingConfig() *InteractiveRendererConfig {
	//nolint:gomnd
	return &InteractiveRendererConfig{
		Colors:      terminal.DefaultColorSchema(),
		RefreshRate: 200 * time.Microsecond,
		ProgressIndicatorFrames: []string{
			"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛",
		},
		ProgressIndicatorCycleDuration: time.Second,
		SuccessStatus:                  "✅",
		FailureStatus:                  "❌",
		DescriptionLinesWhenFailed:     100,
	}
}

// NewDefaultWindowsRenderingConfig returns default config for Windows
func NewDefaultWindowsRenderingConfig() *InteractiveRendererConfig {
	//nolint:gomnd
	return &InteractiveRendererConfig{
		Colors:      terminal.DefaultColorSchema(),
		RefreshRate: 250 * time.Microsecond,
		ProgressIndicatorFrames: []string{
			"\\", "|", "/", "-",
		},
		ProgressIndicatorCycleDuration: time.Second,
		SuccessStatus:                  "+",
		FailureStatus:                  "-",
		DescriptionLinesWhenFailed:     100,
	}
}

// CurrentProgressIndicatorFrame returns current status
func (config *InteractiveRendererConfig) CurrentProgressIndicatorFrame() string {
	amountOfFrames := int64(len(config.ProgressIndicatorFrames))
	nanosPerFrame := int64(config.ProgressIndicatorCycleDuration) / amountOfFrames
	currentNanosTail := time.Now().UnixNano() % int64(config.ProgressIndicatorCycleDuration)
	frameIndex := currentNanosTail / nanosPerFrame
	if frameIndex < amountOfFrames {
		return config.ProgressIndicatorFrames[frameIndex]
	}
	return config.ProgressIndicatorFrames[0]
}
