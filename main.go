package main

import (
	"fmt"

	"github.com/fatih/color"
)

func main() {
	fmt.Print("Selamat datang di ")
	color.New(color.BgBlack).Add(color.FgHiGreen).Add(color.Bold).Print("Wingstop")
	fmt.Println("!")
}
