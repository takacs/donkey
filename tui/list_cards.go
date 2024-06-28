package tui

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"

	"github.com/takacs/donkey/internal/card"
	"github.com/takacs/donkey/internal/review"
)

const (
	columnKeyId    = "id"
	columnKeyFront = "front"
	columnKeyBack  = "back"
	columnKeyDeck  = "deck"
)

var layoutStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder())

type ListCardsModel struct {
	width, height   int
	keys            listCardKeyMap
	help            help.Model
	table           table.Model
	filterTextInput textinput.Model
	cardInspect     bool
}

func (m ListCardsModel) Init() tea.Cmd {
	return nil
}

func (m ListCardsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.filterTextInput.Focused() {
			if msg.String() == "enter" {
				log.Printf("%v", msg)
				m.filterTextInput.Blur()
			} else if key.Matches(msg, m.keys.MainMenu) {
				return newMainMenuModel(m.width, m.height)
			} else {
				log.Printf("%v", msg)
				m.filterTextInput, _ = m.filterTextInput.Update(msg)
			}
			m.table = m.table.WithFilterInput(m.filterTextInput)

			return m, cmd
		}
		switch {
		case key.Matches(msg, m.keys.Search):
			m.filterTextInput.Focus()
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

	bg := "\n" + baseStyle.Render(m.filterTextInput.View()+"\n"+m.table.View()) + "\n" + helpView
	if m.cardInspect {
		rowData := m.table.HighlightedRow().Data
		front, exists := rowData[columnKeyFront]
		if !exists {
			front = "missing front"
		}
		back, exists := rowData[columnKeyBack]
		if !exists {
			back = "missing back"
		}
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
					front.(string)+"\n\n"+back.(string),
				),
			bg,
			false,
		)
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Top,
			cardOverlay,
		)
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Top,
		bg,
	)
}

func (m *ListCardsModel) deleteFocusedCard() error {
	rowData := m.table.HighlightedRow().Data
	id, exists := rowData[columnKeyId]
	if !exists {
		return errors.New("no id for row")
	}
	cardId, err := strconv.Atoi(id.(string))
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

	index := m.table.GetHighlightedRowIndex()
	page := m.table.CurrentPage()
	m.table, err = getTableFromCards(m.width, m.height)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	m.table = m.table.
		WithFilterInput(m.filterTextInput).
		WithCurrentPage(page).
		WithHighlightedRow(index)

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
	columns := []table.Column{
		table.NewColumn(columnKeyId, "ID", int(float32(width)*0.025)).WithFiltered(true).WithStyle(lipgloss.NewStyle().AlignHorizontal(lipgloss.Left)),
		table.NewColumn(columnKeyFront, "Front", int(float32(width)*0.425)).WithFiltered(true).WithStyle(lipgloss.NewStyle().AlignHorizontal(lipgloss.Left)),
		table.NewColumn(columnKeyBack, "Back", int(float32(width)*0.425)).WithFiltered(true).WithStyle(lipgloss.NewStyle().AlignHorizontal(lipgloss.Left)),
		table.NewColumn(columnKeyDeck, "Deck", int(float32(width)*0.05)).WithFiltered(true).WithStyle(lipgloss.NewStyle().AlignHorizontal(lipgloss.Left)),
	}

	rows := []table.Row{}
	for _, card := range cards {
		idstr := strconv.Itoa(int(card.ID))
		rows = append(rows,
			table.NewRow(table.RowData{
				columnKeyId:    idstr,
				columnKeyFront: card.Front,
				columnKeyBack:  card.Back,
				columnKeyDeck:  card.Deck,
			}),
		)
	}

	table := table.
		New(columns).
		Filtered(true).
		Focused(true).
		WithFooterVisibility(false).
		WithPageSize(height - 10).
		WithRows(rows).
		HighlightStyle(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color(secondaryColor)).
			Bold(false),
		).
		HeaderStyle(
			lipgloss.NewStyle().
				Bold(true).
				PaddingBottom(1).
				Foreground(lipgloss.Color(primaryColor)),
		)

	return table, nil
}

func newListCardsModel(width, height int) ListCardsModel {
	table, err := getTableFromCards(width, height)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	textinput := textinput.New()
	textinput.PromptStyle.AlignHorizontal(lipgloss.Left)
	textinput.Placeholder = "press \"s\" to search"

	return ListCardsModel{
		width:           width,
		height:          height,
		help:            help.New(),
		keys:            listCardKeys,
		table:           table,
		filterTextInput: textinput,
	}
}
