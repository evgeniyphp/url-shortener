package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}
	
	query, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, err
	}
	
	_, err = query.Exec()
	if err != nil {
		return nil, err
	}
	
	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(url, alias string) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO url(alias, url) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	
	res, err := stmt.Exec(alias, url)
	if err != nil {
		return 0, err
	}
	
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	stmt, err := s.db.Prepare("SELECT * FROM url WHERE alias = ?")
	if err != nil {
		return "", err
	}
	
	var url string
	
	row := stmt.QueryRow(alias)
	err = row.Scan(&url)
	
	return url, nil
}

// TODO: implement delete method
//func (s *Storage) DeleteUrl(alias string) (string, error) {
//	
//}
