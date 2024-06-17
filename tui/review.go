package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"

	"github.com/takacs/donkey/internal/card"
)

type ReviewModel struct {
	width, height int
	keys          keyMap
	help          help.Model
	name          string
	cards         []card.Card
}

func (m ReviewModel) Init() tea.Cmd {
	return nil
}

func (m ReviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.MainMenu):
			return InitProject(m.width, m.height)
		default:
			return m, tea.Quit
		}
	}
	return m, tea.Batch(cmds...)
}
func (m ReviewModel) View() string {
	helpView := m.help.View(m.keys)

	log.Println(m.cards)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		baseStyle.Render(m.name+"\n"+helpView))
}

func newReviewModel(width, height int) ReviewModel {
	// TODO improve init
	carddb, err := card.New()
	if err != nil {
		log.Fatal(err)
	}

	cards, err := carddb.GetCards(20)
	if err != nil {
		log.Fatal(err)
	}

	return ReviewModel{
		width:  width,
		height: height,
		name:   "review",
		cards:  cards,
		help:   help.New(),
		keys:   keys,
	}
}
