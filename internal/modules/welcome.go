package modules

import (
	"context"
	"math/rand"
	"strconv"
	"strings"

	"github.com/tajamullone106-droid/Blood/internal/database"
	"github.com/tajamullone106-droid/Blood/internal/utils"
)

type NewMember struct {
	UserID      int64
	FirstName   string
	LastName    string
	Username    string
	ChatID      int64
	ChatTitle   string
	MemberCount int
}

func (n NewMember) mention() string {
	name := n.FirstName
	if name == "" {
		name = "there"
	}
	return utils.Mention(name, n.UserID)
}

type RenderedWelcome struct {
	Text    string
	Buttons [][]utils.Btn
}

func fillPlaceholders(tmpl string, n NewMember) string {
	username := "no username"
	if n.Username != "" {
		username = "@" + n.Username
	}
	replacer := strings.NewReplacer(
		"{first}", n.FirstName,
		"{last}", n.LastName,
		"{fullname}", strings.TrimSpace(n.FirstName+" "+n.LastName),
		"{mention}", n.mention(),
		"{username}", username,
		"{id}", strconv.FormatInt(n.UserID, 10),
		"{chatname}", n.ChatTitle,
		"{count}", strconv.Itoa(n.MemberCount),
	)
	return replacer.Replace(tmpl)
}

func BuildWelcome(ctx context.Context, store *database.WelcomeStore, n NewMember) (*RenderedWelcome, error) {
	settings, err := store.Get(ctx, n.ChatID)
	if err != nil {
		return nil, err
	}
	if !settings.Enabled || len(settings.Messages) == 0 {
		return nil, nil
	}

	tmpl := settings.Messages[rand.Intn(len(settings.Messages))]
	text := fillPlaceholders(tmpl, n)

	return &RenderedWelcome{
		Text:    text,
		Buttons: utils.WelcomeButtons(settings.ButtonURL),
	}, nil
}

func StyleChatName(chatName string) string {
	return utils.StyleText(chatName, utils.FontBoldSansSerif)
}
