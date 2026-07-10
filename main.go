package main

import (
	_ "embed"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ryansuhartanto/koda-b8-weekly3/model"
)

//go:embed menu.json
var json string

func main() {
	data := model.NewData([]byte(json))
	data.Restaurant = model.Restaurant(lipgloss.NewStyle().
		Foreground(lipgloss.BrightGreen).
		Background(lipgloss.Black).
		Render(string(data.Restaurant)))

	if _, err := tea.NewProgram(model.NewMain(data)).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
