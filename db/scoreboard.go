package db

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
