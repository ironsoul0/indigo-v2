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

	// moodleClient := moodle.Init()
	// grades := moodleClient.GetGrades(config.UserName, config.UserPassword)

	// conn, err := sql.Open(config.DBDriver, config.DBSource)
	// if err != nil {
	// 	log.Fatal("Can not connect do DB:", err)
	// }

	// store := db.NewStore(conn)
}
