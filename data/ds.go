package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Source interface {
	Query(ctx context.Context, query string) (Set, error)
	Close() error
}

type Set interface {
	Next() bool
	Err() error

	Scan(dest ...interface{}) error

	Columns() ([]string, error)

	Close() error
}

type ReusableSet interface {
	Set

	reusable()
	Count() int
	Reset()
}

func OpenSQL(driver, conn string) (Source, error) {
	db, err := sql.Open(driver, conn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := db.QueryContext(ctx, "select 1")
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	if !r.Next() {
		_ = db.Close()
		return nil, fmt.Errorf("cannot confirm connection")
	}
	return (*sqlDs)(db), nil
}

type sqlDs sql.DB

func (s *sqlDs) Query(ctx context.Context, query string) (Set, error) {
	db := (*sql.DB)(s)
	return db.QueryContext(ctx, query)
}

func (s *sqlDs) Close() error {
	db := (*sql.DB)(s)
	return db.Close()
}
