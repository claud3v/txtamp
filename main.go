package main

import (
	"log"

	"txtamp/auth"
)

func main() {
	if err := auth.Run(); err != nil {
		log.Fatal(err)
	}
}
