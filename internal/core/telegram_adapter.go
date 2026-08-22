package core

import (
	"github.com/amarnathcjd/gogram/telegram"
	"github.com/tajamullone106-droid/Blood/internal/utils"
)

func buildKeyboard(rows [][]utils.Btn) *telegram.ReplyMarkup {
	if len(rows) == 0 {
		return nil
	}
	kb := telegram.NewKeyboard()
	for _, row := range rows {
		var tgButtons []telegram.KeyboardButton
		for _, b := range row {
			if b.IsURL {
				tgButtons = append(tgButtons, telegram.Button.URL(b.Text, b.Data))
			} else {
				tgButtons = append(tgButtons, telegram.Button.Data(b.Text, b.Data))
			}
		}
		kb.AddRow(tgButtons...)
	}
	return kb
}
