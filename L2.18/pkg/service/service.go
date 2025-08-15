package service

import (
	"errors"
	"fmt"
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
func (t *Date) UnmarshalJSON(data []byte) error {

	if string(data) == "" || string(data) == `""` {
		*t = Date{time.Now()}
		return nil
	}

	timeStr := strings.ReplaceAll(string(data), `"`, "")
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

// Storage - структура, в которой хранятся события
type Storage struct {
	storageEvents map[int]*Event
	numberEvent   int
	mu            sync.RWMutex
}

// errMessage  - структура для вывода ошибки запроса
type ErrMessage struct {
	Error string `json:"error"`
}

// NewStorage - конструктор для storage
func NewStorage() Storage {
	return Storage{
		storageEvents: make(map[int]*Event),
		numberEvent:   1,
		mu:            sync.RWMutex{},
	}
}

// CreateEvent - создает событие с уникальным ID и отпарвляет в хранилище
func (s *Storage) CreateEvent(e *Event) {
	s.mu.Lock()
	e.EventId = s.numberEvent
	s.storageEvents[e.EventId] = e
	s.numberEvent++
	s.mu.Unlock()

}

// UpdateEvent - обновляет событие в хранилище
func (s *Storage) UpdateEvent(e *Event) error {
	s.mu.Lock()
	if _, ok := s.storageEvents[e.EventId]; !ok {
		return fmt.Errorf("собыие %d не найдено", e.EventId)
	}
	s.storageEvents[e.EventId] = e
	s.mu.Unlock()
	return nil
}

// DeleteEvent - удалет событие из хранилища
func (s *Storage) DeleteEvent(id int) (e *Event, err error) {
	s.mu.Lock()
	if _, ok := s.storageEvents[id]; !ok {
		return nil, fmt.Errorf("события с id %d не существует", id)
	}
	e = s.storageEvents[id]
	delete(s.storageEvents, id)
	s.mu.Unlock()
	return e, nil

}

// EvensDay - получает события в конкретный день по конкретному id
func (s *Storage) EvensDay(userId int, date time.Time) []Event {
	eventsSlice := make([]Event, 0)
	s.mu.RLock()

	//при итерациипо хранилищу вынимаем совпадения по дате и id и добавлем в слайс
	for _, event := range s.storageEvents {
		if event.UserId == userId && event.Date.Day() == date.Day() && event.Date.Month() == date.Month() && event.Date.Year() == date.Year() {
			eventsSlice = append(eventsSlice, *event)
		}
	}
	s.mu.RUnlock()
	return eventsSlice
}

// EventsWeek - получает события в конткретную неделю по конкретному id
func (s *Storage) EventsWeek(userId int, date time.Time) []Event {
	eventsSlice := make([]Event, 0)
	s.mu.RLock()
	yearEvent, weekEvent := date.ISOWeek()

	//при итерации по хранилищу вынимаем совпадения по дате и id и добавлем в слайс
	for _, event := range s.storageEvents {
		yearWant, weekWant := event.Date.ISOWeek()
		if event.UserId == userId && yearEvent == yearWant && weekEvent == weekWant {
			eventsSlice = append(eventsSlice, *event)
		}

	}
	s.mu.RUnlock()
	return eventsSlice
}

// EventsMonth - получает события в конкретный месяц по конкретному id
func (s *Storage) EventsMonth(userId int, date time.Time) []Event {
	eventsSlice := make([]Event, 0)
	s.mu.RLock()

	//при итерации по хранилищу вынимаем совпадения по месяцу и по конкретному id
	for _, event := range s.storageEvents {
		if event.UserId == userId && event.Date.Year() == date.Year() && event.Date.Month() == date.Month() {
			eventsSlice = append(eventsSlice, *event)
		}
	}
	s.mu.RUnlock()
	return eventsSlice
}
