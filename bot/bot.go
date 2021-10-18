package bot

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	WORKERS = 10
)

func New(token string, db *bolt.DB) *tb.Bot {
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/register", handleRegister(bot, db))
	bot.Handle("/status", handleStatus(bot, db))

	var (
		// Universal markup builders.
		menu     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		selector = &tb.ReplyMarkup{}

		// Reply buttons.
		btnHelp     = menu.Text("ℹ Help")
		btnSettings = menu.Text("⚙ Settings")

		// Inline buttons.
		//
		// Pressing it will cause the client to
		// send the bot a callback.
		//
		// Make sure Unique stays unique as per button kind,
		// as it has to be for callback routing to work.
		//
		btnPrev = selector.Data("⬅", "prev")
		btnNext = selector.Data("➡", "next")
	)

	menu.Reply(
		menu.Row(btnHelp),
		menu.Row(btnSettings),
	)
	selector.Inline(
		selector.Row(btnPrev, btnNext),
	)

	bot.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			return
		}

		bot.Send(m.Sender, "Hello!", menu)
	})

	bot.Handle(&btnHelp, func(m *tb.Message) {
		bot.Send(m.Sender, "Lol!")
	})

	// On inline button pressed (callback)
	bot.Handle(&btnPrev, func(c *tb.Callback) {
		// ...
		// Always respond!
		bot.Respond(c, &tb.CallbackResponse{})
	})

	go bot.Start()
	go notify(bot, db)

	return bot
}
