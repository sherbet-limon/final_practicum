package api

import (
	"net/http"
	"time"
	
)
const FormatDate string = "20060102"
// обрабатывает входящий запрос, возвращает след.дату исполнения задачи
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	dstart := r.FormValue("date")
	if dstart == "" {
		http.Error(w, "Нет данных", http.StatusBadRequest)
	}
	repeat := r.FormValue("repeat")
	if repeat == "" {
		http.Error(w, "Нет данных", http.StatusBadRequest)
	}
	nowStr := r.FormValue("now")
	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	}else{
		now, err = time.Parse(FormatDate, nowStr)
		if err != nil {
			panic(err)
		}
	}
	date, time, err:= PrepareNextDate(repeat, dstart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	nextDate, err := NextDate(now, date, time)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(nextDate))
}