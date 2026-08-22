package modules

import (
	"context"
	"fmt"

	"github.com/tajamullone106-droid/Blood/internal/database"
	"github.com/tajamullone106-droid/Blood/internal/utils"
)

type Menu struct {
	Text    string
	Buttons [][]utils.Btn
}

const (
	cbMain           = "settings:main"
	cbWelcomeMenu    = "settings:welcome"
	cbWelcomeToggle  = "settings:welcome:toggle"
	cbFontMenu       = "settings:font"
	cbFontSetPrefix  = "settings:font:set:"
	cbPlayerMenu     = "settings:player"
	cbPlayModeToggle = "settings:player:toggle_playmode"
	cbCleanToggle    = "settings:player:toggle_clean"
	cbClose          = "settings:close"
)

func onOff(b bool) string {
	if b {
		return "✅ ON"
	}
	return "❌ OFF"
}

func MainSettingsMenu(chatTitle string) Menu {
	return Menu{
		Text: utils.Bold(StyleChatName("Settings")) + " — " + chatTitle + "\n\nChoose what to configure:",
		Buttons: [][]utils.Btn{
			{utils.Colored(utils.ColorGreen, "Welcome Message", cbWelcomeMenu)},
			{utils.Colored(utils.ColorPurple, "Text Font Style", cbFontMenu)},
			{utils.Colored(utils.ColorBlue, "Player", cbPlayerMenu)},
			{utils.Colored(utils.ColorRed, "Close", cbClose)},
		},
	}
}

func welcomeSettingsMenu(cs *database.ChatSettings) Menu {
	return Menu{
		Text: utils.Bold("Welcome Message Settings") + "\n\nStatus: " + onOff(cs.WelcomeOn) +
			"\n\nUse /setwelcome <text> to change the message template.\nSupported placeholders: {first} {mention} {chatname} {count}",
		Buttons: [][]utils.Btn{
			{utils.Colored(utils.ColorGreen, "Toggle Welcome "+onOff(cs.WelcomeOn), cbWelcomeToggle)},
			{utils.Colored(utils.ColorOrange, "⬅ Back", cbMain)},
		},
	}
}

func fontSettingsMenu(cs *database.ChatSettings) Menu {
	current := cs.FontStyle
	if current == "" {
		current = "plain (default)"
	}
	label := func(style utils.FontStyle, name string) utils.Btn {
		prefix := "⬜"
		if cs.FontStyle == string(style) {
			prefix = "✅"
		}
		return utils.Btn{Text: prefix + " " + name, Data: cbFontSetPrefix + string(style)}
	}
	return Menu{
		Text: utils.Bold("Text Font Style") + "\n\nCurrent: " + current +
			"\n\nThis controls how names/titles are styled in bot messages (e.g. welcome header).",
		Buttons: [][]utils.Btn{
			{label(utils.FontBoldSansSerif, utils.StyleText("Bold Sans", utils.FontBoldSansSerif))},
			{label(utils.FontScript, utils.StyleText("Script", utils.FontScript))},
			{label(utils.FontMonospace, utils.StyleText("Monospace", utils.FontMonospace))},
			{label(utils.FontSmallCaps, utils.StyleText("Small Caps", utils.FontSmallCaps))},
			{label(utils.FontDoubleStruck, utils.StyleText("Double Struck", utils.FontDoubleStruck))},
			{utils.Btn{Text: "⬜ Plain (no style)", Data: cbFontSetPrefix + "plain"}},
			{utils.Colored(utils.ColorOrange, "⬅ Back", cbMain)},
		},
	}
}

func playerSettingsMenu(cs *database.ChatSettings) Menu {
	return Menu{
		Text: utils.Bold("Player Settings") + "\n\n" +
			"Admin-only play mode: " + onOff(cs.PlayModeChat) + "\n" +
			"Auto-clean command messages: " + onOff(cs.CleanCommands),
		Buttons: [][]utils.Btn{
			{utils.Colored(utils.ColorBlue, "Toggle Admin-only Play", cbPlayModeToggle)},
			{utils.Colored(utils.ColorYellow, "Toggle Clean Commands", cbCleanToggle)},
			{utils.Colored(utils.ColorOrange, "⬅ Back", cbMain)},
		},
	}
}

func HandleSettingsCallback(ctx context.Context, store *database.ChatSettingsStore, chatID int64, chatTitle, data string) (*Menu, error) {
	cs, err := store.Get(ctx, chatID)
	if err != nil {
		return nil, err
	}

	switch {
	case data == cbMain:
		m := MainSettingsMenu(chatTitle)
		return &m, nil

	case data == cbWelcomeMenu:
		m := welcomeSettingsMenu(cs)
		return &m, nil
	case data == cbWelcomeToggle:
		cs.WelcomeOn = !cs.WelcomeOn
		if err := store.SetField(ctx, chatID, "welcome_on", cs.WelcomeOn); err != nil {
			return nil, err
		}
		m := welcomeSettingsMenu(cs)
		return &m, nil

	case data == cbFontMenu:
		m := fontSettingsMenu(cs)
		return &m, nil
	case len(data) > len(cbFontSetPrefix) && data[:len(cbFontSetPrefix)] == cbFontSetPrefix:
		style := data[len(cbFontSetPrefix):]
		if style == "plain" {
			style = ""
		}
		cs.FontStyle = style
		if err := store.SetField(ctx, chatID, "font_style", style); err != nil {
			return nil, err
		}
		m := fontSettingsMenu(cs)
		return &m, nil

	case data == cbPlayerMenu:
		m := playerSettingsMenu(cs)
		return &m, nil
	case data == cbPlayModeToggle:
		cs.PlayModeChat = !cs.PlayModeChat
		if err := store.SetField(ctx, chatID, "play_mode_chat", cs.PlayModeChat); err != nil {
			return nil, err
		}
		m := playerSettingsMenu(cs)
		return &m, nil
	case data == cbCleanToggle:
		cs.CleanCommands = !cs.CleanCommands
		if err := store.SetField(ctx, chatID, "clean_commands", cs.CleanCommands); err != nil {
			return nil, err
		}
		m := playerSettingsMenu(cs)
		return &m, nil

	case data == cbClose:
		return &Menu{Text: "Settings closed.", Buttons: nil}, nil
	}

	return nil, fmt.Errorf("unknown settings callback: %s", data)
}
