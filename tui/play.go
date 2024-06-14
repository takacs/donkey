package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/takacs/donkey/db"
	"golang.org/x/term"
	"os"
)

type keyMap struct {
	Back key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back}, // first column
	}
}

var keys = keyMap{
	Back: key.NewBinding(
		key.WithKeys("B", "b"),
		key.WithHelp("b/B", "go back to main menu"),
	),
}

type PlayModel struct {
	keys keyMap
	help help.Model
	name string
}

func (m PlayModel) Init() tea.Cmd {
	return nil
}

func (m PlayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			path, err := db.GetDbPath("cards")
			if err != nil {
				fmt.Println("error getting db path")
			}
			termWidth, termHeight, _ := term.GetSize(int(os.Stdin.Fd()))
			return InitProject(path, termWidth, termHeight)
		default:
			fmt.Printf("default press quit %v \n", msg)
			return m, tea.Quit
		}
	}
	return m, tea.Batch(cmds...)
}

func (m PlayModel) View() string {
	helpView := m.help.View(m.keys)

	return m.name + "\n" + helpView
}

func newPlayModel() PlayModel {
	return PlayModel{
		name: "play",
		help: help.New(),
		keys: keys,
	}
}
