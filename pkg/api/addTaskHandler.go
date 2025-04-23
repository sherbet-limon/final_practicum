package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"go1f/pkg/db"
)

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
			if afterDate(timeParseTaskDate, dateTaskNow){
				id, err = AddTask(task)
				if err != nil {
					writeError(w, "ошибка добавления задачи в бд", http.StatusBadRequest)
					log.Println("задача не добавлена в бд")
				}
			} else {
			task.Date = dateTaskNow.Format(FormatDate)
			id, err = AddTask(task)
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
			if afterDate(timeParseTaskDate, dateTaskNow){
				id, err = AddTask(task)
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
			id, err = AddTask(task)
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
//формирует ответ в JSON, записывает заголовок и статус
func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, errorMessage string, statusCode int) {
	writeJson(w, map[string]string{"error": errorMessage}, statusCode)
}
