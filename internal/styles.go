package internal

import "charm.land/lipgloss/v2"

var OwnMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("2"))

var OtherMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("6"))

var SystemMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("8")).
	Italic(true)
