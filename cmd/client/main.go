package main

import (
	"log"
	"net"
	"strings"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	messages []string
	textarea textarea.Model
	viewport viewport.Model
	conn     net.Conn
}

func initialModel() model {
	ta := textarea.New()
	vp := viewport.New(viewport.WithWidth(100), viewport.WithHeight(35))
	ta.Placeholder = "Type a message..."
	ta.SetWidth(60)
	ta.SetHeight(5)
	ta.SetVirtualCursor(false)
	ta.Focus()
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error: can not connect to the server", err)

	}
	return model{textarea: ta, viewport: vp, conn: conn}

}

func (m model) Init() tea.Cmd {

	return tea.Batch(
		textarea.Blink, waitForMsg(m.conn))
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
				//messages.WriteString("\n")
			}
			m.viewport.SetContent(messages.String())
			m.viewport.GotoBottom()
			m.conn.Write([]byte(m.textarea.Value()))
			m.textarea.Reset()
			return m, nil
		}
	case serverMsg:
		m.messages = append(m.messages, string(msg))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		return m, waitForMsg(m.conn)
	}
	return m, nil

}
func (m model) View() tea.View {

	return tea.NewView("LAN Chat\n\n" + m.viewport.View() + "\nmessages: " + m.textarea.View())
}

type serverMsg string

func waitForMsg(conn net.Conn) tea.Cmd {
	return func() tea.Msg {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error: Reading server message %s", err)
		}
		return serverMsg(string(buffer[:n]))
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
