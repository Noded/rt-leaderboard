package db

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

// GetUserRankAndTopScores returns user rank, username, score and top 10 users with scores
func (s *SQLStorage) GetUserRankAndTopScores(userId int) (int, string, int, *[]UserScore, error) {
	// Query to get user's rank, username and score using window functions
	userQuery := `
        WITH RankedUsers AS (
            SELECT 
                username,
                score,
                RANK() OVER (ORDER BY score DESC) as rank
            FROM leaderboard l
            JOIN users u ON l.user_id = u.id
        )
        SELECT rank, username, score 
        FROM RankedUsers 
        WHERE user_id = ?`

	var userRank int
	var username string
	var userScore int

	err := s.db.QueryRow(userQuery, userId).Scan(&userRank, &username, &userScore)
	if err != nil {
		return 0, "", 0, nil, err
	}

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
		return 0, "", 0, nil, err
	}
	defer rows.Close()

	topScores := make([]UserScore, 0, 10)
	for rows.Next() {
		var us UserScore
		if err := rows.Scan(&us.Username, &us.Score, &us.Rank); err != nil {
			return 0, "", 0, nil, err
		}
		topScores = append(topScores, us)
	}

	return userRank, username, userScore, &topScores, nil
}

// UpdateScore implemented method update sql table
func (s *SQLStorage) UpdateScore(task string) error {
	var value = tasks[task]
	var DBValue int

	userId, err := GetCurrentUserID()
	if err != nil {
		return err
	}
	s.db.QueryRow("SELECT score FROM leaderboard WHERE user_id=?", userId).Scan(&DBValue)

	value += DBValue

	s.db.Exec("UPDATE leaderboard SET score=? WHERE user_id=?", value, userId)

	return nil
}
