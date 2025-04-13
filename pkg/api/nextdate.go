package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const formatDate string = "20060102"

// Функция возвращает след.дату исполнения задачи
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
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
		now, err = time.Parse(formatDate, nowStr)
		if err != nil {
			panic(err)
		}
	}
	
	nextDate, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(nextDate))
}


// проверяет, что дата старта позже даты now
func afterDate(date time.Time, nowStr time.Time) bool {
	return date.After(nowStr)
}

// парсит строку, переносит дату на указанное количество дней/год
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	dateStart, err := time.Parse(formatDate, dstart)
	if err != nil {
		return "", errors.New("неверный формат времени, ожидается YYYYMMDD")
	}
	dateParts := strings.Split(repeat, " ")
	if len(dateParts) < 1 {
		return "", errors.New("неверный формат")
	}
	//var nextDate time.Time
	switch dateParts[0] {
	case "d":
		if len(dateParts) < 2 {
			return "", errors.New("неверный формат")
		}
		days, err := strconv.Atoi(dateParts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("недопустимое количество дней, должно быть от 1 до 400")
		}
		for {
			dateStart = dateStart.AddDate(0, 0, days)
			if afterDate(dateStart, now) {
				break
			}
		}
	case "y":
		for {
			dateStart = dateStart.AddDate(1, 0, 0)
			if afterDate(dateStart, now) {
				break
			}
		}
	default:
		return "", errors.New("недопустимый формат записи")
	}
	nextDate := dateStart.Format(formatDate)

	return nextDate, nil
}
