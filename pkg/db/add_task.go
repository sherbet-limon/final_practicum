package db

import (
	"database/sql"
	// "go1f/pkg/db"
	"errors"

	_ "modernc.org/sqlite"
)

// вставляет в бд новую запись и возвращает id
func AddTask(task *Task) (int64, error) {
	var id int64
	if DB == nil {
		return 0, errors.New("БД не инициализирована")
	}
	res, err := DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat);",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}
