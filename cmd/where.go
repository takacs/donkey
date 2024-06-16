package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/takacs/donkey/carddb"
	"log"
)

func init() {
	rootCmd.AddCommand(whereCmd)
}

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show where your cards are stored",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		carddb, err := carddb.New()
		if err != nil {
			log.Fatal("can't open db")
		}
		fmt.Println(carddb.DataDir)
		return err
	},
}
