package sqlite

import (
	"BannerService/cmd/internal/storage"
	"fmt"

	"github.com/mattn/go-sqlite3"
)


func (s *Storage) CreateFeature(feature_name string) (int64, error) {
	const op = "storage.sqlite.CreateFeature"
	stmt, err := s.db.Prepare("INSERT INTO Feature(feature_name) VALUES(?)")
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	res, err := stmt.Exec(feature_name)
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

func (s *Storage) UpdateFeatureName(feature_id int64, feature_name string) error {
	const op = "storage.sqlite.UpdateFeatureName"
	stmt, err := s.db.Prepare("UPDATE Feature SET feature_name = ? WHERE banner_id =?")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	_, err = stmt.Exec(stmt, feature_name, feature_id)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil

}
func (s *Storage) DeleteFeature(feature_id int64) error {
	const op = "storage.sqlite.DeleteFeature"
	stmt, err := s.db.Prepare("DELETE FROM Tag WHERE tag_id = ?")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	_, err = stmt.Exec(stmt, feature_id)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s:%w", op, storage.ErrURLExists)
		}
		return nil
	}
	return nil
}
