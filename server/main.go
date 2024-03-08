package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kacperhemperek/discord-go/api"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}
}
func main() {
	s := api.NewApiServer(8080)
	fmt.Println("Change")
	s.Start()
}
