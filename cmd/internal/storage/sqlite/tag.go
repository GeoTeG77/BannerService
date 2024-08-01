package sqlite

import (
	"BannerService/cmd/internal/storage"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func (s *Storage) CreateTag(tag_name string) (int64, error) {
	const op = "storage.sqlite.CreateTag"
	stmt, err := s.db.Prepare("INSERT INTO Tag(tag_name) VALUES(?)")
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	res, err := stmt.Exec(tag_name)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s:%w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	return id, nil
}
// нужно переделать!!! апдейт тэг
func (s *Storage) UpdateTag(banner_id *int64, tag_ids []int64) error {
	const op = "storage.sqlite.UpdateBanner.Tags_ids"
	stmt, err := s.db.Prepare("UPDATE BannerTag SET tag_ids = ? WHERE banner_id = ?")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	_, err = stmt.Exec(op, tag_ids, banner_id)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s:%w", op, storage.ErrURLExists)
		}
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}
func (s *Storage) UpdateTagName(tag_id int64, tag_name string) error {
	const op = "storage.sqlite.UpdateTagName"
	stmt, err := s.db.Prepare("UPDATE Tag SET tag_name = ? WHERE banner_id =?")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	_, err = stmt.Exec(stmt, tag_name, tag_id)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}
func (s *Storage) DeleteTag(tag_id int64) error {
	const op = "storage.sqlite.DeleteTag"
	stmt, err := s.db.Prepare("DELETE FROM Tag WHERE tag_id = ?")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	_, err = stmt.Exec(stmt, tag_id)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s:%w", op, storage.ErrURLExists)
		}
		return nil
	}
	return nil
}

