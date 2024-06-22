package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/takacs/donkey/internal/card"
	"github.com/takacs/donkey/internal/supermemo"
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
		err = loadFileToDb(args[0], deck, carddb)
		if err != nil {
			return err
		}
		return nil
	},
}

func loadFileToDb(path string, deck string, cdb *card.CardDb) error {
	if strings.HasSuffix(path, ".apkg") {
		return errors.New(`.apkg format is not supported as of now. please export deck with Export Format -> Cards in Plain Text and uncheck Include HTML and media references from Anki and try again.`)
	}
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	inserted := 1
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if len(line) < 1 || line[0] == '#' {
			fmt.Printf("couldn't import line %v\n", inserted)
			continue
		}
		fields := strings.Split(line, "\t")
		cardId, err := cdb.Insert(fields[0], fields[1], deck)
		if err != nil {
			fmt.Printf("failed importing %v | %v\n", fields[0], fields[1])
		} else {
			fmt.Printf("imported card front: %v | back: %v", fields[0], fields[1])
		}
		supermemoDb, err := supermemo.New()
		if err != nil {
			fmt.Printf("couldn't import line %v\n", inserted)
		}
		err = supermemoDb.Insert(cardId)
		if err != nil {
			fmt.Printf("couldn't import line %v\n", inserted)
			continue
		}

		inserted++
	}

	fmt.Printf("inserted %v cards\n", inserted)

	if err := scanner.Err(); err != nil {
		return err
	}
	return err
}
