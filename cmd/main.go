package main

import (
	"log"
	"project/internal/app"
)

func main() {
	log.Println("API start!")

	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	// Запуск приложения
	application.StartServer()

	log.Println("API is shitting down!")
}
