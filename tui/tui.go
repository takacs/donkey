package tui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func StartTea() error {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	termWidth, termHeight, _ := term.GetSize(int(os.Stdin.Fd()))
	m, _ := InitProject(termWidth, termHeight)
	program := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return nil
}
