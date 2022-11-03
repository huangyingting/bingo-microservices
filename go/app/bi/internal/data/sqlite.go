package data

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteBIStore struct {
	db  *sql.DB
	ctx context.Context
}

func NewSqliteBIStore(dsn string) (*SqliteBIStore, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return &SqliteBIStore{db: db, ctx: context.Background()}, nil
}

func (s *SqliteBIStore) Open() error {
	const INITIALIZE_STATEMENT = `
	BEGIN;
	CREATE TABLE IF NOT EXISTS clicks (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		alias VARCHAR(64) NOT NULL,
		ip VARCHAR(15) NOT NULL,
		ua VARCHAR(1024) NOT NULL,
		referer VARCHAR(2048) NOT NULL,
		country VARCHAR(8) DEFAULT '',
		city VARCHAR(128) DEFAULT '',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS IDX_ALIAS ON clicks(alias);
	COMMIT;
	`
	_, err := s.db.ExecContext(s.ctx, INITIALIZE_STATEMENT)
	return err
}

func (s *SqliteBIStore) Close() error {
	return s.db.Close()
}

func (s *SqliteBIStore) Create(click *Click) error {
	const CREATE_STATEMENT = `
	INSERT INTO clicks (
		alias, ip, ua, referer, country, city, created_at
	) VALUES (
		?, ?, ?, ?, ?, ?, ?
	)
	`
	_, err := s.db.ExecContext(
		s.ctx,
		CREATE_STATEMENT,
		click.Alias,
		click.IP,
		click.UA,
		click.Referer,
		click.Country,
		click.City,
		click.CreatedAt,
	)
	return err
}

func (s *SqliteBIStore) Clicks(alias string) (uint64, error) {
	const COUNT_STATEMENT = `
	SELECT COUNT(*) FROM clicks
	WHERE alias = ?
	`
	row := s.db.QueryRowContext(s.ctx, COUNT_STATEMENT, alias)
	var i uint64
	err := row.Scan(&i)
	return i, err
}
