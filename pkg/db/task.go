package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

/*
	type Task struct {
		ID      string `json:"id,omitempty"`
		Date    string `json:"date,omitempty"`
		Title   string `json:"title"`
		Comment string `json:"comment,omitempty"`
		Repeat  string `json:"repeat,omitempty"`
	}
*/
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
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

func Tasks(db *sql.DB, limit int) ([]*Task, error) {
	query := `SELECT * FROM scheduler ORDER BY date LIMIT ` + strconv.Itoa(limit) + `;`
	tasks := []*Task{}
	fmt.Println("NULL tasks")
	rows, err := db.Query(query)
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return []*Task{}, err
		}
		tasks = append(tasks, &t)
	}
	fmt.Println("DB tasks", tasks)
	return tasks, nil
}

func TasksSearch(db *sql.DB, limit int, srchDate string, isDate bool) ([]*Task, error) {
	fmt.Println("START TasksSearch", srchDate, isDate)
	tasks := []*Task{}
	var query string
	if isDate {
		query = `SELECT * FROM scheduler WHERE date = :search ORDER BY date LIMIT :limit;`
	} else {
		query = `SELECT * FROM scheduler WHERE title LIKE concat('%',:search,'%') OR comment LIKE concat('%',:search,'%') ORDER BY date LIMIT :limit;`

	}
	rows, err := db.Query(query,
		sql.Named("search", srchDate),
		sql.Named("limit", limit))
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return []*Task{}, err
		}
		tasks = append(tasks, &t)
	}
	fmt.Println("DB tasks", tasks)
	return tasks, nil

}

// GetTask - запросить задачу из таблицы scheduler по id.
func GetTask(db *sql.DB, id int) (*Task, error) {
	query := `SELECT * FROM scheduler WHERE id = :id`
	var t Task
	err := db.QueryRow(query, sql.Named("id", id)).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return &Task{}, err
	}
	fmt.Println("DB tasks", t)
	return &t, nil
}

// UpdateTask - изменить задачу в таблице scheduler по id.
func UpdateTask(db *sql.DB, task *Task) error {
	// параметры пропущены, не забудьте указать WHERE
	query := `UPDATE scheduler 
				SET date=:date, title=:title, comment=:comment, repeat=:repeat 
				WHERE id=:id`
	res, err := db.Exec(query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}
