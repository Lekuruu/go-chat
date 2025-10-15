package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatMessage struct {
	Timestamp time.Time
	Sender    string
	Content   string
	IsSystem  bool
}

type ChatUI struct {
	program  *tea.Program
	messages []ChatMessage
	users    []string
	mu       sync.Mutex
}

type model struct {
	viewport    viewport.Model
	textarea    textarea.Model
	messages    []ChatMessage
	users       []string
	ready       bool
	width       int
	height      int
	sendMessage func(string)
}

type newMessageMsg ChatMessage
type usersUpdateMsg []string

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("240"))

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	senderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86"))

	systemStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("243"))

	timestampStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	userListStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(lipgloss.Color("240")).
			PaddingLeft(1)

	inputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(lipgloss.Color("240"))
)

func NewChatUI(sendMessage func(string)) *ChatUI {
	ui := &ChatUI{
		messages: make([]ChatMessage, 0),
		users:    make([]string, 0),
	}

	m := model{
		messages:    ui.messages,
		users:       ui.users,
		sendMessage: sendMessage,
	}

	m.textarea = textarea.New()
	m.textarea.Placeholder = "Type a message..."
	m.textarea.Focus()
	m.textarea.Prompt = "> "
	m.textarea.CharLimit = 500
	m.textarea.SetHeight(3)
	m.textarea.ShowLineNumbers = false

	ui.program = tea.NewProgram(m, tea.WithAltScreen())
	return ui
}

func (ui *ChatUI) Run() error {
	_, err := ui.program.Run()
	return err
}

func (ui *ChatUI) AddMessage(sender, content string) {
	ui.mu.Lock()
	msg := ChatMessage{
		Timestamp: time.Now(),
		Sender:    sender,
		Content:   content,
		IsSystem:  false,
	}
	ui.messages = append(ui.messages, msg)
	ui.mu.Unlock()

	if ui.program != nil {
		ui.program.Send(newMessageMsg(msg))
	}
}

func (ui *ChatUI) AddSystemMessage(format string, args ...interface{}) {
	ui.mu.Lock()
	msg := ChatMessage{
		Timestamp: time.Now(),
		Content:   fmt.Sprintf(format, args...),
		IsSystem:  true,
	}
	ui.messages = append(ui.messages, msg)
	ui.mu.Unlock()

	if ui.program != nil {
		ui.program.Send(newMessageMsg(msg))
	}
}

func (ui *ChatUI) SetUsers(users []string) {
	ui.mu.Lock()
	ui.users = users
	ui.mu.Unlock()

	if ui.program != nil {
		ui.program.Send(usersUpdateMsg(users))
	}
}

func (ui *ChatUI) AddUser(user string) {
	ui.mu.Lock()
	ui.users = append(ui.users, user)
	users := make([]string, len(ui.users))
	copy(users, ui.users)
	ui.mu.Unlock()

	if ui.program != nil {
		ui.program.Send(usersUpdateMsg(users))
	}
}

func (ui *ChatUI) RemoveUser(user string) {
	ui.mu.Lock()
	for i, u := range ui.users {
		if u == user {
			ui.users = append(ui.users[:i], ui.users[i+1:]...)
			break
		}
	}
	users := make([]string, len(ui.users))
	copy(users, ui.users)
	ui.mu.Unlock()

	if ui.program != nil {
		ui.program.Send(usersUpdateMsg(users))
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			// Exit the program
			return m, tea.Quit
		case tea.KeyEnter:
			// We want to send a message now
			// -> grab the content and clear the textarea
			content := strings.TrimSpace(m.textarea.Value())
			if content != "" && m.sendMessage != nil {
				m.sendMessage(content)
				m.textarea.Reset()
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		// Handle console window resizing
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(msg.Width-22, msg.Height-6)
			m.viewport.YPosition = 2
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 22
			m.viewport.Height = msg.Height - 6
		}

		m.textarea.SetWidth(msg.Width - 22)
		m.viewport.SetContent(m.renderMessages())

	case newMessageMsg:
		m.messages = append(m.messages, ChatMessage(msg))
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()

	case usersUpdateMsg:
		m.users = []string(msg)
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}

	header := headerStyle.Width(m.width - 22).Render("go-chat")
	userListHeader := headerStyle.Width(20).Render(fmt.Sprintf("Users (%d)", len(m.users)))

	chatArea := m.viewport.View()
	userList := m.renderUserList()

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		chatArea,
		userList,
	)

	input := inputStyle.Width(m.width - 22).Render(m.textarea.View())
	fullHeader := lipgloss.JoinHorizontal(lipgloss.Top, header, userListHeader)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		fullHeader,
		mainContent,
		input,
	)
}

func (m model) renderMessages() string {
	var lines []string

	for _, msg := range m.messages {
		timestamp := timestampStyle.Render(msg.Timestamp.Format("15:04:05"))

		if msg.IsSystem {
			line := fmt.Sprintf("%s %s",
				timestamp,
				systemStyle.Render("* "+msg.Content),
			)
			lines = append(lines, line)
		} else {
			line := fmt.Sprintf("%s %s %s",
				timestamp,
				senderStyle.Render(msg.Sender+":"),
				messageStyle.Render(msg.Content),
			)
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

func (m model) renderUserList() string {
	var lines []string
	lines = append(lines, "")

	for _, user := range m.users {
		lines = append(lines, "• "+user)
	}

	content := strings.Join(lines, "\n")
	return userListStyle.
		Width(20).
		Height(m.viewport.Height).
		Render(content)
}
