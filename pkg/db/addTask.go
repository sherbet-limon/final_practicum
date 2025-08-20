package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)
//вставляет в бд новую запись и возвращает id
func AddTask(task *Task) (int64, error) {
	var id int64
	db, err := sql.Open("sqlite", "scheduler.db")
    if err != nil {
        log.Println("Ошибка подключения к БД")
        return 0, err
    }
    defer db.Close()

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat);",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
		if err == nil {
			id, err = res.LastInsertId()
		}
		return id, err
	} 
