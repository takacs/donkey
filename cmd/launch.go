package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	tui "github.com/takacs/donkey/tui"
)

func init() {
	rootCmd.AddCommand(launchCmd)
}

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "launch donkey app",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := tui.StartTea()
		if err != nil {
			fmt.Println("failed to start donkey")
		}
		return nil
	},
}
