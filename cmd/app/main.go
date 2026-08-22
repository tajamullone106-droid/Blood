package main

import (
	"context"
	"log"

	"github.com/tajamullone106-droid/Blood/internal/config"
	"github.com/tajamullone106-droid/Blood/internal/core"
)

func main() {
	cfg := config.Load()

	bot, err := core.New(cfg)
	if err != nil {
		log.Fatalf("[main] failed to initialize bot: %v", err)
	}

	if err := bot.Start(context.Background()); err != nil {
		log.Fatalf("[main] bot stopped with error: %v", err)
	}
}
