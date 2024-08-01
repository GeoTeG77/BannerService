package sqlite

import (
	banner "BannerService/cmd/internal/http-server/handlers/url/banner"
	"database/sql"
	//"BannerService/cmd/internal/http-server/handlers/url/feature"
	"BannerService/cmd/internal/storage"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func (s *Storage) CreateBanner(feature_id int64, tag_ids []int64, title string, content string, url string, is_active string) (int64, error) {
	const op = "storage.sqlite.CreateBanner"
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	res, err := s.CreateBannerContentStmt.Exec(title, content, url)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s:%w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	_, err = s.CreateBannerFeatureStmt.Exec(id, feature_id)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s:%w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	//INSERT INTO TagBanner TABLE
	for _, tag_id := range tag_ids {
		_, err = s.CreateTagBannerStmt.Exec(id, tag_id)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, fmt.Errorf("%s:%w", op, storage.ErrURLExists)
			}
			return 0, fmt.Errorf("%s:%w", op, err)
		}
	}

	if err = tx.Commit(); err != nil {
		fmt.Println("Error committing transaction:", err)
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}

func (s *Storage) UpdateBanner(banner_id int64, tag_ids []int64, feature_id int64, title string, text string, url string, is_active string) error {
	const op = "storage.sqlite.UpdateBanner"
	tx, err := s.db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	switch {
	case len(tag_ids) != 0:

		s.UpdateTag(&banner_id, tag_ids)
		fallthrough

	case feature_id != 0:

		_, err = s.UpdatBannerFeatureStmt.Exec(banner_id, feature_id)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return fmt.Errorf("%s:%w", op, err)
			}
			return fmt.Errorf("%s:%w", op, err)
		}
		fallthrough

	case title != "":

		_, err = s.UpdateBannerContentTitleStmt.Exec(title, banner_id)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return fmt.Errorf("%s:%w", op, err)
			}
			return fmt.Errorf("%s:%w", op, err)
		}
		fallthrough

	case text != "":

		_, err = s.UpdateBannerContentTextStmt.Exec(text, banner_id)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return fmt.Errorf("%s:%w", op, storage.ErrURLExists)
			}
			return fmt.Errorf("%s:%w", op, err)
		}
		fallthrough

	case url != "":

		_, err = s.UpdateBannerContentURLStmt.Exec(url, banner_id)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return fmt.Errorf("%s:%w", op, storage.ErrURLExists)
			}
			return fmt.Errorf("%s:%w", op, err)
		}

	default:
		fmt.Println("unreachable sit")
	}

	_, err = s.UpdateBannerIsActiveStmt.Exec(is_active, banner_id)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s:%w", op, err)
		}
		return fmt.Errorf("%s:%w", op, err)
	}

	if err = tx.Commit(); err != nil {
		fmt.Println("Error committing transaction:", err)
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (s *Storage) GetBanner(feature_id int64, tag_id int64) (string, string, string, string, error) {
	const op = "storage.sqlite.GetBanner"

	tx, err := s.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var title, text, url, is_active string

	err = s.GetBannerStmt.QueryRow(tag_id, feature_id).Scan(&title, &text, &url, &is_active)
	if err != nil {
		tx.Rollback()
	}

	if err = tx.Commit(); err != nil {
		fmt.Println("Error committing transaction:", err)
		return "", "", "", "", fmt.Errorf("%s:%w", op, err)
	}
	return title, text, url, is_active, nil
}

func (s *Storage) GetBanners(feature_id int64, tag_id int64, limit int64, offset int64) ([]banner.GetBannersResponce, error) {
	const op = "storage.sqlite.GetBanners"
	_ = op //закинуть логгер

	var rows *sql.Rows
	var err error
	var res banner.GetBannersResponce
	ress := make([]banner.GetBannersResponce, 0, limit)
	
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	switch {
	case feature_id !=0 && tag_id == 0:
		rows,err = s.GetBannersByFeatureStmt.Query(feature_id, limit, offset)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return nil, fmt.Errorf("%s:%w", op, err)
			}
			return nil, fmt.Errorf("%s:%w", op, err)
		}
	defer rows.Close()

	case feature_id ==0 && tag_id !=0:
	rows,err = s.GetBannersByTagStmt.Query(tag_id, limit,offset)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s:%w", op, err)
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	
	default:
		rows,err = s.GetBannersByTagStmt.Query(tag_id, limit,offset)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return nil, fmt.Errorf("%s:%w", op, err)
			}
			return nil, fmt.Errorf("%s:%w", op, err)
		}
	}

	for rows.Next() {
		if err := rows.Scan(&res.Banner_id, &res.Tag_ids, &res.Feature_id, &res.Content.Title, &res.Content.Text, &res.Content.URL, &res.Is_active, &res.Created_at, &res.Updated_at); err != nil {
			panic(err)
		}
		ress = append(ress, res)
	}

	if err = tx.Commit(); err != nil {
		fmt.Println("Error committing transaction:", err)
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return ress, nil
}

func (s *Storage) DeleteBanner(banner_id int64) (error) {
	const op = "storage.sqlite.DeleteBanner"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()


	stmt, err := s.db.Prepare("DELETE FROM Content WHERE banner_id = ?")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	_, err = stmt.Exec(stmt, banner_id)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s:%w", op, storage.ErrURLExists)
		}
		return fmt.Errorf("%s:%w", op, storage.ErrURLExists)
	}
	
	if err = tx.Commit(); err != nil {
		fmt.Println("Error committing transaction:", err)
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}
