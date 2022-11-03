package data

import (
	"context"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

type SQLServerBIStore struct {
	db  *sql.DB
	ctx context.Context
}

func NewSQLServerBIStore(dsn string) (*SQLServerBIStore, error) {
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}
	return &SQLServerBIStore{db: db, ctx: context.Background()}, nil
}

func (s *SQLServerBIStore) Open() error {
	const INITIALIZE_STATEMENT = `
	BEGIN TRANSACTION;
	IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='clicks' and xtype='U')
		CREATE TABLE clicks (
		id BIGINT NOT NULL IDENTITY PRIMARY KEY,
		alias VARCHAR(64) NOT NULL,
		ip VARCHAR(15) NOT NULL,
		ua VARCHAR(1024) NOT NULL,
		country VARCHAR(8) DEFAULT '',
		city VARCHAR(128) DEFAULT '',
		referer VARCHAR(2048) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	IF NOT EXISTS(SELECT * FROM sys.indexes WHERE Name = 'IDX_ALIAS')
		CREATE NONCLUSTERED INDEX IDX_ALIAS ON clicks(alias);
	COMMIT;
	`
	_, err := s.db.ExecContext(s.ctx, INITIALIZE_STATEMENT)
	return err
}

func (s *SQLServerBIStore) Close() error {
	return s.db.Close()
}

func (s *SQLServerBIStore) Create(click *Click) error {
	const CREATE_STATEMENT = `
	INSERT INTO clicks (
		alias, ip, ua, referer, country, city, created_at
	) VALUES (
		@p1, @p2, @p3, @p4, @p5, @p6, @p7
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

func (s *SQLServerBIStore) Clicks(alias string) (uint64, error) {
	const COUNT_STATEMENT = `
	SELECT COUNT(*) FROM clicks
	WHERE alias = @p1
	`
	row := s.db.QueryRowContext(s.ctx, COUNT_STATEMENT, alias)
	var i uint64
	err := row.Scan(&i)
	return i, err
}
