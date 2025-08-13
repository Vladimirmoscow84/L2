package service

import (
	"errors"
	"strings"
	"sync"
	"time"
)

// event - структура события
type Event struct {
	EventId int    `json:"event_id,omitempty"`
	UserId  int    `json:"usser_id"`
	Title   string `json:"title"`
	Notice  string `json:"notice"`
	Date    Date   `json:"date"`
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

// Result - структура для вывода результата запроса
type Result struct {
	Message string  `json:"message"`
	Results []Event `json:"results"`
}

// Storage - структура, в которой хрантся события
type Storage struct {
	storageResults map[int]*Event
	numberEvent    int
	mu             *sync.RWMutex
}

// errMessage  - структура для вывода ошибки запроса
type errMessage struct {
	Error string `json:"error"`
}
