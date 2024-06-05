package main

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	db "github.com/takacs/donkey/db"
)

func setupPath() string {
	cardDir, err := db.GetDbPath("cards")
	if err != nil {
		fmt.Println("error getting db path")
	}
	return cardDir
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
