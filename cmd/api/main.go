package main

import (
	"log"

	"github.com/olegrom32/imperial-fleet-api/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		log.Fatal(err)
	}
}
