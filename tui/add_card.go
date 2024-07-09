package tui

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"

	"github.com/takacs/donkey/internal/card"
	"github.com/takacs/donkey/internal/supermemo"
)

const (
	front = iota
	back
	deck
)

const (
	formWidth = 50
)

var (
	inputStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(primaryColor))
	messageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(secondaryColor))
)

var MissingFieldError = errors.New("Insert failed, Front and Back are mandatory fields!")

type AddCardModel struct {
	width, height int
	inputs        []textinput.Model
	focus         int
	message       string
	keys          addCardKeyMap
	help          help.Model
}

func (m AddCardModel) Init() tea.Cmd {
	return nil
}

func (m AddCardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.MainMenu):
			return newMainMenuModel(m.width, m.height)
		case key.Matches(msg, m.keys.Next):
			m.nextFocus()
		case key.Matches(msg, m.keys.Submit):
			err := m.submitCard()
			if errors.Is(err, MissingFieldError) {
				m.message = "Insert failed, Front and Back are mandatory fields!"
				return m, tea.Batch(cmds...)
			}
			if err == nil {
				m.message = "Inserted!"
				return m, tea.Batch(cmds...)
			}
		}
	}
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m AddCardModel) View() string {
	formView := fmt.Sprintf(
		`
 %s

 %s


 %s

 %s


 %s

 %s


 %s

`,
		inputStyle.Width(formWidth).Render("Front"),
		m.inputs[front].View(),
		inputStyle.Width(formWidth).Render("Back"),
		m.inputs[back].View(),
		inputStyle.Width(formWidth).Render("Deck"),
		m.inputs[deck].View(),
		messageStyle.Render(m.message),
	)

	helpView := m.help.View(m.keys)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		baseStyle.Render(formView+"\n"+helpView))
}

func (m *AddCardModel) nextFocus() {
	m.inputs[m.focus].Blur()
	m.focus = (m.focus + 1) % 3
	m.inputs[m.focus].Focus()
}

func (m *AddCardModel) submitCard() error {
	if m.inputs[front].Value() == "" || m.inputs[back].Value() == "" {
		return MissingFieldError
	}
	carddb, err := card.New()
	if err != nil {
		return err
	}
	defer carddb.Close()
	cardId, err := carddb.Insert(m.inputs[front].Value(), m.inputs[back].Value(), m.inputs[deck].Value())
	if err != nil {
		return err
	}

	supermemodb, err := supermemo.New()
	if err != nil {
		return err
	}
	defer supermemodb.Close()
	err = supermemodb.Insert(cardId)
	if err != nil {
		return err
	}

	m.inputs = defaultInputs(0)
	return nil
}

func defaultInputs(focus int) []textinput.Model {
	var inputs []textinput.Model = make([]textinput.Model, 3)

	inputs[front] = textinput.New()
	inputs[front].Placeholder = lipgloss.NewStyle().PaddingRight(formWidth - len("question")).Render("question")
	inputs[front].PlaceholderStyle.AlignHorizontal(lipgloss.Left)
	inputs[front].CharLimit = 1000
	inputs[front].Width = formWidth

	inputs[back] = textinput.New()
	inputs[back].Placeholder = lipgloss.NewStyle().PaddingRight(formWidth - len("answer")).Render("answer")
	inputs[back].CharLimit = 1000
	inputs[back].Width = formWidth

	inputs[deck] = textinput.New()
	inputs[deck].Placeholder = lipgloss.NewStyle().PaddingRight(formWidth - len("default")).Render("default")
	inputs[deck].CharLimit = 100
	inputs[deck].Width = formWidth

	inputs[focus].Focus()

	return inputs

}

func newAddCardModel(width, height int) AddCardModel {
	focus := 0
	inputs := defaultInputs(focus)

	return AddCardModel{
		width:   width,
		height:  height,
		inputs:  inputs,
		focus:   focus,
		message: "",
		help:    help.New(),
		keys:    addCardKeys,
	}
}
