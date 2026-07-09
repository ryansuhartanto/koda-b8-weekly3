package main

import (
	"fmt"
	"log"

	"charm.land/huh/v2"
	"github.com/fatih/color"
)

func main() {
	fmt.Print("Selamat datang di ")
	color.New(color.BgBlack).Add(color.FgHiGreen).Add(color.Bold).Print("Wingstop")
	fmt.Println("!")

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
