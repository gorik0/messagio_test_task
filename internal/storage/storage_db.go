package storage

func (s *Storage) CloseDB() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddRecord(message string, data string) (int64, error) {
	stmt := `INSERT INTO messages (message,data) values ($1, $2) returning id`
	var id int
	err := s.db.QueryRow(stmt, message, data).Scan(&id)
	if err != nil {
		return 0, err
	}
	return int64(id), err
}

func (s *Storage) MessageConsumed(id int64) error {
	stmt := `UPDATE messages set processed=true where id=$1`
	_, err := s.db.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil

}

func (s *Storage) NumberOfProcessedMessages() (total int, processed int, err error) {

	stmt := `SELECT 
    (SELECT COUNT(*) FROM messages) as total,
    (SELECT COUNT(*) FROM messages where processed = true) as processed`

	err = s.db.QueryRow(stmt).Scan(&total, &processed)
	if err != nil {
		return 0, 0, err

	}
	return total, processed, nil
}
