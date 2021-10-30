package bot

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/ironsoul0/indigo-v2/db"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	mainMenu    = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnStatus   = menu.Text("ℹ️  Status")
	btnRegister = menu.Text("🤖  Register")
	btnCodes    = menu.Text("🍧  Codes")
)

var (
	unverifiedMenu = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnVerify      = menu.Text("🗝  Verify")
)

func (bot *Bot) initKeyboard() {
	mainMenu.Reply(
		menu.Row(btnStatus),
		menu.Row(btnRegister),
		menu.Row(btnCodes),
	)

	unverifiedMenu.Reply(
		menu.Row(btnVerify),
	)

	bot.Tg.Handle(&btnStatus, bot.handleStatus)
	bot.Tg.Handle(&btnRegister, bot.handleRegister)
	bot.Tg.Handle(&btnCodes, bot.handleCodes)

	bot.Tg.Handle(&btnVerify, bot.handleVerify)

	bot.Tg.Handle(tb.OnText, func(m *tb.Message) {
		user, scene := bot.getChatInfo(m.Sender.ID)

		if scene != nil && scene.Scene == db.REGISTER_SCENE {
			bot.handleRegisterStep(m, scene)
			return
		}

		if scene != nil && scene.Scene == db.VERIFICATION_SCENE {
			bot.handleVerifyStep(m)
			return
		}

		if user != nil {
			bot.Tg.Send(m.Sender, "Please choose command below", mainMenu)
			return
		}

		bot.Tg.Send(m.Sender, "Please choose command below", unverifiedMenu)
	})
}

func (bot *Bot) addAdmin() {
	bot.db.Update(func(tx *bolt.Tx) error {
		u := db.User{
			Username: "",
			Password: "",
			ChatID:   1,
			Invites:  3,
		}
		buf, _ := json.Marshal(u)
		bucket := tx.Bucket([]byte(db.USERS_BUCKET))
		bucket.Put([]byte("1"), buf)

		return nil
	})
}
