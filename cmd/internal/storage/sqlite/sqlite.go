package sqlite

import (
	"database/sql"
	"fmt"
	"time"
)

type Storage struct {
	db                           *sql.DB
	CreateBannerContentStmt      *sql.Stmt
	CreateBannerFeatureStmt      *sql.Stmt
	CreateTagBannerStmt          *sql.Stmt
	UpdateBannerTagStmt          *sql.Stmt
	UpdatBannerFeatureStmt       *sql.Stmt
	UpdateBannerContentTitleStmt *sql.Stmt
	UpdateBannerContentTextStmt  *sql.Stmt
	UpdateBannerContentURLStmt   *sql.Stmt
	UpdateBannerIsActiveStmt     *sql.Stmt
	GetBannerStmt                *sql.Stmt
	GetBannersByFeatureStmt      *sql.Stmt
	GetBannersByTagStmt          *sql.Stmt
	GetBannersStmt               *sql.Stmt
}

func New(storagePath string, maxOpenConns int, maxIdleConns int, connMaxLifeTime time.Duration) (*Storage, error) {
	const op = "storage.sqlite.New"

	creators := make([]string, 0, 8)

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifeTime)

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

	CreateBannerContentStmt, err := db.Prepare("INSERT INTO Content(title, text, url) VALUES(?, ?, ?)")
	if err != nil {
		return nil, err
	}
	fmt.Println("success CreateBannerContentStmt")

	CreateBannerFeatureStmt, err := db.Prepare("INSERT INTO Banner(banner_id, feature_id) VALUES(?, ?)")
	if err != nil {
		return nil, err
	}
	
	fmt.Println("success CreateBannerFeatureStmt")

	CreateTagBannerStmt, err := db.Prepare("INSERT INTO TagBanner(tag_id, banner_id) VALUES(?, ?)")
	if err != nil {
		return nil, err
	}

	fmt.Println("success CreateTagBannerStmt")
	//написать
	UpdateBannerTagStmt, err := db.Prepare("INSERT INTO TagBanner(tag_id, banner_id) VALUES(?, ?)")
	if err != nil {
		return nil, err
	}

	fmt.Println("success UpdateBannerTagStmt")

	UpdatBannerFeatureStmt, err := db.Prepare("UPDATE Banner SET feature_id = ? WHERE banner_id = ?")
	if err != nil {
		return nil, err
	}

	fmt.Println("success UpdatBannerFeatureStmt")

	UpdateBannerContentTitleStmt, err := db.Prepare("UPDATE Content SET title = ? WHERE banner_id = ?")
	if err != nil {
		return nil, err
	}

	fmt.Println("success UpdateBannerContentTitleStmt")

	UpdateBannerContentTextStmt, err := db.Prepare("UPDATE Content SET text = ? WHERE banner_id = ?")
	if err != nil {
		return nil, err
	}

	fmt.Println("success UpdateBannerContentTextStmt")

	UpdateBannerContentURLStmt, err := db.Prepare("UPDATE Content SET url = ? WHERE banner_id = ?")
	if err != nil {
		return nil, err
	}

	fmt.Println("success UpdateBannerContentURLStmt")

	UpdateBannerIsActiveStmt, err := db.Prepare("UPDATE Content SET is_active = ? WHERE banner_id = ?")
	if err != nil {
		return nil, err
	}

	fmt.Println("success UpdateBannerIsActiveStmt")

	GetBannerStmt, err := db.Prepare(`SELECT c.title, c.text, c.url
        FROM Content c
        INNER JOIN Banner b ON c.banner_id = b.banner_id
        INNER JOIN TagBanner tb ON b.banner_id = tb.banner_id
        WHERE tb.tag_id = ? AND b.feature_id = ?;`)
	if err != nil {
		return nil, err
	}

	fmt.Println("success GetBannerStmt")

	GetBannersByFeatureStmt, err := db.Prepare(`SELECT 
		b.banner_id, 
		GROUP_CONCAT(tb.tag_id, ',') AS tag_ids, 
		b.feature_id, 
		c.title, 
		c.text, 
		c.url, 
		c.is_active, 
		c.version AS created_at, 
		c.updated_at
	FROM 
		Banner b
	INNER JOIN 
		Content c ON b.banner_id = c.banner_id
	INNER JOIN 
		TagBanner tb ON b.banner_id = tb.banner_id
	WHERE 
		b.feature_id = ?
	GROUP BY 
		b.banner_id, 
		b.feature_id, 
		c.title, 
		c.text, 
		c.url, 
		c.is_active, 
		c.version, 
		c.updated_at
		LIMIT ? OFFSET ?;
	`)
	if err != nil {
		return nil, err
	}

	fmt.Println("success GetBannersByFeatureStmt")

	GetBannersByTagStmt, err := db.Prepare(`SELECT 
		b.banner_id, 
		GROUP_CONCAT(tb.tag_id, ',') AS tag_ids, 
		b.feature_id, 
		c.title, 
		c.text, 
		c.url, 
		c.is_active, 
		c.version AS created_at, 
		c.updated_at
	FROM 
		Banner b
	INNER JOIN 
		Content c ON b.banner_id = c.banner_id
	INNER JOIN 
		TagBanner tb ON b.banner_id = tb.banner_id
	WHERE 
		b.feature_id = ?
	GROUP BY 
		b.banner_id, 
		b.feature_id, 
		c.title, 
		c.text, 
		c.url, 
		c.is_active, 
		c.version, 
		c.updated_at
		LIMIT ? OFFSET ?;
	`)
	if err != nil {
		return nil, err
	}

	fmt.Println("success GetBannersByTagStmt")

	GetBannersStmt, err := db.Prepare(`SELECT 
    b.banner_id, 
    GROUP_CONCAT(tb.tag_id, ',') AS tag_ids, 
    b.feature_id, 
    c.title, 
    c.text, 
    c.url, 
    c.is_active, 
    c.version AS created_at, 
    c.updated_at
FROM 
    Banner b
INNER JOIN 
    Content c ON b.banner_id = c.banner_id
INNER JOIN 
    TagBanner tb ON b.banner_id = tb.banner_id
WHERE 
    tb.tag_id = ? AND b.feature_id = ?
GROUP BY 
    b.banner_id, 
    b.feature_id, 
    c.title, 
    c.text, 
    c.url, 
    c.is_active, 
    c.version, 
    c.updated_at
	LIMIT ? OFFSET ?;
	`)
	if err != nil {
		return nil, err
	}

	fmt.Println("success GetBannersStmt")

	return &Storage{
		db:                           db,
		CreateBannerContentStmt:      CreateBannerContentStmt,
		CreateBannerFeatureStmt:      CreateBannerFeatureStmt,
		CreateTagBannerStmt:          CreateTagBannerStmt,
		UpdateBannerTagStmt:          UpdateBannerTagStmt,
		UpdatBannerFeatureStmt:       UpdatBannerFeatureStmt,
		UpdateBannerContentTitleStmt: UpdateBannerContentTitleStmt,
		UpdateBannerContentTextStmt:  UpdateBannerContentTextStmt,
		UpdateBannerContentURLStmt:   UpdateBannerContentURLStmt,
		UpdateBannerIsActiveStmt:     UpdateBannerIsActiveStmt,
		GetBannerStmt:                GetBannerStmt,
		GetBannersByFeatureStmt:      GetBannersByFeatureStmt,
		GetBannersByTagStmt:          GetBannersByTagStmt,
		GetBannersStmt:               GetBannersStmt,
	}, nil
}

func (s *Storage) Close() error {

	s.CreateBannerContentStmt.Close()
	s.CreateBannerFeatureStmt.Close()
	s.CreateTagBannerStmt.Close()
	s.UpdateBannerTagStmt.Close()
	s.UpdateBannerContentTitleStmt.Close()
	s.UpdateBannerContentTextStmt.Close()
	s.UpdateBannerContentURLStmt.Close()
	s.UpdateBannerIsActiveStmt.Close()
	s.GetBannerStmt.Close()
	s.GetBannersByFeatureStmt.Close()
	s.GetBannersByTagStmt.Close()
	s.GetBannersStmt.Close()
	s.db.Close()

	return s.db.Close()
}
