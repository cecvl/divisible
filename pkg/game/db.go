package game

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

type Stats struct {
	Plays     int
	Total     int
	LastScore int
	BestScore int
}

func NewStore(path string) (*Store, error) {
	d, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	// set some sensible pragmas
	_, _ = d.Exec("PRAGMA journal_mode=WAL;")
	_, _ = d.Exec("PRAGMA synchronous=NORMAL;")

	// create tables
	schema := []string{
		`CREATE TABLE IF NOT EXISTS scores (id INTEGER PRIMARY KEY AUTOINCREMENT, score INTEGER NOT NULL, created_at DATETIME DEFAULT CURRENT_TIMESTAMP);`,
		`CREATE TABLE IF NOT EXISTS stats (id INTEGER PRIMARY KEY CHECK(id=1), plays INTEGER, total_score INTEGER, last_score INTEGER, best_score INTEGER);`,
		`INSERT OR IGNORE INTO stats(id, plays, total_score, last_score, best_score) VALUES (1,0,0,0,0);`,
	}
	for _, q := range schema {
		if _, err := d.Exec(q); err != nil {
			d.Close()
			return nil, fmt.Errorf("creating schema: %w", err)
		}
	}

	return &Store{db: d}, nil
}

func (s *Store) SaveScore(score int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`INSERT INTO scores(score) VALUES (?);`, score); err != nil {
		return err
	}

	if _, err := tx.Exec(`UPDATE stats SET plays = plays + 1, total_score = total_score + ?, last_score = ?, best_score = CASE WHEN ? > best_score THEN ? ELSE best_score END WHERE id = 1;`, score, score, score, score); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) GetStats() (Stats, error) {
	var st Stats
	row := s.db.QueryRow(`SELECT plays, total_score, last_score, best_score FROM stats WHERE id = 1;`)
	if err := row.Scan(&st.Plays, &st.Total, &st.LastScore, &st.BestScore); err != nil {
		return st, err
	}
	return st, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
