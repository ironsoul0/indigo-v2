package main

import (
	"fmt"
	"time"

	"github.com/ironsoul0/indigo-v2/bot"
	"github.com/ironsoul0/indigo-v2/db"
	"github.com/ironsoul0/indigo-v2/util"

	_ "github.com/lib/pq"
)

func main() {
	config := util.LoadConfig(".")
	db := db.New(config.DBPath)
	bot := bot.New(config.BotToken, db)

	fmt.Println(bot)

	time.Sleep(100 * time.Second)
}
