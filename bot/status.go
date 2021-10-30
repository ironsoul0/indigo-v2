package bot

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

const statusResponse = `
Username: <b>%s</b>
Activated: <b>%v</b>
`

func (bot *Bot) handleStatus(m *tb.Message) {
	user, _ := bot.getChatInfo(m.Sender.ID)

	if user == nil {
		bot.Tg.Send(m.Sender, "Nothing found 😢")
	} else {
		bot.Tg.Send(m.Sender, fmt.Sprintf(statusResponse, user.Username, user.Activated), tb.ModeHTML)
	}
}
