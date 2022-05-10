package data

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresBIStore struct {
	db  *sql.DB
	ctx context.Context
}

func NewPostgresBIStore(dsn string) (*PostgresBIStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresBIStore{db: db, ctx: context.Background()}, nil
}

func (s *PostgresBIStore) Open() error {
	const INITIALIZE_STATEMENT = `
	BEGIN;
	CREATE TABLE IF NOT EXISTS clicks (
		id BIGSERIAL PRIMARY KEY,
		alias VARCHAR(64) NOT NULL,
		ip VARCHAR(15) NOT NULL,
		ua VARCHAR(1024) NOT NULL,
		referer VARCHAR(2048) NOT NULL,
		country VARCHAR(8),
		city VARCHAR(128),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS IDX_ALIAS ON clicks (alias);
	COMMIT;	
	`
	_, err := s.db.ExecContext(s.ctx, INITIALIZE_STATEMENT)
	return err
}

func (s *PostgresBIStore) Close() error {
	return s.db.Close()
}

func (s *PostgresBIStore) Create(click *Click) error {
	const CREATE_STATEMENT = `
	INSERT INTO clicks (
		alias, ip, ua, referer, country, city, created_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
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

func (s *PostgresBIStore) Clicks(alias string) (uint64, error) {
	const COUNT_STATEMENT = `
	SELECT COUNT(*) FROM clicks
	WHERE alias = $1
	`
	row := s.db.QueryRowContext(s.ctx, COUNT_STATEMENT, alias)
	var i uint64
	err := row.Scan(&i)
	return i, err
}
