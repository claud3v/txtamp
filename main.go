package main

import (
	"log"

	"txtamp/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
