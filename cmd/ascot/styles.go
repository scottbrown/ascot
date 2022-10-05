// A common set of CLI styles
package main

import (
	"github.com/charmbracelet/lipgloss"
)

// A heading that stands out among the other text
var headingStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.AdaptiveColor{Light: "12", Dark: "86"})

// A very bright, loud display of text
var alertStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#FF0000"))

// styling for text that show success
var passStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

// styling for text that shows failure
var failStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

// styling for highlighting text
var highlightStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "12", Dark: "86"})
