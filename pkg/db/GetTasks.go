package db

import (
	"database/sql"
	"errors"
	"log"
	_ "modernc.org/sqlite"
)

const query = "SELECT ID, Date, Title, Comment, Repeat FROM scheduler ORDER BY Date LIMIT $1"
const queryId = "SELECT ID, Date, Title, Comment, Repeat FROM scheduler WHERE ID=?"
const queryUpdate = "UPDATE scheduler SET Date=?, Title=?, Comment=?, Repeat=? WHERE ID=?"
const queryUpdateDate = "UPDATE scheduler SET Date=? WHERE ID=?"
const queryDelete = "DELETE FROM scheduler WHERE ID=?"
type TasksResp struct {
	Tasks []*Task `json:"tasks"`
}

func GetTasksByDate(limit int) ([]*Task, error) {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Println("не удалось открыть БД")
		return nil, errors.New("не удалось открыть БД")
	}
	defer db.Close()

	rows, err := db.Query(query, limit)
	if err != nil {
		log.Println("не удалось получить данные из БД")
		return nil, errors.New("не удалось получить данные из БД")
	}
	defer rows.Close()
	
	tasks := []*Task{}
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Println("не удалось записать данные в структуру")
			return nil, errors.New("не удалось записать данные в структуру")
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func GetTasksByID(id string)(*Task, error){
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Println("не удалось открыть БД")
		return nil, errors.New("не удалось открыть БД")
	}
	defer db.Close()

	task := &Task{}
	row := db.QueryRow(queryId, id) 
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		log.Println("не удалось получить данные из БД")
		return nil, errors.New("не удалось получить данные из БД")
	}
	return task, nil
}

func UpdateTask(task *Task) error {
    db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Println("не удалось открыть БД")
		return errors.New("не удалось открыть БД")
	}
	defer db.Close()

	res, err := db.Exec(queryUpdate, 
        task.Date,
        task.Title,
        task.Comment,
        task.Repeat,
        task.ID)
    if err != nil {
        return errors.New("не удалось обновить данные в БД")
    }
    // проверяем сколько строк было изменено
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return errors.New("некорректный id")
    }
    return nil
}

func UpdateDate(task *Task, id string) error {
    db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Println("не удалось открыть БД")
		return errors.New("не удалось открыть БД")
	}
	defer db.Close()

	res, err := db.Exec(queryUpdateDate, task.Date, id)
    if err != nil {
        return errors.New("не удалось обновить данные в БД")
    }
    // проверяем сколько строк было изменено
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return errors.New("некорректный id")
    }
    return nil
}

func DeleteTask(id string) error {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Println("не удалось открыть БД")
		return errors.New("не удалось открыть БД")
	}
	defer db.Close()

	res, err := db.Exec(queryDelete, id)
	if err != nil {
		log.Println("не удалось удалить данные из БД")
		return errors.New("не удалось удалить данные из БД")
	}
	count, err := res.RowsAffected()
    if err != nil {
        log.Println("не удалась проверка измененных строк", err)
        return errors.New("не удалась проверка измененных строк")
    }
    if count == 0 {
        return errors.New("ничего не удалено")
    }
	return nil
}
