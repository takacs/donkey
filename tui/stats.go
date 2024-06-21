package tui

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type StatsModel struct {
	width, height int
	help          help.Model
}

func (m StatsModel) Init() tea.Cmd {
	return nil
}

func (m StatsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m StatsModel) View() string {
	return ""
}

func newStatsModel(width, height int) StatsModel {
	return StatsModel{
		width:  width,
		height: height,
		help:   help.New(),
	}
}
