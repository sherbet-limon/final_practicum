package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Функция возвращает след.дату исполнения задачи

func dateRequest(w http.ResponseWriter, r *http.Request) {
	dstart = r.FormValue("date")
	if dstart == ""{
		http.Error(w, "Нет данных", http.StatusBadRequest)
	}
	repeat = r.FormValue("repeat")
	if repeat == ""{
		http.Error(w, "Нет данных", http.StatusBadRequest)
	}
	now := time.Now()
	nextDate, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(nextDate))
}
//проверяет, что дата старта позже даты now
func arterDate(date time.Time, now time.Time) bool{
	if date.After(now){
		return true
	}
	return false
}
//парсит строку, переносит дату на указанное количество дней/год 
func NextDate(now time.Time, dstart string, repeat string) (string, error){
	dateStart, err:=time.Parse("20060102", dstart);
	if err !=nil {
		return "", errors.New("Неверный формат времени, ожидается YYYYMMDD")
	}  
	dateParts:=strings.Split(repeat, " ")
	if len(dateParts) < 1 {
		return "", errors.New("Неверный формат")
	}
	//var nextDate time.Time
	switch dateParts[0] { 
	case "d":
		if len(dateParts) < 2 {
			return "", errors.New("Неверный формат")
		}
		days,err:= strconv.Atoi(dateParts [1])
		if err != nil || days <= 0 || days > 400{
			return "", errors.New("Недопустимое количество дней, должно быть от 1 до 400")
		}
		for arterDate(dateStart, now)==false {
			dateStart = dateStart.AddDate(0, 0, days)
		}
	case "y":
		for arterDate(dateStart, now)==false {
			dateStart = dateStart.AddDate(1, 0, 0)
		}
	default:
		return "", errors.New("Недопустимый формат записи")	
	}
	nextDate:=dateStart.Format("20060102")
		
return nextDate, nil
}

func {
err:=http.HandleFunc(`/`, mainHandle); err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }
} 