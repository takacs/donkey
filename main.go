package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/takacs/donkey/cmd"
	"github.com/takacs/donkey/tui"
	"os"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "cli" {
		cmd.Execute()
	} else {
		err := tui.StartTea()
		if err != nil {
			fmt.Println("failed to start donkey")
		}
	}
}
