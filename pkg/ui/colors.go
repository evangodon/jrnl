package ui

import (
	lg "github.com/charmbracelet/lipgloss"
)

var Color = struct {
	Primary   lg.Color
	Secondary lg.Color
	White     lg.Color
	GreyLight lg.Color
}{
	Primary:   lg.Color("#a855f7"),
	Secondary: lg.Color("#e0e0e0"),
	White:     lg.Color("#ffffff"),
	GreyLight: lg.Color("#334155"),
}
