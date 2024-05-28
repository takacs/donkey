package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	ddb "github.com/takacs/donkey/db"
)

var rootCmd = &cobra.Command{
	Use:   "cards",
	Short: "A CLI card management tool for anki style brain training.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var addCmd = &cobra.Command{
	Use:   "add card",
	Short: "Add a new card with an optional deck name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := ddb.OpenDb(setupPath())
		if err != nil {
			return err
		}
		defer c.Db.Close()
		front, err := cmd.Flags().GetString("front")
		if err != nil {
			return err
		}
		back, err := cmd.Flags().GetString("back")
		if err != nil {
			return err
		}
		deck, err := cmd.Flags().GetString("deck")
		if err != nil {
			return err
		}
		if err := c.Insert(front, back, deck); err != nil {
			return err
		}
		return nil
	},
}

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show where your cards are stored",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Println(setupPath())
		return err
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete ID",
	Short: "Delete a card by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := ddb.OpenDb(setupPath())
		if err != nil {
			return err
		}
		defer c.Db.Close()
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		return c.Delete(uint(id))
	},
}

var updateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a card by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := ddb.OpenDb(setupPath())
		if err != nil {
			return err
		}
		defer c.Db.Close()
		front, err := cmd.Flags().GetString("front")
		if err != nil {
			return err
		}
		back, err := cmd.Flags().GetString("back")
		if err != nil {
			return err
		}
		deck, err := cmd.Flags().GetString("deck")
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		var status string
		status = ddb.Todo.String()
		newcard := ddb.Card{uint(id), front, back, deck, status, time.Time{}}
		return c.Update(newcard)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your cards",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := ddb.OpenDb(setupPath())
		if err != nil {
			return err
		}
		defer c.Db.Close()
		cards, err := c.Getcards()
		if err != nil {
			return err
		}
		fmt.Print(setupTable(cards))
		return nil
	},
}

func setupTable(cards []ddb.Card) *table.Table {
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

func cardsToItems(cards []ddb.Card) []list.Item {
	var items []list.Item
	for _, t := range cards {
		items = append(items, t)
	}
	return items
}

func init() {
	addCmd.Flags().StringP(
		"deck",
		"d",
		"",
		"specify a deck for your card",
	)
	addCmd.Flags().StringP(
		"front",
		"f",
		"",
		"specify the front for your card",
	)
	addCmd.Flags().StringP(
		"back",
		"b",
		"",
		"specify the back for your card",
	)
	addCmd.Flags().IntP(
		"status",
		"s",
		int(ddb.Todo),
		"specify a status for your card",
	)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(whereCmd)
	rootCmd.AddCommand(deleteCmd)
}
