package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type AddCardModel struct {
	name string
}

func (m AddCardModel) Init() tea.Cmd {
	return nil
}

func (m AddCardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m AddCardModel) View() string {
	return m.name
}
