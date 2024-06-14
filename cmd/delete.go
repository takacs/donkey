package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	ddb "github.com/takacs/donkey/db"
	"strconv"
)

func init() {
	deleteCmd.Flags().StringP(
		"id",
		"i",
		"",
		"specify a deck for your card",
	)
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a card based on id. get a list of ids with `donkey list`",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := ddb.OpenDb(ddb.SetupPath())
		if err != nil {
			return err
		}
		defer c.Db.Close()
		ids, err := cmd.Flags().GetString("id")
		if err != nil {
			fmt.Println("-id flag is required")
			return err
		}
		id, err := strconv.Atoi(ids)
		if err != nil {
			fmt.Println("-id flag is required")
			fmt.Println()
		}
		return c.Delete(uint(id))
	},
}
