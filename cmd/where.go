package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/takacs/donkey/internal/database"
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
		path, err := database.GetDbPath()
		if err != nil {
			log.Fatal("can't find path")
		}
		fmt.Println(path)
		return err
	},
}
