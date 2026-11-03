package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const tableSchema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(32) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
);`

var DB *sql.DB

func Init(dbFile string) error {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	if _, err = db.Exec(tableSchema); err != nil {
		return err
	}

	_, _ = db.Exec("CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);")

	DB = db
	return nil
}
