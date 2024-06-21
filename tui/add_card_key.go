package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type addCardKeyMap struct {
	MainMenu key.Binding
	Next     key.Binding
	Submit   key.Binding
}

func (k addCardKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Submit}
}

func (k addCardKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next},
		{k.Submit},
	}
}

var addCardKeys = addCardKeyMap{
	MainMenu: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "escape to main menu"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "add card"),
	),
}
