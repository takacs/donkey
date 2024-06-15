package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	MainMenu key.Binding
	Exit     key.Binding
	Tab      key.Binding
	Enter    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.MainMenu, k.Exit, k.Tab, k.Enter}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.MainMenu}, // first column
		{k.Exit},
		{k.Tab, k.Enter},
	}
}

var keys = keyMap{
	MainMenu: key.NewBinding(
		key.WithKeys("m", "M"),
		key.WithHelp("m/M", "main menu"),
	),
	Exit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "escape donkey"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
}
