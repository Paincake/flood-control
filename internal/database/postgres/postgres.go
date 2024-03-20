package postgres

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

func New(dbname, username, password, host, port string) (*Database, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?&sslmode=disable",
		username,
		password,
		host,
		port,
		dbname)
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &Database{db: db}, nil
}
func (d *Database) MakeCallAndGetRecent(ctx context.Context, userId int64, n int) (int, error) {
	resultChan := make(chan int)
	errChan := make(chan error)
	go func() {
		tx, err := d.db.Begin()
		if err != nil {
			errChan <- err
			return
		}
		_, err = tx.Exec("INSERT INTO calls (user_id, call_timestamp) VALUES ($1, current_timestamp)", userId)
		res, err := tx.Query("SELECT COUNT(*) FROM calls WHERE call_timestamp BETWEEN current_timestamp - INTERVAL '$ seconds' AND current_timestamp", n)
		if err != nil {
			errChan <- err
			return
		}
		err = tx.Commit()
		if err != nil {
			errChan <- err
			return
		}
		res.Next()
		var count int
		err = res.Scan(&count)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- count
	}()
	for {
		select {
		case <-ctx.Done():
			return -1, ctx.Err()
		case err := <-errChan:
			return -1, err
		case val := <-resultChan:
			return val, nil
		}
	}
}
