package internal

import "charm.land/lipgloss/v2"

var OwnMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#00C68D"))

var OtherMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#0055DA"))

var SystemMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#e80000")).
	Italic(true)
var PrivateMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#fcff33")).
	Italic(true)
