package bot

import (
	"fmt"

	"github.com/boltdb/bolt"
	tb "gopkg.in/tucnak/telebot.v2"
)

const statusResponse = `
Username: <b>%s</b>
Activated: <b>%v</b>
`

func handleStatus(b *tb.Bot, db *bolt.DB) func(m *tb.Message) {
	return func(m *tb.Message) {
		user, _ := getChatInfo(db, m.Sender.ID)

		if user == nil {
			b.Send(m.Sender, "No data for you")
		} else {
			b.Send(m.Sender, fmt.Sprintf(statusResponse, user.Username, user.Activated), tb.ModeHTML)
		}
	}
}
