package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"

	"github.com/takacs/donkey/internal/card"
	"github.com/takacs/donkey/internal/review"
)

var reviewCardStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(10, 20, 10, 20)

type ReviewModel struct {
	width, height int
	keys          reviewKeyMap
	help          help.Model
	name          string
	cards         []card.Card
	currentCard   int
	flip          bool
	numberOfCards int
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
		case key.Matches(msg, m.keys.Enter):
			m.flip = true
		case key.Matches(msg, m.keys.Easy):
			m.addReview(review.Easy)
			m.flip = false
			m.currentCard += 1
		case key.Matches(msg, m.keys.Good):
			m.addReview(review.Good)
			m.flip = false
			m.currentCard += 1
		case key.Matches(msg, m.keys.Hard):
			m.addReview(review.Hard)
			m.flip = false
			m.currentCard += 1
		case key.Matches(msg, m.keys.Again):
			m.addReview(review.Again)
			m.flip = false
			m.currentCard += 1
		}

	}

	if m.currentCard >= m.numberOfCards {
		return InitProject(m.width, m.height)
	}
	return m, tea.Batch(cmds...)
}
func (m ReviewModel) View() string {
	helpView := m.help.View(m.keys)

	cardView := m.getCardView()

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		reviewCardStyle.Render(cardView)+"\n"+helpView)
}

func (m ReviewModel) getCardView() string {
	centerStyle := lipgloss.NewStyle().Align(lipgloss.Center)
	cardView := centerStyle.Render(m.cards[m.currentCard].Front) + "\n"
	if m.flip {
		cardView = cardView + centerStyle.Render("-----") + "\n"
		cardView = cardView + centerStyle.Render(m.cards[m.currentCard].Back)
	}
	return cardView
}

func newReviewModel(width, height, numberOfCards int) ReviewModel {
	// TODO improve init
	carddb, err := card.New()
	if err != nil {
		log.Fatal(err)
	}

	cards, err := carddb.GetCards(numberOfCards)
	if err != nil {
		log.Fatal(err)
	}

	h := help.New()
	h.ShowAll = true

	return ReviewModel{
		width:         width,
		height:        height,
		name:          "review",
		cards:         cards,
		numberOfCards: numberOfCards,
		help:          h,
		keys:          reviewKeys,
		flip:          false,
	}
}

func (m ReviewModel) addReview(grade review.Grade) error {
	reviewDb, err := review.New()
	if err != nil {
		log.Fatal(err)
	}
	err = reviewDb.Insert(m.cards[m.currentCard].ID, grade)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
