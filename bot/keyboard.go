package bot

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	schema "github.com/ironsoul0/indigo-v2/db"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	mainMenu    = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnStatus   = menu.Text("ℹ️  Status")
	btnRegister = menu.Text("🤖 Register")
	btnCodes    = menu.Text("🍧 Codes")
)

var (
	unverifiedMenu = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnVerify      = menu.Text("🗝  Verify")
)

func getChatInfo(db *bolt.DB, chatID int) (*schema.User, *schema.Scene) {
	var user *schema.User
	var scene *schema.Scene

	db.View(func(tx *bolt.Tx) error {
		chatID := fmt.Sprintf("%d", chatID)
		bucket := tx.Bucket([]byte(schema.USERS_BUCKET))
		userData := bucket.Get([]byte(chatID))

		if userData != nil {
			user = &schema.User{}
			json.Unmarshal(userData, user)
		}

		bucket = tx.Bucket([]byte(schema.STATES_BUCKET))
		sceneData := bucket.Get([]byte(chatID))

		if sceneData != nil {
			scene = &schema.Scene{}
			json.Unmarshal(sceneData, scene)
		}

		return nil
	})

	return user, scene
}

func initKeyboard(bot *tb.Bot, db *bolt.DB) {
	mainMenu.Reply(
		menu.Row(btnStatus),
		menu.Row(btnRegister),
		menu.Row(btnCodes),
	)

	unverifiedMenu.Reply(
		menu.Row(btnVerify),
	)

	bot.Handle(&btnStatus, handleStatus(bot, db))
	bot.Handle(&btnRegister, handleRegister(bot, db))
	bot.Handle(&btnCodes, handleRegister(bot, db))

	bot.Handle(tb.OnText, func(m *tb.Message) {
		user, scene := getChatInfo(db, m.Sender.ID)

		if user != nil {
			bot.Send(m.Sender, "", mainMenu)
			return
		}

		if scene.Scene == schema.REGISTER_SCENE {
			return
		}

		if scene.Scene == schema.VERIFICATION_SCENE {
			return
		}

		bot.Send(m.Sender, "Verify yourself", unverifiedMenu)
	})
}
