package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/takacs/donkey/internal/card"
)

func init() {
	loadCmd.Flags().StringP(
		"deck",
		"d",
		"default",
		"specify a deck for your card",
	)
	rootCmd.AddCommand(loadCmd)
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "load exported anki file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		carddb, err := card.New()
		if err != nil {
			return err
		}
		defer carddb.Close()
		deck, err := cmd.Flags().GetString("deck")
		if err != nil {
			return err
		}
		loadFileToDb(args[0], deck, carddb)
		return nil
	},
}

func loadFileToDb(path string, deck string, cdb *card.CardDb) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var inserted int
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 || line[0] == '#' {
			continue
		}
		fields := strings.Split(line, "\t")
		err := cdb.Insert(fields[0], fields[1], deck)
		if err != nil {
			log.Printf("failed importing %v | %v", fields[0], fields[1])
		}
		inserted++
	}

	fmt.Printf("inserted %v cards\n", inserted)

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
}
