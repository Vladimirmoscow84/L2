package service

import (
	"errors"
	"strings"
	"time"
)

// event - структура события
type event struct {
	EventId int
	UserId  int
	Title   string
	Note    string
	Date    Date
}

// Date - структура для распарсивания дат, введенных пользователем
type Date struct {
	time.Time
}

// UnmarshalJSON - парсит переданные пользователем даты
func (t *Date) UnmarshalJSON(data string) error {

	if data == "" || data == `""` {
		*t = Date{time.Now()}
		return nil
	}

	timeStr := strings.ReplaceAll(data, `"`, "")
	parsTime, err := time.Parse("2006-01-02T15:04", timeStr)
	if err != nil {
		parsTime, err = time.Parse("2006-01-02T15:04:00Z", timeStr)
		if err != nil {
			parsTime, err = time.Parse("2006-01-02", timeStr)
			if err != nil {
				return errors.New("неправильный формат времени, попробуте ввести время в формате 2006-01-02T15:04")
			}
		}
	}
	*t = Date{parsTime}
	return nil

}
