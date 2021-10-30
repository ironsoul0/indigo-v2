package bot

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/ironsoul0/indigo-v2/db"
	schema "github.com/ironsoul0/indigo-v2/db"
	"github.com/ironsoul0/indigo-v2/scrapers/moodle"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (bot *Bot) handleRegister(m *tb.Message) {
	bot.updateScene(m.Sender.ID, db.REGISTER_SCENE)
	bot.Tg.Send(m.Sender, "Welcome aboard! What is your username?", &tb.ReplyMarkup{ReplyKeyboardRemove: true})
}

func (bot *Bot) handleRegisterStep(m *tb.Message, scene *db.Scene) {
	if scene.Step == 0 {
		bot.db.Update(func(tx *bolt.Tx) error {
			chatID := fmt.Sprintf("%d", m.Sender.ID)
			scene.Step += 1
			scene.Username = m.Text
			buf, _ := json.Marshal(scene)

			bucket := tx.Bucket([]byte(db.STATES_BUCKET))
			bucket.Put([]byte(chatID), buf)

			return nil
		})

		bot.Tg.Send(m.Sender, "What is your password?")
		return
	}

	password := m.Text
	username := scene.Username
	client := moodle.Init()
	response := client.GetGrades(username, password)

	if response.InvalidCredentials || response.RequestFailure {
		bot.Tg.Send(m.Sender, "Invalid credentials!", mainMenu)
		bot.updateScene(m.Sender.ID, db.NO_SCENE)
		return
	}

	bot.updateScene(m.Sender.ID, db.NO_SCENE)

	err := bot.db.Update(func(tx *bolt.Tx) error {
		chatID := fmt.Sprintf("%d", m.Sender.ID)
		bucket := tx.Bucket([]byte(schema.USERS_BUCKET))

		user := schema.User{
			ChatID:    m.Sender.ID,
			Username:  username,
			Password:  password,
			Activated: true,
			Courses:   response.Courses,
			Invites:   3,
		}
		buf, _ := json.Marshal(user)
		err := bucket.Put([]byte(chatID), buf)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		bot.Tg.Send(m.Sender, "Something wrong happended. Try again.", mainMenu)
		return
	}

	bot.Tg.Send(m.Sender, "Notifications are on", mainMenu)
}
