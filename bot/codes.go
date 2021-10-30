package bot

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

const noCodesLeft = `
You have <b>%d</b> invites left.

Your secret code: <b>%d</b>
`

func (bot *Bot) handleCodes(m *tb.Message) {
	user, _ := bot.getChatInfo(m.Sender.ID)

	if user == nil {
		bot.Tg.Send(m.Sender, "You are not registered to share codes!")
	} else {
		bot.Tg.Send(m.Sender, fmt.Sprintf(noCodesLeft, user.Invites, m.Sender.ID), tb.ModeHTML)
	}
}
