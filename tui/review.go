package tui

import (
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"

	"github.com/takacs/donkey/internal/card"
	"github.com/takacs/donkey/internal/review"
	"github.com/takacs/donkey/internal/supermemo"
)

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
		case key.Matches(msg, m.keys.Space):
			m.flip = true
		case key.Matches(msg, m.keys.Easy):
			m.handleGrade(review.Easy)
		case key.Matches(msg, m.keys.Good):
			m.handleGrade(review.Good)
		case key.Matches(msg, m.keys.Hard):
			m.handleGrade(review.Hard)
		case key.Matches(msg, m.keys.Again):
			m.handleGrade(review.Again)
		}

	}

	if m.currentCard >= m.numberOfCards {
		return InitProject(m.width, m.height)
	}

	return m, tea.Batch(cmds...)
}

func (m ReviewModel) View() string {
	helpView := m.help.View(m.keys)

	cardView := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Bold(true).
		Foreground(lipgloss.Color(primaryColor)).
		Render(m.cards[m.currentCard].Front) + "\n\n"
	if m.flip {
		cardViewBack := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Bold(false).
			Foreground(lipgloss.Color("231")).
			Render(m.cards[m.currentCard].Back)
		cardView += cardViewBack + "\n\n"
	} else {
		cardView += lipgloss.NewStyle().Align(lipgloss.Center).Render("\n\n")
	}
	cStr := strconv.Itoa(m.currentCard)
	nStr := strconv.Itoa(m.numberOfCards)

	counter := lipgloss.NewStyle().Align(lipgloss.Center).Foreground(lipgloss.Color(secondaryColor)).Render(cStr + " / " + nStr + "\n\n")

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		cardView+counter+helpView)
}

func newReviewModel(width, height, numberOfCards int) ReviewModel {
	supermemodb, err := supermemo.New()
	if err != nil {
		log.Fatal(err)
	}
	cardIds := supermemodb.GetXSoonestReviewTimeCardIds(numberOfCards)
	// TODO improve init
	carddb, err := card.New()
	if err != nil {
		log.Fatal(err)
	}

	cards, err := carddb.GetCardsFromIds(cardIds)
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

func (m ReviewModel) addReview(grade review.Grade) {
	reviewDb, err := review.New()
	if err != nil {
		log.Fatal(err)
	}
	err = reviewDb.Insert(m.cards[m.currentCard].ID, grade)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *ReviewModel) handleGrade(grade review.Grade) {
	m.addReview(grade)
	err := supermemo.UpdateCardParams(m.cards[m.currentCard].ID, grade)
	if err != nil {
		log.Fatal(err)
	}
	m.flip = false
	m.currentCard++
}
