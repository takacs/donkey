package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	table table.Model
}

func InitProject(path string) (tea.Model, tea.Cmd) {
	m := Model{table: createTable()}
	return m, func() tea.Msg { return "hi" }
}

func (m Model) Init() tea.Cmd {
	return nil
}

func createTable() table.Model {
	columns := []table.Column{
		{Title: "Main Menu", Width: 20},
	}

	rows := []table.Row{
		{"Add Card"},
		{"List Cards"},
		{"Play"},
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
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			cmd := tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()),
			)
			switch m.table.SelectedRow()[0] {
			case "Add Card":
				return newAddCardModel(), cmd
			case "List Cards":
				return ListCardsModel{name: "list cards"}, cmd
			case "Play":
				return newPlayModel(), cmd
			case "Stats":
				return StatsModel{name: "stats"}, cmd
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}
