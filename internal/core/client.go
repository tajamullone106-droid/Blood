package core

import (
	"context"
	"log"

	"github.com/tajamullone106-droid/Blood/internal/config"
	"github.com/tajamullone106-droid/Blood/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Bot struct {
	Cfg               *config.Config
	Mongo             *mongo.Client
	DB                *mongo.Database
	WelcomeStore      *database.WelcomeStore
	ChatSettingsStore *database.ChatSettingsStore
}

func New(cfg *config.Config) (*Bot, error) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	db := client.Database("blood_bot")

	return &Bot{
		Cfg:               cfg,
		Mongo:             client,
		DB:                db,
		WelcomeStore:      database.NewWelcomeStore(db),
		ChatSettingsStore: database.NewChatSettingsStore(db),
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	log.Printf("[core] Blood bot (@%s) starting — owner=%d, sudo=%v", b.Cfg.BotUsername, b.Cfg.OwnerID, b.Cfg.SudoUsers)
	log.Println("[core] TODO: wire in your Telegram client (gogram/gotdbot) here and register handlers from internal/modules")
	select {}
}
