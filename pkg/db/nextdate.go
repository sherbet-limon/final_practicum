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
func AfterDate(date, now time.Time) bool {
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

// проверяет дату на соответствие ГГГГММДД и приводит строку к формату time
func PrepareDate(date string) (time.Time, error) {
	if len(date) != 8 {
		return time.Time{}, errors.New("дата должна быть в формате YYYYMMDD")
	}
	dateFormat, err := time.Parse(formatDate, date)
	if err != nil {
		return time.Time{}, errors.New("неверный формат даты, ожидается YYYYMMD")
	}
	if dateFormat.Format("20060102") != date {
		return time.Time{}, errors.New("несуществующая дата")
	}
	return dateFormat, nil
}

//проверяет repeat, разбивает на части. Если нет date, подставляет now. Возвращает nextDate
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
    if repeat == "" {
        return "", errors.New("правило повторения не задано")
    }
    date, err := CheckDate(dstart)
    if err != nil {
        return "", errors.New("неверный формат dstart, ожидается YYYYMMDD")
    }
    parts := strings.Split(repeat, " ")
    if len(parts) < 2 {
        return "", errors.New("неверный формат repeat")
    }
    rule := parts[0]
    switch rule {
    case "d":
        if len(parts) != 2 {
            return "", errors.New("не указан интервал в днях")
        }
        interval, err := strconv.Atoi(parts[1])
        if err != nil {
            return "", errors.New("интервал дней должен быть числом")
        }
        if interval <= 0 || interval > 400 {
            return "", errors.New("количество дней должно быть от 1 до 400")
        }
        // прибавляем интервал дней, пока дата не станет больше now
        for {
            date = date.AddDate(0, 0, interval)
            if AfterDate(date, now) {
                break
            }
        }
    case "y":
        if len(parts) != 1 {
            return "", errors.New("неверный формат правила для y")
        }
        dateOrig := date
        for {
            date = date.AddDate(1, 0, 0)
            if AfterDate(date, now) {
                break
            }
            // Предотвратим бесконечный цикл на случай неверных данных:
            if date.Equal(dateOrig) {
                return "", errors.New("ошибка вычисления следующей даты для 'y'")
            }
        }
    default:
        return "", errors.New("недопустимый формат правила повторения")
    }

    return date.Format(formatDate), nil
}


// делит repeat из запроса на части
// func PrepareRepeat(repeat string) ([]string, error) {
// 	dateParts := strings.Split(repeat, " ")
// 	if len(dateParts)  {
// 		return nil, errors.New("неверный формат правила повторения")
// 	}
// 	return dateParts, nil
// }
// //определяет количество дней переноса, dstart - из url, checkDate-дата из url или now
// func Rule(repeat, dstart string) (string, error){
// 	dateParts, err:= PrepareRepeat(repeat)
// 	if err!=nil{
// 		return "", errors.New("ошибка чтения правила повторения")
// 	}
// 	points, err:= strconv.Atoi(dateParts[1])
// 	if err!=nil{
// 		return "", errors.New("ошибка форматирования правила повторения")
// 	}
// 	if points <= 0 || points > 400{
// 		return "", errors.New("количество дней должно быть от 1 до 400")
// 		}
// 	checkDate, err := CheckDate(dstart)
// 	if err!=nil{
// 		return "", errors.New("ошибка проверки даты")
// 	}
// 	dataFin, err := NextDate(checkDate, points, dateParts, dstart)
// 	if err!=nil{
// 		return "", errors.New("ошибка переноса даты")
// 	}
// 	return dataFin, nil
// }
// //переносит дату на указанное количество дней/год. dstart - из url
// func NextDate(now time.Time, points int, dateParts []string, dstart string) (string, error) {
// 	var nextDate time.Time
// 	switch dateParts[0]{
// 	case "d":
// 		for {
// 			nextDate = now.AddDate(0, 0, points)
// 			if AfterDate(nextDate, now) {
// 				break
// 			}
// 		}
// 	case "y":
// 		for {
// 			nextDate = now.AddDate(1, 0, 0)
// 			if AfterDate(nextDate, now) {
// 				break
// 			}
// 		} 
// 	default:
// 		return "", errors.New("недопустимый формат правила повторения")
// 	}
// 	nextDateStr := nextDate.Format(formatDate)
// 	return nextDateStr, nil
// }
