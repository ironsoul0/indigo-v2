package bot

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/ironsoul0/indigo-v2/db"
)

func (bot *Bot) updateScene(chatID int, scene db.SceneID) {
	bot.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(db.STATES_BUCKET))

		state := db.Scene{
			ChatID: chatID,
			Scene:  scene,
			Step:   0,
		}
		buf, _ := json.Marshal(state)
		bucket.Put([]byte(fmt.Sprintf("%d", chatID)), buf)

		return nil
	})
}

func (bot *Bot) getChatInfo(chatID int) (*db.User, *db.Scene) {
	var user *db.User
	var scene *db.Scene

	bot.db.View(func(tx *bolt.Tx) error {
		chatID := fmt.Sprintf("%d", chatID)
		bucket := tx.Bucket([]byte(db.USERS_BUCKET))
		userData := bucket.Get([]byte(chatID))

		if userData != nil {
			user = &db.User{}
			json.Unmarshal(userData, user)
		}

		bucket = tx.Bucket([]byte(db.STATES_BUCKET))
		sceneData := bucket.Get([]byte(chatID))

		if sceneData != nil {
			scene = &db.Scene{}
			json.Unmarshal(sceneData, scene)
		}

		return nil
	})

	return user, scene
}
