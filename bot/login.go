package bot

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	schema "github.com/ironsoul0/indigo-v2/db"
	tb "gopkg.in/tucnak/telebot.v2"
)

const loginResponse = `
Username: <b>%s</b>
Activated: <b>%v</b>
`

var (
	menu     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	selector = &tb.ReplyMarkup{}

	btnHelp     = menu.Text("ℹ Help")
	btnSettings = menu.Text("⚙ Settings")
)

func handleLogin(b *tb.Bot, db *bolt.DB) func(m *tb.Message) {
	return func(m *tb.Message) {
		db.View(func(tx *bolt.Tx) error {
			chatID := fmt.Sprintf("%d", m.Sender.ID)
			bucket := tx.Bucket([]byte(schema.USERS_BUCKET))
			userData := bucket.Get([]byte(chatID))

			if userData == nil {
				b.Send(m.Sender, "No data for you")
			} else {
				var user schema.User
				json.Unmarshal(userData, &user)

				b.Send(m.Sender, fmt.Sprintf(loginResponse, user.Username, user.Activated), tb.ModeHTML)
			}

			return nil
		})
	}
}
