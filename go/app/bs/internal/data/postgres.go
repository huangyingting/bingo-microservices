package data

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresBSStore struct {
	db  *sql.DB
	ctx context.Context
}

func NewPostgresBSStore(dsn string) (*PostgresBSStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresBSStore{db: db, ctx: context.Background()}, nil
}

func (q *PostgresBSStore) Open() error {
	const INITIALIZE_STATEMENT = `
	BEGIN;
	CREATE TABLE IF NOT EXISTS short_url (
		id BIGSERIAL PRIMARY KEY,
		alias VARCHAR(64) NOT NULL UNIQUE,
		url VARCHAR(2048) NOT NULL,
		oid CHAR(36) NOT NULL,
		title VARCHAR(256) DEFAULT '',
		tags VARCHAR(8000) DEFAULT '',
		flags INT DEFAULT 0,
		utm_source VARCHAR(128) DEFAULT '',
		utm_medium VARCHAR(128) DEFAULT '',
		utm_campaign VARCHAR(128) DEFAULT '',
		utm_term VARCHAR(128) DEFAULT '',
		utm_content VARCHAR(128) DEFAULT '',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS IDX_ALIAS_OID ON short_url (alias, oid);
	COMMIT;
	`
	_, err := q.db.ExecContext(q.ctx, INITIALIZE_STATEMENT)
	return err
}

func (q *PostgresBSStore) Close() error {
	return q.db.Close()
}

func (q *PostgresBSStore) CreateShortUrl(
	alias string,
	customized bool,
	url string,
	oid string,
) error {
	const CREATE_SHORT_URL = `
	INSERT INTO short_url (
		alias, url, oid, flags
	) VALUES (
		$1, $2, $3, $4
	)
	`
	flags := Bits(0).Set(FLAG_CUSTOMIZED, customized)
	_, err := q.db.ExecContext(q.ctx, CREATE_SHORT_URL, alias, url, oid, flags)
	return err
}

func (q *PostgresBSStore) DeleteShortUrl(alias string, oid string) error {
	const DELETE_SHORT_URL = `
	DELETE FROM short_url
	WHERE alias = $1 AND oid = $2
	`
	result, err := q.db.ExecContext(q.ctx, DELETE_SHORT_URL, alias, oid)
	if err != nil {
		return err
	}
	if row, _ := result.RowsAffected(); row == 0 {
		err = ErrNoRowsDeleted
	}
	return err
}

func (q *PostgresBSStore) GetShortUrl(alias string) (*ShortUrl, error) {
	const GET_SHORT_URL = `
	SELECT alias, url, oid, title, tags, flags, utm_source, utm_medium, utm_campaign, utm_term, utm_content, created_at FROM short_url
	WHERE alias = $1 LIMIT 1
	`
	row := q.db.QueryRowContext(q.ctx, GET_SHORT_URL, alias)
	var i ShortUrl
	err := row.Scan(
		&i.Alias,
		&i.Url,
		&i.Oid,
		&i.Title,
		&i.Tags,
		&i.Flags,
		&i.UtmSource,
		&i.UtmMedium,
		&i.UtmCampaign,
		&i.UtmTerm,
		&i.UtmContent,
		&i.CreatedAt,
	)
	return &i, err
}

func (q *PostgresBSStore) GetShortUrlByOid(alias string, oid string) (*ShortUrl, error) {
	const GET_SHORT_URL = `
	SELECT alias, url, oid, title, tags, flags, utm_source, utm_medium, utm_campaign, utm_term, utm_content, created_at FROM short_url
	WHERE alias = $1 AND oid = $2 LIMIT 1
	`
	row := q.db.QueryRowContext(q.ctx, GET_SHORT_URL, alias, oid)
	var i ShortUrl
	err := row.Scan(
		&i.Alias,
		&i.Url,
		&i.Oid,
		&i.Title,
		&i.Tags,
		&i.Flags,
		&i.UtmSource,
		&i.UtmMedium,
		&i.UtmCampaign,
		&i.UtmTerm,
		&i.UtmContent,
		&i.CreatedAt,
	)
	return &i, err
}

func (q *PostgresBSStore) ListShortUrl(
	oid string,
	start int64,
	count int64,
) ([]*ShortUrl, error) {
	const LIST_SHORT_URL = `
	SELECT alias, url, oid, title, tags, flags, utm_source, utm_medium, utm_campaign, utm_term, utm_content, created_at FROM short_url
	WHERE oid = $1 ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := q.db.QueryContext(q.ctx, LIST_SHORT_URL, oid, count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ShortUrl
	for rows.Next() {
		var i ShortUrl
		if err := rows.Scan(
			&i.Alias,
			&i.Url,
			&i.Oid,
			&i.Title,
			&i.Tags,
			&i.Flags,
			&i.UtmSource,
			&i.UtmMedium,
			&i.UtmCampaign,
			&i.UtmTerm,
			&i.UtmContent,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *PostgresBSStore) UpdateShortUrl(
	alias string,
	oid string,
	updateShortUrl UpdateShortUrl,
) error {
	const UPDATE_SHORT_URL = `
	UPDATE short_url SET url = $1, title = $2, tags = $3, flags = (flags & 1) + $4, utm_source = $5, utm_medium = $6, utm_campaign = $7, utm_term = $8, utm_content = $9 
	where alias = $10 AND oid = $11
	`
	result, err := q.db.ExecContext(
		q.ctx,
		UPDATE_SHORT_URL,
		updateShortUrl.Url,
		updateShortUrl.Title,
		updateShortUrl.Tags,
		updateShortUrl.Flags,
		updateShortUrl.UtmSource,
		updateShortUrl.UtmMedium,
		updateShortUrl.UtmCampaign,
		updateShortUrl.UtmTerm,
		updateShortUrl.UtmContent,
		alias,
		oid,
	)
	if err != nil {
		return err
	}
	if row, _ := result.RowsAffected(); row == 0 {
		err = ErrNoRowsUpdated
	}
	return err
}
