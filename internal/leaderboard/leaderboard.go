package leaderboard

import (
	"rt-leaderboard/db"
)

var (
	data = db.NewSQLStorage()
)

// ShowBoard Shows top 10 users per score and user rank
func ShowBoard() (string, *[]string, error) {
	userID, err := db.GetCurrentUserID()
	if err != nil {
		return "", nil, err
	}
	userScore, topUsers, err := data.GetScores(userID)
	if err != nil {
		return "", nil, err
	}

	return userScore, topUsers, nil
}
