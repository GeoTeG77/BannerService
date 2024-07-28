package sqlite

import (
	banner "BannerService/cmd/internal/http-server/handlers/url/banner"
	"BannerService/cmd/internal/storage"
	"fmt"

	"github.com/mattn/go-sqlite3"
)



func (s *Storage) CreateBanner(req banner.CreateRequest) (int64, error) {
	const op = "storage.sqlite.CreateBanner"

	//INSERT INTO Content TABLE
	stmt, err := s.db.Prepare("INSERT INTO Content(title, text, url) VALUES(?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	res, err := stmt.Exec(req.Content.Title, req.Content.Text, req.Content.URL)
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

	//INSERT INTO Banner TABLE
	stmt, err = s.db.Prepare("INSERT INTO Banner(banner_id,feature_id) VALUES(?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	_, err = stmt.Exec(id, req.Feature_id)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s:%w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	//INSERT INTO TagBanner TABLE
	stmt, err = s.db.Prepare("INSERT INTO TagBanner(tag_id,banner_id) VALUES(?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	for _, v := range req.Tag_ids {
		_, err = stmt.Exec(v, id)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, fmt.Errorf("%s:%w", op, storage.ErrURLExists)
			}
			return 0, fmt.Errorf("%s:%w", op, err)
		}
	}

	return id, nil
}

func (s *Storage) UpdateBanner(req banner.UpdateRequest) error {
	const op = "storage.sqlite.UpdateBanner"

	switch {
	case req.Tag_ids != nil:
		s.UpdateTag(req.Banner_id, req.Tag_ids)
		fallthrough

	case req.Feature_id != nil:
		s.UpdateFeature(req.Banner_id, req.Feature_id)
		fallthrough

	//INSERT INTO Content TABLE
	case req.Content.Title != nil:

		stmt, err := s.db.Prepare("UPDATE Content SET title = ? WHERE banner_id = ?")

		if err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}

		_, err = stmt.Exec(req.Content.Title, req.Banner_id)

		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return fmt.Errorf("%s:%w", op, storage.ErrURLExists)
			}
			return fmt.Errorf("%s:%w", op, err)
		}
		fallthrough

	case req.Content.Text != nil:

		stmt, err := s.db.Prepare("UPDATE Content SET text = ? WHERE banner_id = ?")

		if err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}

		_, err = stmt.Exec(req.Content.Text, req.Banner_id)

		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return fmt.Errorf("%s:%w", op, storage.ErrURLExists)
			}
			return fmt.Errorf("%s:%w", op, err)
		}
		fallthrough

	case req.Content.URL != nil:

		stmt, err := s.db.Prepare("UPDATE Content SET url = ? WHERE banner_id = ?")

		if err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}

		_, err = stmt.Exec(req.Content.URL, req.Banner_id)

		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return fmt.Errorf("%s:%w", op, storage.ErrURLExists)
			}
			return fmt.Errorf("%s:%w", op, err)
		}
		fallthrough
	case req.Is_active != nil:
		stmt, err := s.db.Prepare("UPDATE Content SET is_active = ? WHERE banner_id = ?")

		if err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}

		_, err = stmt.Exec(req.Is_active, req.Banner_id)

		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return fmt.Errorf("%s:%w", op, storage.ErrURLExists)
			}
			return fmt.Errorf("%s:%w", op, err)
		}
		return nil
	default:
		fmt.Println("unreachable sit")
	}
	return nil
}

func (s *Storage) GetBanner(req banner.GetBannerRequest) (banner.Content, error) {

	query := `
        SELECT c.title, c.text, c.url
        FROM Content c
        INNER JOIN Banner b ON c.banner_id = b.banner_id
        INNER JOIN TagBanner tb ON b.banner_id = tb.banner_id
        WHERE tb.tag_id = ? AND b.feature_id = ?;
    `
	stmt, err := s.db.Query(query, req.Tag_id, req.Feature_id)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var res banner.Content

	var title, text, url string

	for stmt.Next() {
		if err := stmt.Scan(&title, &text, &url); err != nil {
			fmt.Println("db error")
		}
	}

	res.Title = &title
	res.Text = &text
	res.URL = &url

	return res, nil

}

