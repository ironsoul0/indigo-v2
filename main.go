package main

import (
	"database/sql"
	"log"

	db "github.com/ironsoul0/indigo-v2/db/sqlc"
	"github.com/ironsoul0/indigo-v2/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Can not read config file:", err)
	}

	// moodleClient := moodle.Init()
	// _ = moodleClient.GetGrades(config.UserName, config.UserPassword)

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can not connect do DB:", err)
	}

	store := db.NewStore(conn)
}
