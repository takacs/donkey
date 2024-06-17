package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type reviewKeyMap struct {
	MainMenu key.Binding
	Enter    key.Binding
	Again    key.Binding
	Hard     key.Binding
	Good     key.Binding
	Easy     key.Binding
}

func (k reviewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.MainMenu, k.Enter, k.Again, k.Hard, k.Good, k.Easy}
}

func (k reviewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Easy, k.Good, k.Hard, k.Again},
		{k.Enter, k.MainMenu},
	}
}

var reviewKeys = reviewKeyMap{
	MainMenu: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "escape to main menu"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("enter/space", "reveal back of card"),
	),
	Easy: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "Easy"),
	),
	Good: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "Good"),
	),
	Hard: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "Hard"),
	),
	Again: key.NewBinding(
		key.WithKeys("4"),
		key.WithHelp("4", "Again"),
	),
}
