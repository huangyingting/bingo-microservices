package data

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteBSStore struct {
	db  *sql.DB
	ctx context.Context
}

func NewSqliteBSStore(dsn string) (*SqliteBSStore, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return &SqliteBSStore{db: db, ctx: context.Background()}, nil
}

func (q *SqliteBSStore) Open() error {
	const INITIALIZE_STATEMENT = `
	BEGIN;
	CREATE TABLE IF NOT EXISTS short_url (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		alias VARCHAR(64) NOT NULL UNIQUE,
		url VARCHAR(2048) NOT NULL,
		oid CHAR(36) NOT NULL,
		title VARCHAR(256) DEFAULT '',
		tags VARCHAR(8000) DEFAULT '',
		flags UNSIGNED INT DEFAULT 0,
		utm_source VARCHAR(128) DEFAULT '',
		utm_medium VARCHAR(128) DEFAULT '',
		utm_campaign VARCHAR(128) DEFAULT '',
		utm_term VARCHAR(128) DEFAULT '',
		utm_content VARCHAR(128) DEFAULT '',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS IDX_ALIAS_OID ON short_url (alias,oid);
	COMMIT;
	`
	_, err := q.db.ExecContext(q.ctx, INITIALIZE_STATEMENT)
	return err
}

func (q *SqliteBSStore) Close() error {
	return q.db.Close()
}

// ShortUrl
func (q *SqliteBSStore) CreateShortUrl(
	alias string,
	customized bool,
	url string,
	oid string,
) error {
	const CREATE_SHORT_URL = `
	INSERT INTO short_url (
		alias, url, oid, flags
	) VALUES (
		?, ?, ?, ?
	)
	`
	flags := Bits(0).Set(FLAG_CUSTOMIZED, customized)
	_, err := q.db.ExecContext(q.ctx, CREATE_SHORT_URL, alias, url, oid, flags)
	return err
}

func (q *SqliteBSStore) DeleteShortUrl(alias string, oid string) error {
	const DELETE_SHORT_URL = `
	DELETE FROM short_url
	WHERE alias = ? AND oid = ?
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

func (q *SqliteBSStore) GetShortUrl(alias string) (*ShortUrl, error) {
	const GET_SHORT_URL = `
	SELECT alias, url, oid, title, tags, flags, utm_source, utm_medium, utm_campaign, utm_term, utm_content, created_at FROM short_url
	WHERE alias = ? LIMIT 1
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

func (q *SqliteBSStore) GetShortUrlByOid(alias string, oid string) (*ShortUrl, error) {
	const GET_SHORT_URL = `
		SELECT alias, url, oid, title, tags, flags, utm_source, utm_medium, utm_campaign, utm_term, utm_content, created_at FROM short_url
		WHERE alias = ? AND oid = ? LIMIT 1
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

func (q *SqliteBSStore) ListShortUrl(oid string, start int64, count int64) ([]*ShortUrl, error) {
	const LIST_SHORT_URL = `
	SELECT alias, url, oid, title, tags, flags, utm_source, utm_medium, utm_campaign, utm_term, utm_content, created_at FROM short_url
	WHERE oid = ? ORDER BY created_at DESC
	LIMIT ? OFFSET ?
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

func (q *SqliteBSStore) UpdateShortUrl(
	alias string,
	oid string,
	updateShortUrl UpdateShortUrl,
) error {
	const UPDATE_SHORT_URL = `
	UPDATE short_url SET url = ?, title = ?, tags = ?, flags = (flags & 1) + ?, utm_source = ?, utm_medium = ?, utm_campaign = ?, utm_term = ?, utm_content = ? 
	where alias = ? AND oid = ?
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
