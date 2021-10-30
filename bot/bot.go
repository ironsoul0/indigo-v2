package bot

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	Tg *tb.Bot
	db *bolt.DB
}

func New(token string, db *bolt.DB) *Bot {
	tgBot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	bot := &Bot{Tg: tgBot, db: db}
	bot.addAdmin()
	bot.initKeyboard()

	go bot.Tg.Start()
	go bot.notify()

	return bot
}
