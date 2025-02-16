package db

import (
	"database/sql"
	"fmt"
)

// tasks is a map containing the scores for each task.
var tasks = map[string]int{
	"Cook":  10,
	"Study": 20,
}

// UserScore holds the leaderboard data for a user.
type UserScore struct {
	Username string
	Score    int
	Rank     int
}

// GetUserRank retrieves the user's rank, username, and score.
// The rank is calculated with 1 being the highest score.
func (s *SQLStorage) GetUserRank(userId int) (int, string, int, error) {
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

// GetTopUsers retrieves the top 10 users with their scores and ranks.
func (s *SQLStorage) GetTopUsers() ([]UserScore, error) {
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

	topScores := make([]UserScore, 0, 10)

	for rows.Next() {
		var us UserScore
		if err := rows.Scan(&us.Username, &us.Score, &us.Rank); err != nil {
			return nil, fmt.Errorf("failed to scan user score: %w", err)
		}
		topScores = append(topScores, us)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return topScores, nil
}

// UpdateScore updates the user's score based on the specified task.
func (s *SQLStorage) UpdateScore(task string) error {
	value, exists := tasks[task]
	if !exists {
		return fmt.Errorf("invalid task: %s", task)
	}

	userId, err := GetCurrentUserID()
	if err != nil {
		return fmt.Errorf("failed to get current user id: %w", err)
	}

	// Begin a transaction.
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	var DBValue int
	err = tx.QueryRow("SELECT score FROM leaderboard WHERE user_id=?", userId).Scan(&DBValue)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get current score: %w", err)
	}

	// Add the task's score to the current score.
	value += DBValue

	_, err = tx.Exec("UPDATE leaderboard SET score=? WHERE user_id=?", value, userId)
	if err != nil {
		return fmt.Errorf("failed to update score: %w", err)
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
