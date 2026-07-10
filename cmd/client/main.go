package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"lan-chat/internal"
	"lan-chat/protocol"
	"log"
	"net"
	"strings"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	messages   []string
	textarea   textarea.Model
	viewport   viewport.Model
	conn       net.Conn
	registered bool
	reader     *bufio.Reader
	username   string
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
	return model{
		textarea: ta,
		viewport: vp,
		conn:     conn,
		reader:   bufio.NewReader(conn),
	}

}

func (m model) Init() tea.Cmd {

	return tea.Batch(
		textarea.Blink,
		waitForMsg(m.reader),
	)
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
			m.conn.Close()
			return m, tea.Quit

		case "enter":
			text := strings.TrimSpace(m.textarea.Value())
			if text == "" {
				return m, nil
			}

			if text == "/who" {
				cmd := protocol.NewCommandMessage("who")

				data, err := json.Marshal(cmd)
				if err != nil {
					return m, nil
				}

				data = append(data, '\n')

				_, err = m.conn.Write(data)
				if err != nil {
					return m, tea.Quit
				}

				m.textarea.Reset()
				return m, nil
			}

			if strings.HasPrefix(text, "/msg ") {

				cmdContent := strings.TrimPrefix(text, "/")
				cmd := protocol.NewCommandMessage(cmdContent)

				data, err := json.Marshal(cmd)
				if err != nil {
					return m, nil
				}

				data = append(data, '\n')

				_, err = m.conn.Write(data)
				if err != nil {
					return m, tea.Quit
				}

				m.textarea.Reset()
				return m, nil
			}

			if !m.registered {
				m.username = text
				m.conn.Write([]byte(text + "\n"))
				m.textarea.Reset()
				return m, nil
			}

			msg := protocol.NewChatMessage(m.username, text)

			data, err := json.Marshal(msg)
			if err != nil {
				return m, nil
			}
			data = append(data, '\n')
			_, err = m.conn.Write(data)
			if err != nil {
				log.Printf("Error sending message: %v", err)
				return m, tea.Quit
			}

			m.textarea.Reset()

			return m, nil

		}
	case serverMsg:
		msgStr := string(msg)
		if msgStr == "REGISTERED_OK" {
			m.registered = true
			return m, waitForMsg(m.reader)
		}

		var wireMsg protocol.WireMessage
		err := json.Unmarshal([]byte(msgStr), &wireMsg)
		if err != nil {
			m.messages = append(m.messages, msgStr)
		} else {
			formatted := fmt.Sprintf("[%s]: %s", wireMsg.Sender, wireMsg.Content)
			switch {
			case wireMsg.Type == protocol.SystemMessage:
				formatted = internal.SystemMessageStyle.Render(formatted)
			case wireMsg.Type == protocol.PrivateMessage:
				formatted = internal.PrivateMessageStyle.Render(
					fmt.Sprintf("[PM from %s]: %s", wireMsg.Sender, wireMsg.Content),
				)
			case wireMsg.Sender == m.username:
				formatted = internal.OwnMessageStyle.Render(formatted)
			default:
				formatted = internal.OtherMessageStyle.Render(formatted)
			}

			m.messages = append(m.messages, formatted)
		}

		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()

		return m, waitForMsg(m.reader)

	case errMsg:
		log.Printf("Server disconnect: %v", msg)
		return m, tea.Quit
	}
	return m, nil

}
func (m model) View() tea.View {

	return tea.NewView("LAN Chat\n\n" + m.viewport.View() + "\nmessages: " + m.textarea.View())
}

type serverMsg string
type errMsg error

func waitForMsg(reader *bufio.Reader) tea.Cmd {
	return func() tea.Msg {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return errMsg(err)
		}
		return serverMsg(strings.TrimSuffix(msg, "\n"))
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
