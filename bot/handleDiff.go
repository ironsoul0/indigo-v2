package bot

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	schema "github.com/ironsoul0/indigo-v2/db"
	tb "gopkg.in/tucnak/telebot.v2"
)

const newGradeNotification = `
New grade!

Course: <b>%s</b>
Name: <b>%s</b>
Grade: <b>%s</b>
Range: <b>%s</b>
Percentage: <b>%s</b>
`

func handleDiff(bot *tb.Bot, db *bolt.DB, parseResult ParseResult) {
	if parseResult.deactivate {
		bot.Send(tb.ChatID(parseResult.chatID), "Your credentials are wrong. Notifications are off for you.")

		db.Update(func(tx *bolt.Tx) error {
			chatID := fmt.Sprintf("%d", parseResult.chatID)
			bucket := tx.Bucket([]byte(schema.USERS_BUCKET))
			userData := bucket.Get([]byte(chatID))

			var user schema.User
			json.Unmarshal(userData, &user)
			user.Activated = false
			buf, _ := json.Marshal(user)

			bucket.Put([]byte(chatID), buf)
			return nil
		})

		return
	}

	newGrades := 0
	for _, courseDiff := range parseResult.diff {
		for _, newGrade := range courseDiff.Grades {
			bot.Send(
				tb.ChatID(parseResult.chatID),
				fmt.Sprintf(
					newGradeNotification,
					courseDiff.Name,
					newGrade.Name,
					newGrade.Grade,
					newGrade.Range,
					newGrade.Percentage,
				),
				tb.ModeHTML,
			)

			newGrades++
		}
	}

	if newGrades > 0 {
		db.Update(func(tx *bolt.Tx) error {
			chatID := fmt.Sprintf("%d", parseResult.chatID)
			bucket := tx.Bucket([]byte(schema.USERS_BUCKET))
			userData := bucket.Get([]byte(chatID))

			var user schema.User
			json.Unmarshal(userData, &user)
			user.Courses = parseResult.recentCoruses
			buf, _ := json.Marshal(user)

			bucket.Put([]byte(chatID), buf)
			return nil
		})
	}
}
