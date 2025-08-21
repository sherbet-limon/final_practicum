package api

import (
	"bytes"
	"encoding/json"
	"go1f/pkg/db"
	"log"
	"net/http"
	"strconv"
	"time"
)

//post - добавляет запись в бд после проверки заголовка, даты и repeat
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task *db.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeError(w, "ошибка чтения", http.StatusBadRequest)
		log.Println("Не читается тело запроса")
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeError(w, "ошибка десериализации", http.StatusBadRequest)
		log.Println("JSON не десериализовался")
		return
	}
	if task.ID != "" {
		writeError(w, "в запросе не должно быть ID", http.StatusBadRequest)
		log.Println("в запросе не должно быть ID")
		return
	}
	if task.Title == nil || *task.Title == "" {
		writeError(w, `{"error":"заголовок обязателен и не должен быть пустым"}`, http.StatusBadRequest)
		log.Println("Заголовка нет или пустой")
		return
	}

	// dateTaskNow - текущее время, если task.Date пустое или отсутствует
	var dateTaskNow time.Time
	var timeParseTaskDate time.Time
	if task.Date == "" {
		dateTaskNow, err = CheckDate(task.Date)
		if err != nil {
			writeError(w, "не верный формат даты", http.StatusBadRequest)
			log.Println("даты нет или пустая")
			return
		}
		task.Date = dateTaskNow.Format(FormatDate)
		timeParseTaskDate = dateTaskNow
	} else {
		dateTaskNow = time.Now()
		timeParseTaskDate, err = PrepareDate(task.Date)
		if err != nil {
			writeError(w, "не верный формат даты или дата не существует", http.StatusBadRequest)
			log.Println("строка не форматировалась в дату")
			return
		}
	}
	var nextDate string
	var id int64

	switch {
	case *task.Repeat == "" || task.Repeat == nil:
		if afterDate(timeParseTaskDate, dateTaskNow) {
			id, err = db.AddTask(task)
			if err != nil {
				writeError(w, "ошибка добавления задачи в бд", http.StatusBadRequest)
				log.Println("задача не добавлена в бд")
			}
		} else {
			task.Date = dateTaskNow.Format(FormatDate)
			id, err = db.AddTask(task)
			if err != nil {
				writeError(w, "ошибка добавления задачи в бд", http.StatusBadRequest)
				log.Println("задача не добавлена в бд")
			}
		}
	default:
		dateParts, err := PrepareRepeat(*task.Repeat)
		if err != nil {
			writeError(w, "не верный формат правила повторения", http.StatusBadRequest)
			return
		}
		if afterDate(timeParseTaskDate, dateTaskNow) {
			id, err = db.AddTask(task)
			if err != nil {
				writeError(w, "ошибка добавления задачи в бд", http.StatusBadRequest)
				log.Println("задача не добавлена в бд")
			}

		} else {
			nextDate, err = NextDate(dateTaskNow, dateParts, timeParseTaskDate)
			if err != nil {
				writeError(w, err.Error(), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
			id, err = db.AddTask(task)
			if err != nil {
				writeError(w, "ошибка добавления задачи в бд", http.StatusBadRequest)
			}

		}
	}
	idRes := strconv.Itoa(int(id))
	writeJson(w, map[string]string{"id": idRes}, http.StatusOK)
	timeParseTaskDate = dateTaskNow
	task.Date = timeParseTaskDate.Format(FormatDate)
}

//get - получаем task из бд по id из Url
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTasksByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJson(w, task, http.StatusOK)
}

