package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	MainMenu key.Binding
	Exit     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.MainMenu, k.Exit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.MainMenu}, // first column
		{k.Exit},     // first column
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
}
