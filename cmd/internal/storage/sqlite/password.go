package sqlite

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *Storage) CheckPassword(name string) (string, string, string, string, error) {
	const op = "storage.sqlite.CheckPassword"
	var password_hash, salt, isAdmin, useLastRevision string
	stmt, err := s.db.Prepare("SELECT password_hash, salt, is_admin, use_last_revision FROM Users WHERE user_name = ?")
	if err != nil {
		return "", "", "", "", fmt.Errorf("%s:%w", op, err)
	}
	err = stmt.QueryRow(name).Scan(&password_hash, &salt, &isAdmin, &useLastRevision)
	fmt.Println(password_hash, salt, isAdmin, useLastRevision)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", "", fmt.Errorf("user not found")
		}
		return "", "", "", "", fmt.Errorf("db error")
	}
	return password_hash, salt, isAdmin, useLastRevision, nil
}

func (s *Storage) CreateUser(username string, password_hash string, salt string, isAdmin string, useLastRevision string) error {
	const op = "storage.sqlite.CreateUser"
	query := "SELECT user_name FROM Users WHERE user_name = ?"
	var res string
	log.Println("check if user exists", username)
	err := s.db.QueryRow(query, username).Scan(&res)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("пользователь существует %s:%w", op, err)
	}

	if res != "" {
		return fmt.Errorf("пользователь существует %s:%w", op, err)
	}
	log.Println("Inserting new user:", username)

	query = "INSERT INTO Users(user_name, password_hash, salt, is_admin, use_last_revision) VALUES (?,?,?,?,?)"
	_, err = s.db.Exec(query, username, password_hash, salt, isAdmin, useLastRevision)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: %w", op, err)
		} else {
			return fmt.Errorf("%s:%w", op, err)
		}

	}
	return nil
}

func (s *Storage) ChangePassword(username string, password_hash string, salt string) error {
	return nil
}
