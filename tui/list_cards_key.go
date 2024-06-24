package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type listCardKeyMap struct {
	MainMenu key.Binding
	Down     key.Binding
	Up       key.Binding
	Delete   key.Binding
	Inspect  key.Binding
}

func (k listCardKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Down, k.Up, k.Delete, k.Inspect}
}

func (k listCardKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Down},
		{k.Up},
		{k.Delete},
		{k.Inspect},
	}
}

var listCardKeys = listCardKeyMap{
	MainMenu: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "escape to main menu"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "down"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "up"),
	),
	Delete: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "delete"),
	),
	Inspect: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "toggle card inspect"),
	)}
