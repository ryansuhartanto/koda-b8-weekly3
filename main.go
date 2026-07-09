package main

import (
	"fmt"
	"log"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

func main() {
	restaurant := lipgloss.NewStyle().
		Foreground(lipgloss.BrightGreen).
		Background(lipgloss.Black).
		Render("Wingstop")

	lipgloss.Printf("Selamat datang di %v!\n", restaurant)

	var (
		exit bool
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Exit?").
				Value(&exit),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	if !exit {
		fmt.Println("?")
	}

}
