package bot

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
	schema "github.com/ironsoul0/indigo-v2/db"
	"github.com/ironsoul0/indigo-v2/scrapers/moodle"
	tb "gopkg.in/tucnak/telebot.v2"
)

func handleRegister(b *tb.Bot, db *bolt.DB) func(m *tb.Message) {
	return func(m *tb.Message) {
		parts := strings.Split(m.Payload, " ")
		if len(parts) < 2 {
			b.Send(m.Sender, "Invalid credentials")
			return
		}
		username, password := parts[0], parts[1]

		client := moodle.Init()
		response := client.GetGrades(username, password)

		if response.InvalidCredentials || response.RequestFailure {
			b.Send(m.Sender, "Invalid credentials")
			return
		}

		err := db.Update(func(tx *bolt.Tx) error {
			chatID := fmt.Sprintf("%d", m.Sender.ID)
			bucket := tx.Bucket([]byte(schema.USERS_BUCKET))

			user := schema.User{
				ChatID:    m.Sender.ID,
				Username:  username,
				Password:  password,
				Activated: true,
				Courses:   response.Courses,
			}
			buf, _ := json.Marshal(user)
			err := bucket.Put([]byte(chatID), buf)

			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			b.Send(m.Sender, "Something wrong happended. Try again.")
			return
		}

		b.Send(m.Sender, "Notifications are on")
	}
}
