package bot

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/ironsoul0/indigo-v2/db"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (bot *Bot) handleVerify(m *tb.Message) {
	bot.updateScene(m.Sender.ID, db.VERIFICATION_SCENE)
	bot.Tg.Send(m.Sender, "Do you have what it takes?..", &tb.ReplyMarkup{ReplyKeyboardRemove: true})
}

func (bot *Bot) handleVerifyStep(m *tb.Message) {
	chatID, err := strconv.Atoi(m.Text)
	defer bot.updateScene(m.Sender.ID, db.NO_SCENE)

	if err != nil {
		bot.Tg.Send(m.Sender, "You sent the wrong secret phrase! 1", unverifiedMenu)
		return
	}
	user, _ := bot.getChatInfo(chatID)
	if user == nil || user.Invites == 0 {
		bot.Tg.Send(m.Sender, "You sent the wrong secret phrase! 2", unverifiedMenu)
		return
	}

	bot.Tg.Send(m.Sender, "You are in!", mainMenu)
	bot.db.Update(func(tx *bolt.Tx) error {
		user.Invites--
		bucket := tx.Bucket([]byte(db.USERS_BUCKET))
		buf, _ := json.Marshal(*user)
		bucket.Put([]byte(m.Text), buf)

		newUser := &db.User{ChatID: m.Sender.ID}
		buf, _ = json.Marshal(*newUser)
		bucket.Put([]byte(fmt.Sprintf("%d", m.Sender.ID)), buf)

		return nil
	})
}