// put - проверка входящ данных и обновление записи в бд
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task *db.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeError(w, "ошибка чтения", http.StatusBadRequest)
		log.Println("Не читается тело запроса")
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeError(w, "ошибка десериализации", http.StatusBadRequest)
		log.Println("JSON не десериализовался")
		return
	}
	if task.Title == nil || *task.Title == "" {
		writeError(w, `{"error":"заголовок обязателен и не должен быть пустым"}`, http.StatusBadRequest)
		log.Println("Заголовка нет или пустой")
		return
	}

	// dateTaskNow - текущее время, если task.Date пустое или отсутствует
	var dateTaskNow time.Time
	var timeParseTaskDate time.Time
	if task.Date == "" {
		dateTaskNow, err = CheckDate(task.Date)
		if err != nil {
			writeError(w, "не верный формат даты", http.StatusBadRequest)
			log.Println("даты нет или пустая")
			return
		}
		task.Date = dateTaskNow.Format(FormatDate)
		timeParseTaskDate = dateTaskNow
	} else {
		dateTaskNow = time.Now()
		timeParseTaskDate, err = PrepareDate(task.Date)
		if err != nil {
			writeError(w, "не верный формат даты или дата не существует", http.StatusBadRequest)
			log.Println("строка не форматировалась в дату")
			return
		}
	}
	var nextDate string
	switch {
	case *task.Repeat == "" || task.Repeat == nil:
		if afterDate(timeParseTaskDate, dateTaskNow) {
			err = db.UpdateTask(task)
			if err != nil {
				writeError(w, "ошибка обновления задачи в бд", http.StatusBadRequest)
				log.Println("задача не обновлена в бд")
			}
		} else {
			task.Date = dateTaskNow.Format(FormatDate)
			err = db.UpdateTask(task)
			if err != nil {
				writeError(w, "ошибка добавления задачи в бд", http.StatusBadRequest)
				log.Println("задача не обновлена в бд")
			}
		}
	default:
		dateParts, err := PrepareRepeat(*task.Repeat)
		if err != nil {
			writeError(w, "не верный формат правила повторения", http.StatusBadRequest)
			return
		}
		if afterDate(timeParseTaskDate, dateTaskNow) {
			err = db.UpdateTask(task)
			if err != nil {
				writeError(w, "ошибка обновления задачи в бд", http.StatusBadRequest)
				log.Println("задача не обновлена в бд")
			}
		} else {
			nextDate, err = NextDate(dateTaskNow, dateParts, timeParseTaskDate)
			if err != nil {
				writeError(w, err.Error(), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
			err = db.UpdateTask(task)
			if err != nil {
				writeError(w, "ошибка обновления задачи в бд", http.StatusBadRequest)
			}
		}
		writeJson(w, map[string]string{}, http.StatusOK)
	}
}

// получаем из бд записи заданным кол-вом
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetTasksByDate(10)
	if err != nil {
		writeError(w, "ошибка получения записей из базы данных", http.StatusBadRequest)
		return
	}
	writeJson(w, db.TasksResp{Tasks: tasks}, http.StatusOK)
}

//удаляет задачи
func TaskDelHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	writeJson(w, map[string]string{}, http.StatusOK)
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTasksByID(id)
	if err != nil {
		writeError(w, "ошибка обновления задачи в бд", http.StatusBadRequest)
		log.Println("задача не обновлена в бд")
	}
	if task.Repeat == nil || *task.Repeat == "" {
		TaskDelHandler(w, r)
		return
	} else {
		now, err := CheckDate(task.Date)
		if err != nil {
			writeError(w, "ошибка обновления задачи в бд", http.StatusBadRequest)
			log.Println("задача не обновлена в бд")
		}
		dateParts, err := PrepareRepeat(*task.Repeat)
		if err != nil {
			writeError(w, "ошибка обновления задачи в бд", http.StatusBadRequest)
			log.Println("задача не обновлена в бд")
		}
		dateStart, err := PrepareDate(task.Date)
		if err != nil {
			writeError(w, "ошибка обновления задачи в бд", http.StatusBadRequest)
			log.Println("задача не обновлена в бд")
		}
		finDate, err := NextDate(now, dateParts, dateStart)
		if err != nil {
			writeError(w, "ошибка обновления задачи в бд", http.StatusBadRequest)
			log.Println("задача не обновлена в бд")
		}
		task.Date = finDate
		err = db.UpdateDate(task, id)
		if err != nil {
			writeError(w, "ошибка обновления задачи в бд", http.StatusBadRequest)
			log.Println("задача не обновлена в бд")
		}
		writeJson(w, map[string]string{}, http.StatusOK)
	}
}

//формирует ответ в JSON, записывает заголовок и статус
func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// формирует ошибку в JSON, записывает заголовок и статус
func writeError(w http.ResponseWriter, errorMessage string, statusCode int) {
	writeJson(w, map[string]string{"error": errorMessage}, statusCode)
}
