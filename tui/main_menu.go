package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle()

type Model struct {
	width, height int
	keys          keyMap
	help          help.Model
	table         table.Model
}

func InitProject(width, height int) (tea.Model, tea.Cmd) {
	m := Model{
		width:  width,
		height: height,
		help:   help.New(),
		keys:   keys,
		table:  createTable(),
	}
	return m, func() tea.Msg { return "hi" }
}

func (m Model) Init() tea.Cmd {
	return nil
}

func createTable() table.Model {
	columns := []table.Column{
		{Title: "donkey", Width: 50},
	}

	rows := []table.Row{
		{"Add Card"},
		{"List Cards"},
		{"Review"},
		{"Stats"},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t

}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			cmd := tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()),
			)
			switch m.table.SelectedRow()[0] {
			case "Add Card":
				return newAddCardModel(m.width, m.height), cmd
			case "List Cards":
				return newListCardsModel(m.width, m.height), cmd
			case "Review":
				return newReviewModel(m.width, m.height, 20), cmd
			case "Stats":
				return newStatsModel(m.width, m.height), cmd
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	helpView := m.help.View(m.keys)

	style := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		baseStyle.Render(m.table.View())+"\n"+helpView)
	fmt.Printf("%v", style)
	return style
}
