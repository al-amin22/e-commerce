package main

import (
	"log"

	"ecommerce-backend/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
