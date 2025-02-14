package db

import (
	"database/sql"
	"fmt"
)

// tasks its map of key string and value int 16
var tasks = map[string]int{
	"Cook":  10,
	"Study": 20,
}

// UserScore Structure to hold user score data
type UserScore struct {
	Username string
	Score    int
	Rank     int
}

// GetUserRank retrieves user's rank, username and score
// Returns rank starting from 1 (highest score)
func (s *SQLStorage) GetUserRank(userId int) (int, string, int, error) {
	// Single query to get rank and user data
	rankQuery := `
        SELECT 
            (SELECT COUNT(*) + 1 
             FROM leaderboard l2 
             WHERE l2.score > l1.score) as rank,
            u.username,
            l1.score
        FROM leaderboard l1
        JOIN users u ON l1.user_id = u.id
        WHERE l1.user_id = ?`

	var userRank int
	var username string
	var userScore int

	err := s.db.QueryRow(rankQuery, userId).Scan(&userRank, &username, &userScore)
	if err == sql.ErrNoRows {
		return 0, "", 0, fmt.Errorf("user with id %d not found", userId)
	}
	if err != nil {
		return 0, "", 0, fmt.Errorf("failed to get user rank: %w", err)
	}

	return userRank, username, userScore, nil
}

// GetTopUsers retrieves top 10 users with their scores and ranks
func (s *SQLStorage) GetTopUsers() ([]UserScore, error) {
	// Query to get top 10 users with ranks
	topQuery := `
        SELECT 
            u.username,
            l.score,
            RANK() OVER (ORDER BY l.score DESC) as rank
        FROM leaderboard l
        JOIN users u ON l.user_id = u.id
        ORDER BY l.score DESC
        LIMIT 10`

	rows, err := s.db.Query(topQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query top users: %w", err)
	}
	defer rows.Close()

	// Pre-allocate slice with capacity of 10
	topScores := make([]UserScore, 0, 10)

	for rows.Next() {
		var us UserScore
		if err := rows.Scan(&us.Username, &us.Score, &us.Rank); err != nil {
			return nil, fmt.Errorf("failed to scan user score: %w", err)
		}
		topScores = append(topScores, us)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return topScores, nil
}

// UpdateScore implemented method update sql table
func (s *SQLStorage) UpdateScore(task string) error {
	if task != "Cook" && task != "Study" {
		return fmt.Errorf("invalid task: %s", task)
	}
	var value = tasks[task]
	var DBValue int
	userId, err := GetCurrentUserID()
	if err != nil {
		return fmt.Errorf("failed to get current user id: %w", err)
	}
	s.db.QueryRow("SELECT score FROM leaderboard WHERE user_id=?", userId).Scan(&DBValue)

	value += DBValue

	s.db.Exec("UPDATE leaderboard SET score=? WHERE user_id=?", value, userId)

	return nil
}
