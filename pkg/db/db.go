package db
import (
	"modernc.org/sqlite"
)
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
id
)
`
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
func createSchemaProcess(db *sql.DB) error{
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