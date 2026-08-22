package core

import (
	"context"
	"log"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/tajamullone106-droid/Blood/internal/config"
	"github.com/tajamullone106-droid/Blood/internal/database"
	"github.com/tajamullone106-droid/Blood/internal/modules"
	"github.com/tajamullone106-droid/Blood/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	mopts "go.mongodb.org/mongo-driver/mongo/options"
)

type Bot struct {
	Cfg               *config.Config
	Mongo             *mongo.Client
	DB                *mongo.Database
	WelcomeStore      *database.WelcomeStore
	ChatSettingsStore *database.ChatSettingsStore
	Client            *telegram.Client
}

func New(cfg *config.Config) (*Bot, error) {
	ctx := context.Background()

	mclient, err := mongo.Connect(ctx, mopts.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}
	if err := mclient.Ping(ctx, nil); err != nil {
		return nil, err
	}
	db := mclient.Database("blood_bot")

	return &Bot{
		Cfg:               cfg,
		Mongo:             mclient,
		DB:                db,
		WelcomeStore:      database.NewWelcomeStore(db),
		ChatSettingsStore: database.NewChatSettingsStore(db),
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	client, err := telegram.NewClient(telegram.ClientConfig{
		AppID:   b.Cfg.APIID,
		AppHash: b.Cfg.APIHash,
	})
	if err != nil {
		return err
	}
	b.Client = client

	if err := client.Conn(); err != nil {
		return err
	}
	if _, err := client.LoginBot(b.Cfg.BotToken); err != nil {
		return err
	}

	b.registerHandlers(ctx)

	log.Printf("[core] Blood bot (@%s) is live — owner=%d, sudo=%v", b.Cfg.BotUsername, b.Cfg.OwnerID, b.Cfg.SudoUsers)
	client.Idle()
	return nil
}

func (b *Bot) registerHandlers(ctx context.Context) {
	client := b.Client

	client.On(telegram.OnCommand("start"), func(m *telegram.NewMessage) error {
		return m.Reply(utils.Bold("Kokomimusicbot") + " is up! Use /settings to configure this chat, or /play to start music.")
	})

	client.On(telegram.OnCommand("settings"), func(m *telegram.NewMessage) error {
		chatTitle := ""
		if chat := m.Chat(); chat != nil {
			chatTitle = chat.Title
		}
		menu := modules.MainSettingsMenu(chatTitle)
		_, err := m.Reply(menu.Text, telegram.SendOptions{ReplyMarkup: buildKeyboard(menu.Buttons)})
		return err
	})

	client.On("callback", func(c *telegram.CallbackQuery) error {
		data := c.DataString()
		if len(data) < 9 || data[:9] != "settings:" {
			return nil
		}

		chatID := c.ChatID
		chatTitle := ""
		if chat := c.Chat(); chat != nil {
			chatTitle = chat.Title
		}

		menu, err := modules.HandleSettingsCallback(ctx, b.ChatSettingsStore, chatID, chatTitle, data)
		if err != nil {
			c.Answer("Something went wrong, try again.")
			return err
		}
		c.Edit(menu.Text, &telegram.SendOptions{ReplyMarkup: buildKeyboard(menu.Buttons)})
		c.Answer("")
		return nil
	})

	// VERIFY: exact event name/payload for join events can differ by
	// gogram version — check https://gogram.fun/handlers if this doesn't
	// compile as-is, and adjust field names (UserJoined / Users / etc.)
	client.On("chataction", func(ca *telegram.ChatAction) error {
		if !ca.UserJoined && !ca.UserAdded {
			return nil
		}
		user := ca.User()
		chat := ca.Chat()
		if user == nil || chat == nil {
			return nil
		}

		member := modules.NewMember{
			UserID:      user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Username:    user.Username,
			ChatID:      chat.ID,
			ChatTitle:   chat.Title,
			MemberCount: int(chat.ParticipantsCount),
		}

		welcome, err := modules.BuildWelcome(ctx, b.WelcomeStore, member)
		if err != nil || welcome == nil {
			return err
		}
		_, err = client.SendMessage(chat.ID, welcome.Text, telegram.SendOptions{
			ReplyMarkup: buildKeyboard(welcome.Buttons),
			ParseMode:   "html",
		})
		return err
	})
}
