package main

import (
	"strings"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	messages []string
	textarea textarea.Model
	viewport viewport.Model
}

func initialModel() model {
	ta := textarea.New()
	vp := viewport.New(viewport.WithWidth(60), viewport.WithHeight(10))
	ta.Placeholder = "Type a message..."
	ta.SetWidth(60)
	ta.SetHeight(10)
	ta.SetVirtualCursor(false)
	ta.Focus()
	return model{textarea: ta, viewport: vp}

}

func (m model) Init() tea.Cmd {

	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		default:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			m.viewport, _ = m.viewport.Update(msg)
			return m, cmd
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			m.messages = append(m.messages, m.textarea.Value())
			var messages strings.Builder
			for _, message := range m.messages {
				messages.WriteString(message)
				messages.WriteString("\n")
			}
			m.viewport.SetContent(messages.String())
			m.viewport.GotoBottom()
			m.textarea.Reset()
			return m, nil
		}
	}
	return m, nil

}
func (m model) View() tea.View {

	return tea.NewView("LAN Chat\n\n" + m.textarea.View() + "\nmessages: " + m.viewport.View())
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
