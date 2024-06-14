package tui

import (
	"fmt"

	"errors"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/takacs/donkey/db"
	"golang.org/x/term"
	"os"
)

type ListCardsModel struct {
	keys       keyMap
	help       help.Model
	cardsTable table.Model
	name       string
}

func (m ListCardsModel) Init() tea.Cmd {
	return nil
}

func (m ListCardsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			path, err := db.GetDbPath("cards")
			if err != nil {
				fmt.Println("error getting db path")
			}
			termWidth, termHeight, _ := term.GetSize(int(os.Stdin.Fd()))
			return InitProject(path, termWidth, termHeight)
		default:
			fmt.Printf("default press quit %v \n", msg)
			return m, tea.Quit
		}
	}
	return m, tea.Batch(cmds...)
}

func (m ListCardsModel) View() string {
	helpView := m.help.View(m.keys)

	return baseStyle.Render(m.cardsTable.View()) + "\n" + helpView
}

func newListCardsModel() ListCardsModel {
	table, err := getTableFromCards()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	return ListCardsModel{
		name:       "list_cards",
		help:       help.New(),
		keys:       keys,
		cardsTable: table,
	}
}

func getTableFromCards() (table.Model, error) {
	path, err := db.GetDbPath("cards")
	if err != nil {
		return table.Model{}, errors.New("error getting db path")
	}
	cards_db, err := db.OpenDb(path)
	if err != nil {
		return table.Model{}, errors.New("error opening db path")
	}
	cards, err := cards_db.Getcards()
	if err != nil {
		return table.Model{}, errors.New("error getting cards")
	}

	// table setup
	columns := []table.Column{
		{Title: "Front", Width: 20},
		{Title: "Back", Width: 20},
		{Title: "Deck", Width: 20},
	}

	rows := []table.Row{}
	for _, card := range cards {
		rows = append(rows, table.Row{card.Front, card.Back, card.Deck})
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

	return t, nil
}
