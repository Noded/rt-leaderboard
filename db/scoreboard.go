package db

// tasks its map of key string and value int 16
var tasks = map[string]int16{
	"Cook":  10,
	"Study": 20,
}

// GetScores implemented method returns scores
func (s *SQLStorage) GetScores(userId int) (int, error) {
	var err error
	row := s.db.QueryRow("SELECT score FROM scoreboard WHERE user_id=?", userId)
	if row != nil {
		return 0, err
	}
	var scores int
	err = row.Scan(&scores)
	if err != nil {
		return 0, err
	}

	return scores, nil
}

// UpdateScore implemented method update sql table
func (s *SQLStorage) UpdateScore(task string) error {
	var value = tasks[task]

	userId, err := GetCurrentUserID()
	if err != nil {
		return err
	}

	s.db.Exec("UPDATE scoreboard SET score=? WHERE user_id=?", value, userId)

	return nil
}
