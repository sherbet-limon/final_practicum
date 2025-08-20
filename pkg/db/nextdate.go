package db

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const formatDate string = "20060102"

// при отсутствии даты, подставляет текущее время
func CheckDate(date string) (time.Time, error) {
	now := time.Now()
	if date == "" {
		return now, nil
	} 
	now, err := PrepareDate(date)
	if err != nil {
		return time.Time{}, errors.New("неверный формат даты, ожидается YYYYMMD")
	}
	return now, nil
}

// проверяет, что дата из запроса позже даты текущей
func AfterDate(date time.Time, now time.Time) bool {
	dateF:= date.Format(formatDate)
	nowDate := now.Format(formatDate)
	if nowDate == dateF { 	
		return true 		
	} 				
	if date.After(now) {
		return true
	}
	return false
}

// приводит строку к формату даты
func PrepareDate(dataStr string) (time.Time, error) {
	if len(dataStr) != 8 {
		return time.Time{}, errors.New("дата должна быть в формате YYYYMMDD")
	}
	date, err := time.Parse(formatDate, dataStr)
	if err != nil {
		return time.Time{}, errors.New("неверный формат даты, ожидается YYYYMMD")
	}
	if date.Format("20060102") != dataStr {
		return time.Time{}, errors.New("несуществующая дата")
	}
	return date, nil
}

// делит repeat из запроса на части
func PrepareRepeat(repeat string) ([]string, error) {
	dateParts := strings.Split(repeat, " ")
	if len(dateParts) < 1 {
		return nil, errors.New("неверный формат правила повторения")
	}
	return dateParts, nil
}

//определяет количество дней переноса, dstart - из url, checkDate-дата из url или now
func Rule(repeat, dstart string) (string, error){
	dateParts, err:= PrepareRepeat(repeat)
	if err!=nil{
		return "", errors.New("ошибка чтения правила повторения")
	}
	points, err:= strconv.Atoi(dateParts[1])
	if err!=nil{
		return "", errors.New("ошибка форматирования правила повторения")
	}
	if points <= 0 || points > 400{
		return "", errors.New("количество дней должно быть от 1 до 400")
		}
	checkDate, err := CheckDate(dstart)
	if err!=nil{
		return "", errors.New("ошибка проверки даты")
	}
	dataFin, err := NextDate(checkDate, points, dateParts, dstart)
	if err!=nil{
		return "", errors.New("ошибка переноса даты")
	}
	return dataFin, nil
}
//переносит дату на указанное количество дней/год. dstart - из url
func NextDate(now time.Time, points int, dateParts []string, dstart string) (string, error) {
	var nextDate time.Time
	switch dateParts[0]{
	case "d":
		for {
			nextDate = now.AddDate(0, 0, points)
			if AfterDate(nextDate, now) {
				break
			}
		}
	case "y":
		for {
			nextDate = now.AddDate(1, 0, 0)
			if AfterDate(nextDate, now) {
				break
			}
		} 
	default:
		return "", errors.New("недопустимый формат правила повторения")
	}
	nextDateStr := nextDate.Format(formatDate)
	return nextDateStr, nil
}
