package cmd

import (
	"github.com/spf13/cobra"
	ddb "github.com/takacs/donkey/db"
)

func init() {
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
	addCmd.Flags().StringP(
		"deck",
		"d",
		"default",
		"specify a deck for your card",
	)
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new card with an optional deck name",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := ddb.OpenDb(ddb.SetupPath())
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
