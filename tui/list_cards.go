package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ListCardsModel struct {
	name string
}

func (m ListCardsModel) Init() tea.Cmd {
	return nil
}

func (m ListCardsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m ListCardsModel) View() string {
	return m.name
}
