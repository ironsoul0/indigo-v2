package main

import (
	"fmt"
	"log"

	"github.com/ironsoul0/indigo-v2/scrapers/moodle"
	"github.com/ironsoul0/indigo-v2/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Can not read config file:", err)
	}

	fmt.Println(config)

	moodleClient := moodle.Init()
	grades := moodleClient.GetGrades(config.UserName, config.UserPassword)
	fmt.Println(grades)

	// conn, err := sql.Open(config.DBDriver, config.DBSource)
	// if err != nil {
	// 	log.Fatal("Can not connect do DB:", err)
	// }

	// store := db.NewStore(conn)
}
