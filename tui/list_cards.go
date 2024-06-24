package tui

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	charmansi "github.com/charmbracelet/x/exp/term/ansi"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/truncate"
	"github.com/muesli/termenv"

	"github.com/takacs/donkey/internal/card"
	"github.com/takacs/donkey/internal/review"
)

type whitespace struct {
	chars string
	style termenv.Style
}

func (w whitespace) render(width int) string {
	if w.chars == "" {
		w.chars = " "
	}

	r := []rune(w.chars)
	j := 0
	b := strings.Builder{}

	// Cycle through runes and print them into the whitespace.
	for i := 0; i < width; {
		b.WriteRune(r[j])
		j++
		if j >= len(r) {
			j = 0
		}
		i += charmansi.StringWidth(string(r[j]))
	}

	// Fill any extra gaps white spaces. This might be necessary if any runes
	// are more than one cell wide, which could leave a one-rune gap.
	short := width - charmansi.StringWidth(b.String())
	if short > 0 {
		b.WriteString(strings.Repeat(" ", short))
	}

	return w.style.Styled(b.String())
}

type WhitespaceOption func(*whitespace)

var layoutStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	AlignVertical(lipgloss.Center)

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
			if m.cardInspect {
				m.cardInspect = false
			} else {
				return newMainMenuModel(m.width, m.height)
			}
		case key.Matches(msg, m.keys.Delete):
			err := m.deleteFocusedCard()
			if err != nil {
				log.Println(err)
				log.Println("delete failed")
			}
		case key.Matches(msg, m.keys.Inspect):
			m.cardInspect = true
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ListCardsModel) View() string {
	helpView := m.help.View(m.keys)
	cardOverlay := ""

	if m.cardInspect {
		cardOverlay = PlaceOverlay(
			m.width/10,
			5*m.height/100,
			layoutStyle.
				Copy().
				Width(8*m.width/10).
				Height(90*m.height/100).
				AlignHorizontal(lipgloss.Center).
				BorderForeground(lipgloss.Color("#209fb5")).
				Render(
					m.help.View(m.keys),
				),
			helpView,
		)
	}
	return lipgloss.JoinVertical(
		lipgloss.Center,
		baseStyle.Render(m.table.View()),
		helpView,
		cardOverlay,
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

func PlaceOverlay(x, y int, fg, bg string, opts ...WhitespaceOption) string {
	fgLines, fgWidth := getLines(fg)
	bgLines, bgWidth := getLines(bg)
	bgHeight := len(bgLines)
	fgHeight := len(fgLines)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		return fg
	}

	x = clamp(x, 0, bgWidth-fgWidth)
	y = clamp(y, 0, bgHeight-fgHeight)

	ws := &whitespace{}
	for _, opt := range opts {
		opt(ws)
	}

	var b strings.Builder
	for i, bgLine := range bgLines {
		if i > 0 {
			b.WriteByte('\n')
		}
		if i < y || i >= y+fgHeight {
			b.WriteString(bgLine)
			continue
		}

		pos := 0
		if x > 0 {
			left := truncate.String(bgLine, uint(x))
			pos = ansi.PrintableRuneWidth(left)
			b.WriteString(left)
			if pos < x {
				b.WriteString(ws.render(x - pos))
				pos = x
			}
		}

		fgLine := fgLines[i-y]
		b.WriteString(fgLine)
		pos += ansi.PrintableRuneWidth(fgLine)

		right := cutLeft(bgLine, pos)
		bgWidth := ansi.PrintableRuneWidth(bgLine)
		rightWidth := ansi.PrintableRuneWidth(right)
		if rightWidth <= bgWidth-pos {
			b.WriteString(ws.render(bgWidth - rightWidth - pos))
		}

		b.WriteString(right)
	}

	return b.String()
}

func clamp(v, lower, upper int) int {
	return min(max(v, lower), upper)
}

func getLines(s string) (lines []string, widest int) {
	lines = strings.Split(s, "\n")

	for _, l := range lines {
		w := charmansi.StringWidth(l)
		if widest < w {
			widest = w
		}
	}

	return lines, widest
}
func cutLeft(s string, cutWidth int) string {
	var (
		pos    int
		isAnsi bool
		ab     bytes.Buffer
		b      bytes.Buffer
	)
	for _, c := range s {
		var w int
		if c == ansi.Marker || isAnsi {
			isAnsi = true
			ab.WriteRune(c)
			if ansi.IsTerminator(c) {
				isAnsi = false
				if bytes.HasSuffix(ab.Bytes(), []byte("[0m")) {
					ab.Reset()
				}
			}
		} else {
			w = runewidth.RuneWidth(c)
		}

		if pos >= cutWidth {
			if b.Len() == 0 {
				if ab.Len() > 0 {
					b.Write(ab.Bytes())
				}
				if pos-cutWidth > 1 {
					b.WriteByte(' ')
					continue
				}
			}
			b.WriteRune(c)
		}
		pos += w
	}
	return b.String()
}
