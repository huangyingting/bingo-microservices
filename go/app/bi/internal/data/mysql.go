package data

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlBIStore struct {
	db  *sql.DB
	ctx context.Context
}

func NewMysqlBIStore(dsn string) (*MysqlBIStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &MysqlBIStore{db: db, ctx: context.Background()}, nil
}

func (s *MysqlBIStore) Open() error {
	const INITIALIZE_STATEMENT = `
	CREATE TABLE IF NOT EXISTS clicks (
		id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		alias VARCHAR(64) NOT NULL,
		ip VARCHAR(15) NOT NULL,
		ua VARCHAR(1024) NOT NULL,
		referer VARCHAR(2048) NOT NULL,
		country VARCHAR(8) DEFAULT '',
		city VARCHAR(128) DEFAULT '',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX (alias)
	);
	`
	_, err := s.db.ExecContext(s.ctx, INITIALIZE_STATEMENT)
	return err
}

func (s *MysqlBIStore) Close() error {
	return s.db.Close()
}

func (s *MysqlBIStore) Create(click *Click) error {
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

func (s *MysqlBIStore) Clicks(alias string) (uint64, error) {
	const COUNT_STATEMENT = `
	SELECT COUNT(*) FROM clicks
	WHERE alias = ?
	`
	row := s.db.QueryRowContext(s.ctx, COUNT_STATEMENT, alias)
	var i uint64
	err := row.Scan(&i)
	return i, err
}
