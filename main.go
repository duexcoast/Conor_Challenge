package main

import (
	"log"
	"os"
)

func main() {
	if err := validateInput(os.Stdin, os.Stdout, validate); err != nil {
		log.Fatal(err)
	}
}
