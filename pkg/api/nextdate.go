package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const formatDate string = "20060102"

// при отсутствии даты, подставляет текущее время
func CheckDate(nowStr string) (time.Time, error) {
	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(formatDate, nowStr)
		if err != nil {
			panic(err)
		}
	}
	return now, nil
}

// проверяет, что дата старта позже даты now
func afterDate(date time.Time, nowStr time.Time) bool {
	return date.After(nowStr)
}

// парсит дату, делит repeat на части
func PrepareNextDate(repeat string, dstart string) ([]string, time.Time, error) {
	dateStart, err := time.Parse(formatDate, dstart)
	if err != nil {
		return nil, time.Time{}, errors.New("неверный формат времени, ожидается YYYYMMD")
	}
	dateParts := strings.Split(repeat, " ")
	if len(dateParts) < 1 {
		return nil, time.Time{}, errors.New("неверный формат")
	}
	return dateParts, dateStart, nil
}

// переносит дату на указанное количество дней/год
func NextDate(now time.Time, dateParts []string, dateStart time.Time) (string, error) {
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
