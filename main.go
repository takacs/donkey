package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/takacs/donkey/cmd"
)

func main() {
	cmd.Execute()
}
