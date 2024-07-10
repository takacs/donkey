package tui

import (
	"errors"
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
			return newMainMenuModel(m.width, m.height)
		case key.Matches(msg, m.keys.Space):
			m.flip = true
		case key.Matches(msg, m.keys.Easy):
			err := m.handleGrade(review.Easy)
			if err != nil {
				return m, tea.Batch(cmds...)
			}
		case key.Matches(msg, m.keys.Good):
			err := m.handleGrade(review.Good)
			if err != nil {
				return m, tea.Batch(cmds...)
			}
		case key.Matches(msg, m.keys.Hard):
			err := m.handleGrade(review.Hard)
			if err != nil {
				return m, tea.Batch(cmds...)
			}
		case key.Matches(msg, m.keys.Again):
			err := m.handleGrade(review.Again)
			if err != nil {
				return m, tea.Batch(cmds...)
			}
		}

	}

	if m.currentCard >= m.numberOfCards {
		return newMainMenuModel(m.width, m.height)
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

func newReviewModel(width, height, numberOfCards int) (ReviewModel, error) {
	supermemodb, err := supermemo.New()
	if err != nil {
		return ReviewModel{}, err
	}
	cardIds := supermemodb.GetXSoonestReviewTimeCardIds(numberOfCards)

	log.Println(cardIds)
	if len(cardIds) == 0 {
		log.Println("badbad")
		return ReviewModel{}, errors.New("no cards to review yet")
	} else {

		log.Println("goodgood")
	}
	// TODO improve init
	carddb, err := card.New()
	if err != nil {
		return ReviewModel{}, err
	}

	cards, err := carddb.GetCardsFromIds(cardIds)
	if err != nil {
		return ReviewModel{}, err
	}

	cardCount := numberOfCards
	if len(cards) < numberOfCards {
		cardCount = len(cards)
	}

	h := help.New()
	h.ShowAll = true

	return ReviewModel{
		width:         width,
		height:        height,
		name:          "review",
		cards:         cards,
		numberOfCards: cardCount,
		help:          h,
		keys:          reviewKeys,
		flip:          false,
	}, nil
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

func (m *ReviewModel) handleGrade(grade review.Grade) error {
	m.addReview(grade)
	err := supermemo.UpdateCardParams(m.cards[m.currentCard].ID, grade)
	if err != nil {
		return errors.New("can't update card params")
	}
	m.flip = false
	m.currentCard++
	return nil
}
