package main

import (
	"fmt"
	"github.com/kacperhemperek/discord-go/api"
)

func main() {
	s := api.NewApiServer(8080)
	fmt.Println("Change")
	s.Start()
}
