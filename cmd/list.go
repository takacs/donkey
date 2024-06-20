package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/takacs/donkey/internal/card"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list cards in db",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		carddb, err := card.New()
		if err != nil {
			return err
		}
		defer carddb.Close()
		cards, err := carddb.GetCards(0)
		if err != nil {
			return err
		}
		fmt.Print(setupTable(cards))
		fmt.Print("\n")
		return nil
	},
}

func setupTable(cards []card.Card) *table.Table {
	columns := []string{"ID", "Front", "Back", "Deck", "Created At"}
	var rows [][]string
	for _, card := range cards {
		rows = append(rows, []string{
			fmt.Sprintf("%d", card.ID),
			card.Front,
			card.Back,
			card.Deck,
			card.Status,
			card.Created.Format("2006-01-02"),
		})
	}
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		Headers(columns...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("212")).
					Border(lipgloss.NormalBorder()).
					BorderTop(false).
					BorderLeft(false).
					BorderRight(false).
					BorderBottom(true).
					Bold(true)
			}
			if row%2 == 0 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
			}
			return lipgloss.NewStyle()
		})
	return t
}
