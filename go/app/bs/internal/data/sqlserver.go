package data

import (
	"context"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

type SQLServerBSStore struct {
	db  *sql.DB
	ctx context.Context
}

func NewSQLServerBSStore(dsn string) (*SQLServerBSStore, error) {
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}
	return &SQLServerBSStore{db: db, ctx: context.Background()}, nil
}

func (q *SQLServerBSStore) Open() error {
	const INITIALIZE_STATEMENT = `
	BEGIN TRANSACTION; 
	IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='short_url' and xtype='U')
		CREATE TABLE short_url (
			id BIGINT NOT NULL IDENTITY PRIMARY KEY,
			alias VARCHAR(64) NOT NULL,
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
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT AK_ALIAS UNIQUE(alias)
		);
	IF NOT EXISTS(SELECT * FROM sys.indexes WHERE Name = 'IDX_ALIAS_OID')
		CREATE NONCLUSTERED INDEX IDX_ALIAS_OID ON short_url(alias, oid);
	COMMIT;
	`
	_, err := q.db.ExecContext(q.ctx, INITIALIZE_STATEMENT)
	return err
}

func (q *SQLServerBSStore) Close() error {
	return q.db.Close()
}

func (q *SQLServerBSStore) CreateShortUrl(
	alias string,
	customized bool,
	url string,
	oid string,
) error {
	const CREATE_SHORT_URL = `
	INSERT INTO short_url (
		alias, url, oid, flags
	) VALUES (
		@p1, @p2, @p3, @p4
	)
	`
	flags := Bits(0).Set(FLAG_CUSTOMIZED, customized)
	_, err := q.db.ExecContext(q.ctx, CREATE_SHORT_URL, alias, url, oid, flags)
	return err
}

func (q *SQLServerBSStore) DeleteShortUrl(alias string, oid string) error {
	const DELETE_SHORT_URL = `
	DELETE FROM short_url
	WHERE alias = @p1 AND oid = @p2
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

func (q *SQLServerBSStore) GetShortUrl(alias string) (*ShortUrl, error) {
	const GET_SHORT_URL = `
	SELECT TOP 1 alias, url, oid, title, tags, flags, utm_source, utm_medium, utm_campaign, utm_term, utm_content, created_at FROM short_url
	WHERE alias = @p1
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

func (q *SQLServerBSStore) GetShortUrlByOid(alias string, oid string) (*ShortUrl, error) {
	const GET_SHORT_URL = `
	SELECT TOP 1 alias, url, oid, title, tags, flags, utm_source, utm_medium, utm_campaign, utm_term, utm_content, created_at FROM short_url
	WHERE alias = @p1 AND oid = @p2
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

func (q *SQLServerBSStore) ListShortUrl(
	oid string,
	start int64,
	count int64,
) ([]*ShortUrl, error) {
	const LIST_SHORT_URL = `
	SELECT alias, url, oid, title, tags, flags, utm_source, utm_medium, utm_campaign, utm_term, utm_content, created_at FROM short_url
	WHERE oid = @p1 ORDER BY created_at DESC
	OFFSET @p2 ROWS
	FETCH NEXT @p3 ROWS ONLY
	`
	rows, err := q.db.QueryContext(q.ctx, LIST_SHORT_URL, oid, start, count)
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

func (q *SQLServerBSStore) UpdateShortUrl(
	alias string,
	oid string,
	updateShortUrl UpdateShortUrl,
) error {
	const UPDATE_SHORT_URL = `
	UPDATE short_url SET url = @p1, title = @p2, tags = @p3, flags = (flags & 1) + @p4, utm_source = @p5, utm_medium = @p6, utm_campaign = @p7, utm_term = @p8, utm_content = @p9
	where alias = @p10 AND oid = @p11
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
