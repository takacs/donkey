package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	MainMenu key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.MainMenu}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.MainMenu}, // first column
	}
}

var keys = keyMap{
	MainMenu: key.NewBinding(
		key.WithKeys("m", "M"),
		key.WithHelp("m/M", "main menu"),
	),
}
