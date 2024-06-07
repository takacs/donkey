package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/takacs/donkey/db"
)

type AddCardModel struct {
	keys keyMap
	help help.Model
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
		case key.Matches(msg, m.keys.Back):
			path, err := db.GetDbPath("cards")
			if err != nil {
				fmt.Println("error getting db path")
			}
			return InitProject(path)
		default:
			fmt.Printf("default press quit %v \n", msg)
			return m, tea.Quit
		}
	}
	return m, tea.Batch(cmds...)
}

func (m AddCardModel) View() string {
	helpView := m.help.View(m.keys)

	return m.name + "\n" + helpView
}

func newAddCardModel() AddCardModel {
	return AddCardModel{
		name: "play",
		help: help.New(),
		keys: keys,
	}
}
