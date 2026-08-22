package utils

type ButtonColor string

const (
	ColorRed    ButtonColor = "🔴"
	ColorGreen  ButtonColor = "🟢"
	ColorBlue   ButtonColor = "🔵"
	ColorYellow ButtonColor = "🟡"
	ColorPurple ButtonColor = "🟣"
	ColorOrange ButtonColor = "🟠"
	ColorWhite  ButtonColor = "⚪"
	ColorBlack  ButtonColor = "⚫"
)

type Btn struct {
	Text  string
	Data  string
	IsURL bool
}

func Colored(color ButtonColor, label, callbackData string) Btn {
	return Btn{Text: string(color) + " " + label, Data: callbackData}
}

func Row(buttons ...Btn) []Btn { return buttons }

func PlayerControls(chatID int64) [][]Btn {
	return [][]Btn{
		{
			Colored(ColorGreen, "Resume", "player:resume"),
			Colored(ColorYellow, "Pause", "player:pause"),
			Colored(ColorBlue, "Skip", "player:skip"),
			Colored(ColorRed, "Stop", "player:stop"),
		},
		{
			Colored(ColorPurple, "Loop", "player:loop"),
			Colored(ColorOrange, "Shuffle", "player:shuffle"),
			Colored(ColorWhite, "Queue", "player:queue"),
		},
	}
}

func WelcomeButtons(rulesURL string) [][]Btn {
	row := []Btn{}
	if rulesURL != "" {
		row = append(row, Btn{Text: string(ColorBlue) + " Rules", Data: rulesURL, IsURL: true})
	}
	row = append(row, Colored(ColorGreen, "Help", "welcome:help"))
	return [][]Btn{row}
}
