package bot

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	schema "github.com/ironsoul0/indigo-v2/db"
	"github.com/ironsoul0/indigo-v2/scrapers/moodle"
)

type ParseResult struct {
	chatID        int
	diff          []moodle.Course
	recentCoruses []moodle.Course
	deactivate    bool
}

const (
	WORKERS = 10
)

func checkGrades(taskPool chan schema.User, resultChan chan ParseResult, wg *sync.WaitGroup) {
	defer wg.Done()

	moodleClient := moodle.Init()

	for userToCheck := range taskPool {
		response := moodleClient.GetGrades(userToCheck.Username, userToCheck.Password)

		if response.InvalidCredentials {
			resultChan <- ParseResult{
				chatID:     userToCheck.ChatID,
				deactivate: true,
			}
			return
		}

		if !response.Success {
			continue
		}

		resultChan <- ParseResult{
			chatID:        userToCheck.ChatID,
			diff:          moodle.DetectNewGrades(userToCheck.Courses, response.Courses),
			recentCoruses: response.Courses,
		}
	}
}

func (bot *Bot) notify() {
	for {
		usersToCheck := make([]schema.User, 0)

		bot.db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(schema.USERS_BUCKET))
			c := bucket.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				var user schema.User
				json.Unmarshal(v, &user)

				if user.Activated {
					usersToCheck = append(usersToCheck, user)
				}
			}

			return nil
		})

		taskPool := make(chan schema.User)
		resultChan := make(chan ParseResult)
		var wg sync.WaitGroup

		wg.Add(WORKERS)
		for i := 0; i < WORKERS; i++ {
			go checkGrades(taskPool, resultChan, &wg)
		}

		for _, user := range usersToCheck {
			taskPool <- user
		}
		close(taskPool)

		go func() {
			wg.Wait()
			close(resultChan)
		}()

		for moodleDiff := range resultChan {
			bot.handleDiff(moodleDiff)
		}

		time.Sleep(15 * time.Second)
	}
}
