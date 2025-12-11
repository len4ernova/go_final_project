// db package organizes work with the database.
package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
	CREATE TABLE IF NOT EXISTS scheduler (
		id integer PRIMARY KEY AUTOINCREMENT,
		date char(8) NOT NULL DEFAULT "",
		title varchar(256) NOT NULL DEFAULT "",
		comment text NOT NULL DEFAULT "",
		repeat varchar(128) NOT NULL DEFAULT ""
	);`
const index = `CREATE INDEX schedule_date ON scheduler (date);`

var PlanDB *sql.DB

// Init - инициализация БД.
func Init(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		return db, err
	}

	if !checkExist(dbName) {
		return db, create(db, schema+index)
	}
	return db, nil
}

// checkExist - проверка существование БД.
func checkExist(dbFile string) bool {
	_, err := os.Stat(dbFile)
	if err != nil {
		return false
	}

	return true
}

// create - выполнить запрос к БД.
func create(db *sql.DB, schm string) error {
	_, err := db.Exec(schm)
	if err != nil {
		return err
	}
	return nil
}
