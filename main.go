package main

import (
	"fmt"
	"log"

	"github.com/ryansuhartanto/koda-b8-weekly3/form"
)

func main() {
	var main form.Main

	err := main.Form().Run()
	if err != nil {
		log.Fatal(err)
	}

	if !main.Exit {
		fmt.Println("?")
	}

}
