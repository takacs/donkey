package tui

import (
	"fmt"

	"errors"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/takacs/donkey/internal/card"
)

type ListCardsModel struct {
	width, height int
	keys          keyMap
	help          help.Model
	table         table.Model
	name          string
}

func (m ListCardsModel) Init() tea.Cmd {
	return nil
}

func (m ListCardsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Exit):
			return InitProject(m.width, m.height)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ListCardsModel) View() string {
	helpView := m.help.View(m.keys)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		baseStyle.Render(m.table.View())+"\n"+helpView)
}

func newListCardsModel(width, height int) ListCardsModel {
	table, err := getTableFromCards(width, height)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	return ListCardsModel{
		width:  width,
		height: height,
		name:   "list_cards",
		help:   help.New(),
		keys:   keys,
		table:  table,
	}
}

func getTableFromCards(width, height int) (table.Model, error) {
	carddb, err := card.New()
	if err != nil {
		return table.Model{}, errors.New("error getting db")
	}
	cards, err := carddb.GetXCards(0)
	if err != nil {
		return table.Model{}, errors.New("error getting cards")
	}

	// table setup
	columns := []table.Column{
		{Title: "Front", Width: int(float32(width) * 0.425)},
		{Title: "Back", Width: int(float32(width) * 0.425)},
		{Title: "Deck", Width: int(float32(width) * 0.05)},
	}

	rows := []table.Row{}
	for _, card := range cards {
		rows = append(rows, table.Row{card.Front, card.Back, card.Deck})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(int(float32(height)*0.8)),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		Bold(true).
		PaddingBottom(1).
		Foreground(lipgloss.Color(primaryColor))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color(secondaryColor)).
		Bold(false)
	t.SetStyles(s)

	return t, nil
}
