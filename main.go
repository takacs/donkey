package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/takacs/donkey/cmd"
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
	cmd.Execute()
}
