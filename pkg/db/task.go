package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// AddTask - добавление задачи в таблицу scheduler и возврат идентификатора добавленной записи.
func AddTask(db *sql.DB, task *Task) (int64, error) {

	// определите запрос
	fmt.Println("task.Date, task.Title, task.Comment, task.Repeat:", task.Date, task.Title, task.Comment, task.Repeat)

	query := `INSERT INTO scheduler
	(date, title, comment, repeat)
	VALUES (:date, :title, :comment, :repeat);
	`
	fmt.Println("2")
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()

	return id, err
}

func Tasks(limit int) ([]*Task, error) {

	// избежать {"tasks":null} в ответе JSON, следите за результирующим слайсом.
	//  Если он равен nil, нужно создавать пустой слайс.
	//  В этом случае ответ будет {"tasks":[]}
}
