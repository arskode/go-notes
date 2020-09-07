package main

import (
	"log"

	"github.com/arskode/go-notes/api/rest"
)

func main() {

	err := rest.Start()
	if err != nil {
		log.Fatal(err)
	}
}
