package api

import (
	"net/http"
	"time"
)

const FormatDate string = "20060102"

// обрабатывает входящий запрос, возвращает след.дату исполнения задачи
func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	//дата выполнения получена из запроса (может быть любой)
	dstart := r.FormValue("date")
	if dstart == "" {
		http.Error(w, "Нет данных", http.StatusBadRequest)
	}
	repeat := r.FormValue("repeat")
	if repeat == "" {
		http.Error(w, "Нет данных", http.StatusBadRequest)
	}
	// дата "сегодняшняя", которую достаем из запроса для проверки на пустоту
	nowStr := r.FormValue("now")

	// дата "сегодняшняя", которую используем далее
	var now time.Time
	var err error
	// if nowStr == "" {
	// now = time.Now()
	// } 
		now, err = time.Parse(FormatDate, nowStr)
		if err != nil {
			panic(err)
		}
	
	rule, err := PrepareRepeat(repeat)
	if err != nil {
		http.Error(w, "неверный формат правила повторения", http.StatusBadRequest)
		return
	}
	date, err := PrepareDate(dstart)
	if err != nil {
		http.Error(w, "неверный формат времени, ожидается YYYYMMD", http.StatusBadRequest)
		return
	}
	nextDate, err := NextDate(now, rule, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		http.Error(w, "ошибка записи ответа", http.StatusBadRequest)
		return
	}
}