func (s *Storage) GetBanners(req banner.GetBannersRequest) ([]banner.GetBannersResponce, error) {
	const op = "storage.sqlite.GetBanners"
	_ = op //закинуть логгер
	if req.Offset == nil {
		*req.Offset = 1
	}

	if req.Limit == nil {
		*req.Limit = 10
	}

	var res banner.GetBannersResponce
	ress := make([]banner.GetBannersResponce, 0, *req.Limit)

	switch {
	case req.Tag_id == nil && req.Feature_id == nil:
		query := `SELECT
    b.banner_id,
    GROUP_CONCAT(tb.tag_id) AS tag_ids,
    b.feature_id,
    c.title,
    c.text,
    c.url,
    c.version,
    c.updated_at,
    c.is_active
	FROM Banner b
	JOIN Content c ON b.banner_id = c.banner_id
	LEFT JOIN TagBanner tb ON b.banner_id = tb.banner_id
	GROUP BY b.banner_id, b.feature_id, c.title, c.text, c.url, c.version, c.updated_at, c.is_active
	LIMIT ? OFFSET ?;`

		rows, err := s.db.Query(query, req.Limit, req.Offset)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&res.Banner_id, &res.Tag_ids, &res.Feature_id, &res.Content.Title, &res.Content.Text, &res.Content.URL, &res.Is_active, &res.Created_at, &res.Updated_at); err != nil {
				panic(err)
			}
			ress = append(ress, res)
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}

	case req.Feature_id != nil && req.Tag_id == nil:
		query := `SELECT
    b.banner_id,
    GROUP_CONCAT(tb.tag_id) AS tag_ids,
    c.title,
    c.text,
    c.url,
    c.version,
    c.updated_at,
    c.is_active
	FROM Banner b
	JOIN Content c ON b.banner_id = c.banner_id WHERE b.feature_id = ?
	LEFT JOIN TagBanner tb ON b.banner_id = tb.banner_id WHERE b.feature_id = ?
	GROUP BY b.banner_id, b.feature_id, c.title, c.text, c.url, c.version, c.updated_at, c.is_active
	LIMIT ? OFFSET ?;`

		rows, err := s.db.Query(query, req.Feature_id, req.Limit, req.Offset)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&res.Banner_id, &res.Tag_ids, &res.Feature_id, &res.Content.Title, &res.Content.Text, &res.Content.URL, &res.Is_active, &res.Created_at, &res.Updated_at); err != nil {
				panic(err)
			}
			ress = append(ress, res)
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	case req.Feature_id == nil && req.Tag_id != nil:
		query := `SELECT
    b.banner_id,
    GROUP_CONCAT(tb.tag_id) AS tag_ids,
    b.feature_id,
    c.title,
    c.text,
    c.url,
    c.version,
    c.updated_at,
    c.is_active
	FROM Banner b
	JOIN Content c ON b.banner_id = c.banner_id
	LEFT JOIN TagBanner tb ON b.banner_id = tb.banner_id
	WHERE tb.tag_id = ?
	GROUP BY b.banner_id, b.feature_id, c.title, c.text, c.url, c.version, c.updated_at, c.is_active
	LIMIT ? OFFSET ?;
	`

		rows, err := s.db.Query(query, req.Tag_id, req.Limit, req.Offset)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&res.Banner_id, &res.Tag_ids, &res.Feature_id, &res.Content.Title, &res.Content.Text, &res.Content.URL, &res.Is_active, &res.Created_at, &res.Updated_at); err != nil {
				panic(err)
			}
			ress = append(ress, res)
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	default:
		fmt.Println("unreachable sit")
	}
	return ress, nil
}

func (s *Storage) DeleteBanner(id int64) (banner.Response, error) {
	var res banner.Response
	const op = "storage.sqlite.DeleteBanner"

	stmt, err := s.db.Prepare("DELETE FROM Content WHERE banner_id = ?")
	if err != nil {
		return res, fmt.Errorf("%s:%w", op, err)
	}
	_, err = stmt.Exec(stmt, id)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return res, fmt.Errorf("%s:%w", op, storage.ErrURLExists)
		}
		return res, nil
	}
	return res, nil
}
