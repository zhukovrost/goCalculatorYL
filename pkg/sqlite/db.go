package sqlite

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	const (
		usersTable = `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
		login VARCHAR(255) UNIQUE NOT NULL CHECK ( length(login) > 3 ),
		password VARCHAR(255) NOT NULL
	);`

		expressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions(
		id INTEGER NOT NULL PRIMARY KEY, 
		expression TEXT NOT NULL,
		creator INTEGER NOT NULL,
		result INTEGER NOT NULL DEFAULT 0,
		status varchar(16) NOT NULL DEFAULT 'pending',
	
		FOREIGN KEY (creator)  REFERENCES expressions (id)
	);`
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, expressionsTable); err != nil {
		return err
	}

	return nil
}
