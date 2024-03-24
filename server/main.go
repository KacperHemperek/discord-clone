package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kacperhemperek/discord-go/api"
	"os"
	"os/signal"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	s := api.NewApiServer(8080)
	return s.Start()
}

func main() {

	ctx := context.Background()
	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
