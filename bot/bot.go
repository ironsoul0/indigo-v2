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

	go bot.Start()
	go notify(bot, db)

	return bot
}
