package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	gap "github.com/muesli/go-app-paths"
)

/*
Here's the plan, we're going to store our data in a dedicated data directory at
`XDG_DATA_HOME/.tasks`. That's where we will store a copy of our SQLite DB.
*/

// setupPath uses XDG to create the necessary data dirs for the program.
func setupPath() string {
	// get XDG paths
	scope := gap.NewScope(gap.User, "cards")
	dirs, err := scope.DataDirs()
	if err != nil {
		log.Fatal(err)
	}
	// create the app base dir, if it doesn't exist
	var cardDir string
	if len(dirs) > 0 {
		cardDir = dirs[0]
	} else {
		cardDir, _ = os.UserHomeDir()
	}
	if err := initCardDir(cardDir); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cardDir)
	return cardDir
}

func initCardDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o750)
		}
		return err
	}
	return nil
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
