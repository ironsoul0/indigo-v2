package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	schema "github.com/ironsoul0/indigo-v2/db"
	tb "gopkg.in/tucnak/telebot.v2"
)

func New(token string, db *bolt.DB) *tb.Bot {
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/hello", func(m *tb.Message) {
		bot.Send(m.Sender, "Hello World!")
		chatID := fmt.Sprintf("%d", m.Sender.ID)

		db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(schema.USERS_BUCKET))
			user := schema.User{
				ChatID: chatID,
			}
			buf, _ := json.Marshal(user)
			err := bucket.Put([]byte(chatID), buf)

			if err != nil {
				return err
			}

			return nil
		})
	})

	bot.Handle("/get", func(m *tb.Message) {
		chatID := fmt.Sprintf("%d", m.Sender.ID)

		db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(schema.USERS_BUCKET))
			bufUser := bucket.Get([]byte(chatID))

			var user schema.User
			json.Unmarshal(bufUser, &user)

			bot.Send(m.Sender, user.ChatID)

			return nil
		})
	})

	go bot.Start()

	return bot
}
