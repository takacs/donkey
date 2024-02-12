package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
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
	Use:   "add Card",
	Short: "Add a new card with an optional deck name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := openDB(setupPath())
		if err != nil {
			return err
		}
		defer c.db.Close()
		deck, err := cmd.Flags().GetString("deck")
		if err != nil {
			return err
		}
		if err := c.insert(args[0], args[0], deck); err != nil {
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
		c, err := openDB(setupPath())
		if err != nil {
			return err
		}
		defer c.db.Close()
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		return c.delete(uint(id))
	},
}

var updateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a card by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := openDB(setupPath())
		if err != nil {
			return err
		}
		defer c.db.Close()
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}
		prog, err := cmd.Flags().GetInt("status")
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		deck := "default"
		var status string
		switch prog {
		case int(inProgress):
			status = inProgress.String()
		case int(done):
			status = done.String()
		default:
			status = todo.String()
		}
		newcard := card{uint(id), name, project, status, deck, time.Time{}}
		return c.update(newcard)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your cards",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := openDB(setupPath())
		if err != nil {
			return err
		}
		defer c.db.Close()
		cards, err := c.getcards()
		if err != nil {
			return err
		}
		fmt.Print(setupTable(cards))
		return nil
	},
}

func setupTable(cards []card) *table.Table {
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

func cardsToItems(cards []card) []list.Item {
	var items []list.Item
	for _, t := range cards {
		items = append(items, t)
	}
	return items
}

func init() {
	addCmd.Flags().StringP(
		"project",
		"p",
		"",
		"specify a project for your card",
	)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	updateCmd.Flags().StringP(
		"name",
		"n",
		"",
		"specify a name for your card",
	)
	updateCmd.Flags().StringP(
		"project",
		"p",
		"",
		"specify a project for your card",
	)
	updateCmd.Flags().IntP(
		"status",
		"s",
		int(todo),
		"specify a status for your card",
	)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(whereCmd)
	rootCmd.AddCommand(deleteCmd)
}
