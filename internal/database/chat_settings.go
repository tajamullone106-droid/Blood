package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatSettings struct {
	ChatID        int64  `bson:"chat_id"`
	WelcomeOn     bool   `bson:"welcome_on"`
	PlayModeChat  bool   `bson:"play_mode_chat"`
	FontStyle     string `bson:"font_style"`
	CleanCommands bool   `bson:"clean_commands"`
}

func defaultChatSettings(chatID int64) *ChatSettings {
	return &ChatSettings{
		ChatID:       chatID,
		WelcomeOn:    true,
		PlayModeChat: false,
		FontStyle:    "",
		CleanCommands: false,
	}
}

type ChatSettingsStore struct {
	col *mongo.Collection
}

func NewChatSettingsStore(db *mongo.Database) *ChatSettingsStore {
	return &ChatSettingsStore{col: db.Collection("chat_settings")}
}

func (s *ChatSettingsStore) Get(ctx context.Context, chatID int64) (*ChatSettings, error) {
	var cs ChatSettings
	err := s.col.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&cs)
	if err == mongo.ErrNoDocuments {
		return defaultChatSettings(chatID), nil
	}
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

func (s *ChatSettingsStore) SetField(ctx context.Context, chatID int64, field string, value interface{}) error {
	_, err := s.col.UpdateOne(ctx,
		bson.M{"chat_id": chatID},
		bson.M{
			"$set":         bson.M{field: value},
			"$setOnInsert": bson.M{"chat_id": chatID},
		},
		options.Update().SetUpsert(true),
	)
	return err
}
