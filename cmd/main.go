package main

import (
	"log"
	"project/internal/api"
)

func main() {
	log.Println("API start!")
	api.StartServer()
	log.Println("API is shitting down!")
}
