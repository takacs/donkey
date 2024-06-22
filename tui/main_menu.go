package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle()
var primaryColor = "#f1b376"
var secondaryColor = "#abaf74"

type Model struct {
	width, height int
	errorMessage  string
	table         table.Model
}

func InitProject(width, height int) (tea.Model, tea.Cmd) {
	m := Model{
		width:  width,
		height: height,
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
		{"Review"},
		{"Add Card"},
		{"List Cards"},
		{"Settings"},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(12),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		Bold(true).
		PaddingBottom(2).
		Foreground(lipgloss.Color(primaryColor))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(secondaryColor)).
		Bold(false)
	s.Cell = s.Cell.Height(2).Align(lipgloss.Center).Padding(0, 1)
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
				model, err := newReviewModel(m.width, m.height, 20)
				if err != nil {
					m.errorMessage = "no cards to review yet!"
				} else {
					return model, cmd
				}
			case "Settings":
				return newStatsModel(m.width, m.height), cmd
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {

	style := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center-0.02,
		lipgloss.Center,
		baseStyle.Render(m.table.View())+"\n"+m.errorMessage)
	return style
}
