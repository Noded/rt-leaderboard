package db

// tasks its map of key string and value int 16
var tasks = map[string]int{
	"Cook":  10,
	"Study": 20,
}

// GetScores implemented method returns scores
func (s *SQLStorage) GetScores(userId int) (int, *[]int, error) {
	var userScore int
	err := s.db.QueryRow("SELECT score FROM leaderboard WHERE user_id=?", userId).Scan(&userScore)
	if err != nil {
		return 0, nil, err
	}

	rows, err := s.db.Query("SELECT score FROM leaderboard ORDER BY score DESC LIMIT 10")
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	topScores := make([]int, 0, 10) // Предварительное выделение памяти
	for rows.Next() {
		var score int
		if err := rows.Scan(&score); err != nil {
			return 0, nil, err
		}
		topScores = append(topScores, score)
	}

	return userScore, &topScores, nil
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
