package service

import (
	"final_project/pkg/db"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat rule is empty")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid repeat rule")
	}

	switch parts[0] {
	case "d":
		return nextDay(date, now, parts)

	case "y":
		return nextYear(date, now, parts)

	case "w":
		return nextWeek(date, now, parts)

	case "m":
		return nextMonth(date, now, parts)

	default:
		return "", fmt.Errorf("invalid repeat rule: %s", repeat)
	}
}

func nextDay(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid day rule")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid day interval: %w", err)
	}

	if interval < 1 || interval > 400 {
		return "", fmt.Errorf("day interval must be between 1 and 400")
	}

	for {
		date = date.AddDate(0, 0, interval)

		if afterNow(date, now) {
			return date.Format(dateFormat), nil
		}
	}
}

func nextYear(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 1 {
		return "", fmt.Errorf("invalid yearly rule")
	}

	for {
		date = date.AddDate(1, 0, 0)

		if afterNow(date, now) {
			return date.Format(dateFormat), nil
		}
	}
}

func nextWeek(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid weekly rule")
	}

	weekdays, err := parseWeekdays(parts[1])
	if err != nil {
		return "", err
	}

	for {
		date = date.AddDate(0, 0, 1)

		if !afterNow(date, now) {
			continue
		}

		weekday := int(date.Weekday())

		if weekday == 0 {
			weekday = 7
		}

		if weekdays[weekday] {
			return date.Format(dateFormat), nil
		}
	}
}

func nextMonth(date, now time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("invalid monthly rule")
	}

	days, err := parseMonthDays(parts[1])
	if err != nil {
		return "", err
	}

	var months []int

	if len(parts) == 3 {
		months, err = parseMonths(parts[2])
		if err != nil {
			return "", err
		}
	}

	start := date

	for {
		year := date.Year()
		month := int(date.Month())

		// Проверяем текущий месяц.
		if len(months) == 0 || contains(months, month) {
			var candidates []time.Time

			for _, day := range days {
				candidate, ok := monthDate(year, month, day)
				if !ok {
					continue
				}

				if candidate.After(now) && candidate.After(start) {
					candidates = append(candidates, candidate)
				}
			}

			// Из всех подходящих дат берём самую раннюю.
			if len(candidates) > 0 {
				result := candidates[0]

				for _, candidate := range candidates[1:] {
					if candidate.Before(result) {
						result = candidate
					}
				}

				return result.Format(dateFormat), nil
			}
		}

		// Подходящей даты в текущем месяце нет —
		// переходим к следующему.
		date = date.AddDate(0, 1, 0)
	}
}

func parseWeekdays(value string) (map[int]bool, error) {
	result := make(map[int]bool)

	for _, item := range strings.Split(value, ",") {
		weekday, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("invalid weekday: %s", item)
		}

		if weekday < 1 || weekday > 7 {
			return nil, fmt.Errorf("weekday must be between 1 and 7")
		}

		result[weekday] = true
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no weekdays specified")
	}

	return result, nil
}

func parseMonthDays(value string) ([]int, error) {
	var days []int

	for _, item := range strings.Split(value, ",") {
		day, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("invalid month day: %s", item)
		}

		if day == 0 || day < -2 || day > 31 {
			return nil, fmt.Errorf("month day must be from 1 to 31, -1 or -2")
		}

		days = append(days, day)
	}

	if len(days) == 0 {
		return nil, fmt.Errorf("no month days specified")
	}

	return days, nil
}

func parseMonths(value string) ([]int, error) {
	var months []int

	for _, item := range strings.Split(value, ",") {
		month, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("invalid month: %s", item)
		}

		if month < 1 || month > 12 {
			return nil, fmt.Errorf("month must be between 1 and 12")
		}

		months = append(months, month)
	}

	if len(months) == 0 {
		return nil, fmt.Errorf("no months specified")
	}

	return months, nil
}

func monthDate(year, month, day int) (time.Time, bool) {
	lastDay := daysInMonth(year, month)

	if day == -1 {
		day = lastDay
	} else if day == -2 {
		if lastDay < 2 {
			return time.Time{}, false
		}

		day = lastDay - 1
	}

	if day < 1 || day > lastDay {
		return time.Time{}, false
	}

	return time.Date(
		year,
		time.Month(month),
		day,
		0, 0, 0, 0,
		time.UTC,
	), true
}

func daysInMonth(year, month int) int {
	firstDayNextMonth := time.Date(
		year,
		time.Month(month)+1,
		1,
		0, 0, 0, 0,
		time.UTC,
	)

	return firstDayNextMonth.AddDate(0, 0, -1).Day()
}

func contains(values []int, value int) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}

	return false
}

func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0,
		date.Location(),
	)

	nowOnly := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	return dateOnly.After(nowOnly)
}
func CheckDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	_, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	today := now.Format(dateFormat)

	if task.Date < today {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date = next
		}
	}

	return nil
}
