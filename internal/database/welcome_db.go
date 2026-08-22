package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WelcomeSetting struct {
	ChatID         int64    `bson:"chat_id"`
	Enabled        bool     `bson:"enabled"`
	Messages       []string `bson:"messages"`
	ButtonURL      string   `bson:"button_url"`
	DeletePrevious bool     `bson:"delete_previous"`
	CleanTimer     int      `bson:"clean_timer"`
}

const DefaultWelcomeMessage = "👋 Welcome {mention} to {chatname}!\nHave a great time here 🎶"

type WelcomeStore struct {
	col *mongo.Collection
}

func NewWelcomeStore(db *mongo.Database) *WelcomeStore {
	return &WelcomeStore{col: db.Collection("welcome_settings")}
}

func (s *WelcomeStore) Get(ctx context.Context, chatID int64) (*WelcomeSetting, error) {
	var w WelcomeSetting
	err := s.col.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&w)
	if err == mongo.ErrNoDocuments {
		return &WelcomeSetting{
			ChatID:   chatID,
			Enabled:  true,
			Messages: []string{DefaultWelcomeMessage},
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (s *WelcomeStore) Save(ctx context.Context, w *WelcomeSetting) error {
	_, err := s.col.UpdateOne(ctx,
		bson.M{"chat_id": w.ChatID},
		bson.M{"$set": w},
		options.Update().SetUpsert(true),
	)
	return err
}

func (s *WelcomeStore) AddMessage(ctx context.Context, chatID int64, msg string) error {
	_, err := s.col.UpdateOne(ctx,
		bson.M{"chat_id": chatID},
		bson.M{
			"$addToSet":    bson.M{"messages": msg},
			"$setOnInsert": bson.M{"enabled": true},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

func (s *WelcomeStore) SetEnabled(ctx context.Context, chatID int64, enabled bool) error {
	_, err := s.col.UpdateOne(ctx,
		bson.M{"chat_id": chatID},
		bson.M{"$set": bson.M{"enabled": enabled}},
		options.Update().SetUpsert(true),
	)
	return err
}
