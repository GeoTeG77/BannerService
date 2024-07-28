package sqlite

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	creators := make([]string, 0, 8)

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ReferencesON := `PRAGMA foreign_keys = ON;`
	CreateTagBannerTable := `CREATE TABLE IF NOT EXISTS TagBanner(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tag_id INTEGER,
		banner_id INTEGER
		);`
	CreateTagTable := `CREATE TABLE IF NOT EXISTS Tag(
		tag_id INTEGER PRIMARY KEY AUTOINCREMENT,
		tag_name TEXT NOT NULL,
		FOREIGN KEY (tag_id) REFERENCES TagBanner(tag_id) ON DELETE CASCADE ON UPDATE CASCADE
		);`
	CreateContentTable := `CREATE TABLE IF NOT EXISTS Content(
		banner_id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		text TEXT,
		url TEXT,
		version TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		is_active INTEGER NOT NULL CHECK (is_active IN (0, 1))
		);`
	CreateBannerTable := `CREATE TABLE IF NOT EXISTS Banner(
		banner_id INTEGER PRIMARY KEY,
		feature_id INTEGER,
		FOREIGN KEY (banner_id) REFERENCES Content(banner_id) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (banner_id) REFERENCES TagBanner(banner_id) ON DELETE CASCADE ON UPDATE CASCADE
		);`
	CreateFeatureTable := `CREATE TABLE IF NOT EXISTS Feature(
		feature_id INTEGER PRIMARY KEY,
		feature_name TEXT NOT NULL,
		FOREIGN KEY (feature_id) REFERENCES Banner(feature_id) ON DELETE CASCADE ON UPDATE CASCADE
		);`
	CreateUsersTable := `CREATE TABLE IF NOT EXISTS Users(
    user_name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
	salt TEXT NOT NULL,
	is_admin INTEGER NOT NULL CHECK (is_admin IN (0, 1)),
	use_last_revision INTEGER NOT NULL CHECK (use_last_revision IN (0, 1))
	);`
	
	creators = append(creators, ReferencesON, CreateTagBannerTable, CreateTagTable, CreateContentTable, CreateBannerTable, CreateFeatureTable, CreateUsersTable)
	for _, v := range creators {
		_, err = db.Exec(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		fmt.Println("success")
	}

	return &Storage{db: db}, nil
}


func (s *Storage) Close() error {
	return s.db.Close()
}