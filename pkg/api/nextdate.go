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
	}
	now, err = time.Parse(formatDate, nowStr)
	if err != nil {
		panic(err)
	}
	return now, nil
}

// проверяет, что дата из запроса позже даты текущей
func afterDate(dateReq time.Time, now time.Time) bool {
	dateR := dateReq.Format(FormatDate)
	dateNow := now.Format(FormatDate)
	if dateNow == dateR { //пришлось добавить такую проверку из-за разногласий
		return true //today в тестах (присваивалась текущая дата и 00:00ч без секунд)
	} //и получаемых через time.Now()(присваивалась текущая дата и время до секунд)
	if dateReq.After(now) {
		return true
	}
	return false
}

// приводит строку к формату даты
func PrepareDate(dstart string) (time.Time, error) {
	if len(dstart) != 8 {
		return time.Time{}, errors.New("дата должна быть в формате YYYYMMDD")
	}
	dateStart, err := time.Parse(formatDate, dstart)
	if err != nil {
		return time.Time{}, errors.New("неверный формат даты, ожидается YYYYMMD")
	}
	if dateStart.Format("20060102") != dstart {
		return time.Time{}, errors.New("несуществующая дата")
	}
	return dateStart, nil
}

// делит repeat из запроса на части
func PrepareRepeat(repeat string) ([]string, error) {
	dateParts := strings.Split(repeat, " ")
	if len(dateParts) < 1 {
		return nil, errors.New("неверный формат правила повторения")
	}
	return dateParts, nil
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
		return "", errors.New("недопустимый формат правила повторения")
	}
	nextDate := dateStart.Format(formatDate)
	return nextDate, nil
}
