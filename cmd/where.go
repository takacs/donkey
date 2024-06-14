package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	ddb "github.com/takacs/donkey/db"
)

func init() {
	rootCmd.AddCommand(whereCmd)
}

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show where your cards are stored",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Println(ddb.SetupPath())
		return err
	},
}
