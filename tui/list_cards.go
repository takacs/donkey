package tui

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"

	"github.com/takacs/donkey/internal/card"
	"github.com/takacs/donkey/internal/review"
)

var layoutStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder())

type ListCardsModel struct {
	width, height int
	keys          listCardKeyMap
	help          help.Model
	table         table.Model
	cardInspect   bool
}

func (m ListCardsModel) Init() tea.Cmd {
	return nil
}

func (m ListCardsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.MainMenu):
			return newMainMenuModel(m.width, m.height)
		case key.Matches(msg, m.keys.Delete):
			err := m.deleteFocusedCard()
			if err != nil {
				log.Println(err)
				log.Println("delete failed")
			}
		case key.Matches(msg, m.keys.Inspect):
			m.cardInspect = !m.cardInspect
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ListCardsModel) View() string {
	helpView := m.help.View(m.keys)
	cardOverlay := ""

	bg := baseStyle.Render(m.table.View()) + "\n" + helpView
	if m.cardInspect {
		cardOverlay = PlaceOverlay(
			m.width/4, m.height/4,
			layoutStyle.
				Copy().
				Width(m.width/2).
				Height(m.height/2).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				BorderForeground(lipgloss.Color("#209fb5")).
				Render(
					m.table.SelectedRow()[1]+"\n\n"+m.table.SelectedRow()[2],
				),
			bg,
			false,
		)
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			cardOverlay,
		)
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		bg,
	)
}

func (m *ListCardsModel) deleteFocusedCard() error {
	cardId, err := strconv.Atoi(m.table.SelectedRow()[0])
	log.Printf("deleting %v", cardId)
	if err != nil {
		return err
	}
	carddb, err := card.New()
	if err != nil {
		return err
	}
	err = carddb.Delete(uint(cardId))
	if err != nil {
		return err
	}
	reviewdb, err := review.New()
	if err != nil {
		return err
	}
	err = reviewdb.Delete(uint(cardId))
	if err != nil {
		return err
	}

	cursor := m.table.Cursor()
	m.table, err = getTableFromCards(m.width, m.height)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	m.table.SetCursor(cursor)

	return nil
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
		{Title: "ID", Width: int(float32(width) * 0.025)},
		{Title: "Front", Width: int(float32(width) * 0.425)},
		{Title: "Back", Width: int(float32(width) * 0.425)},
		{Title: "Deck", Width: int(float32(width) * 0.05)},
	}

	rows := []table.Row{}
	for _, card := range cards {
		idstr := strconv.Itoa(int(card.ID))
		rows = append(rows, table.Row{idstr, card.Front, card.Back, card.Deck})
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

func newListCardsModel(width, height int) ListCardsModel {
	table, err := getTableFromCards(width, height)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	return ListCardsModel{
		width:  width,
		height: height,
		help:   help.New(),
		keys:   listCardKeys,
		table:  table,
	}
}
