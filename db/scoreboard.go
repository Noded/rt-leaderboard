package db

import (
	"database/sql"
	"log"
)

// tasks its map of key string and value int 16
var tasks = map[string]int16{
	"Cook":  10,
	"Study": 20,
}

// AddScore method adding score to db
func (s *SQLStorage) AddScore(task string) error {
	var err error
	if task != "" {
		var value = tasks[task]
		_, err = s.db.Exec("INSERT INTO scoreboard(score) VALUES(?)", value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLStorage) GetScores(userId string) (*[]int, error) {
	var err error
	if userId == "" {
		// Returns nothing
		return nil, nil
	}
	rows, err := s.db.Query("SELECT score FROM scoreboard WHERE user_id=?", userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var scores []int
	for rows.Next() {
		var score int
		if err := rows.Scan(&score); err != nil {
			log.Fatal(err)
		}
		scores = append(scores, score)
	}
	scoresPtr := &scores

	return scoresPtr, nil
}
