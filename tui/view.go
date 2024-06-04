package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	text string
}

func InitProject() (tea.Model, tea.Cmd) {
	m := Model{text: "hello, tea"}
	return m, func() tea.Msg { return "hi" }
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.text
}
