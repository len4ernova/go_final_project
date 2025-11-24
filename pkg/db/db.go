package db

import (
	"database/sql"
	"os"
)

const schema = `
	CREATE TABLE IF NOT EXISTS scheduler (
		id integer PRIMARY KEY AUTOINCREMENT,
		date char(8) NOT NULL DEFAULT "",
		title varchar(256) NOT NULL DEFAULT "",
		comment text NOT NULL DEFAULT "",
		repeat varchar(128) NOT NULL DEFAULT ""
	);`
const ind = `CREATE INDEX schedule_date ON scheduler (date);`

// Init - инициализация БД.gibt
func Init(dbFile string) {
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	if bool {
		// TODO create
		return
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

}
func createSchemaProcess(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS processes (
		id_proc integer PRIMARY KEY AUTOINCREMENT,
		name_proc text NOT NULL DEFAULT ""  UNIQUE
	);`)
	if err != nil {
		return err
	}
	return nil
}
