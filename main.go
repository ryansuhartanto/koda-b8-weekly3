package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ryansuhartanto/koda-b8-weekly3/model"
)

func main() {
	restaurant := lipgloss.NewStyle().
		Foreground(lipgloss.BrightGreen).
		Background(lipgloss.Black).
		Render("Wingstop")

	if _, err := tea.NewProgram(model.NewMain(restaurant)).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
