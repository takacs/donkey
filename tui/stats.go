package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type StatsModel struct {
	name string
}

func (m StatsModel) Init() tea.Cmd {
	return nil
}

func (m StatsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m StatsModel) View() string {
	return m.name
}
