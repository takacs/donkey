package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	db "github.com/takacs/donkey/db"
)

type Model struct {
	cards []db.Card
}

func InitProject(path string) (tea.Model, tea.Cmd) {
	m := StatsModel{name: "stats"}
	return m, func() tea.Msg { return "hi" }
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		default:
			fmt.Printf("default press quit %v \n", msg)
			return m, tea.Quit
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.cards[0].Front
}
