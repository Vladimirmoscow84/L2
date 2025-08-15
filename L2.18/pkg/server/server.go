package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"L2.18/pkg/service"
)

//структура для работы с хранилищем

type handler struct {
	storage service.Storage
}

// Newhandler - коструктор экземпляра handler
func Newhandler() handler {
	return handler{
		storage: service.NewStorage(),
	}
}

// Routers - регистрация хэндлеров в роутере
func (h *handler) Routers(mux *http.ServeMux) {
	mux.HandleFunc("/create_event", h.createEvent)
	mux.HandleFunc("/update_event", h.updateEvent)
	mux.HandleFunc("/delete_event", h.deleteEvent)
	mux.HandleFunc("/events_for_day", h.eventsForDay)
	mux.HandleFunc("/events_for_week", h.eventsForWeek)
	mux.HandleFunc("/events_for_month", h.eventsForMonth)

}

// unmarshalJSON - парсит и валидирует параметры запросов
func (h *handler) unmarshalJSON(r *http.Request) (*service.Event, error) {
	var event service.Event
	file, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	err = json.Unmarshal(file, &event)
	if err != nil {
		return nil, err
	}
	if event.UserId < 1 {
		return nil, errors.New("userId должен быть положительным")
	}
	return &event, nil
}

// succesResponse - возвращает успешный ответ
func (h *handler) succesResponse(w http.ResponseWriter, answer string, events []service.Event) {
	w.Header().Set("Content-Type", "application/json")
	result := service.Result{
		Message: answer,
		Results: events,
	}
	data, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		h.errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// errorResponse - возвращает JSON с описанием ошибки
func (h *handler) errorResponse(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	errMessage := err.Error()
	data, er := json.MarshalIndent(service.ErrMessage{Error: errMessage}, "", "\t")
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ошибка сериализации:" + errMessage))
		return
	}
	w.WriteHeader(status)
	w.Write(data)
}

// createEvent  - хендлер, который создает событие
func (h *handler) createEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		event, err := h.unmarshalJSON(r)
		if err != nil {
			h.errorResponse(w, err, http.StatusBadRequest)
			return
		}
		h.storage.CreateEvent(event)
		h.succesResponse(w, "событие создалось", []service.Event{*event})
		return
	}
	err := fmt.Errorf("использован неверный метод %s, необходимо  использовать метод POST", r.Method)
	h.errorResponse(w, err, http.StatusBadRequest)
}

// updaatevent - хенделр, который обновляет событие
func (h *handler) updateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		event, err := h.unmarshalJSON(r)
		if err != nil {
			h.errorResponse(w, err, http.StatusBadRequest)
		}
		err = h.storage.UpdateEvent(event)
		if err != nil {
			h.errorResponse(w, err, http.StatusBadRequest)
		}
		h.succesResponse(w, "событие изменено", []service.Event{*event})
		return
	}
	err := fmt.Errorf("использован неверный метод %s, необходио использовать метод POST", r.Method)
	h.errorResponse(w, err, http.StatusBadRequest)
}

// deleteEvent - хендлер, котрый удаляет событие
func (h *handler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			h.errorResponse(w, err, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var event *service.Event
		err = json.Unmarshal(data, &event)
		if err != nil {
			h.errorResponse(w, err, http.StatusBadRequest)
			return
		}
		h.succesResponse(w, "событие удалено", []service.Event{*event})
		return
	}
	err := fmt.Errorf("использован неверный метод %s, необходио использовать метод POST", r.Method)
	h.errorResponse(w, err, http.StatusBadRequest)
}

// eventsForDay - хендлер, показывающий события конкретного id в конкретный день
func (h *handler) eventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//из запроса получаем id и дату
		id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			h.errorResponse(w, err, http.StatusBadRequest)
			return
		}
		dateStr := r.URL.Query().Get("date")
		dateTime, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			h.errorResponse(w, err, http.StatusBadRequest)
			return
		}
		events := h.storage.EvensDay(id, dateTime)
		if len(events) == 0 {
			h.succesResponse(w, "нет событий на данную дату", events)
			return
		}
		h.succesResponse(w, "получен список событий на данную дату", events)
		return
	}
	err := fmt.Errorf("использован неверный метод %s, необходио использовать метод GET", r.Method)
	h.errorResponse(w, err, http.StatusBadRequest)
}

// eventsForWeek - хендлер, показывающий события конкретного id в конкретную днеделю
func (h *handler) eventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//из запроса получаем id и неделю
		id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			h.errorResponse(w, err, http.StatusBadRequest)
			return
		}
		dateStr := r.URL.Query().Get("date")
		dateTime, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			h.errorResponse(w, err, http.StatusBadRequest)
			return
		}
		events := h.storage.EventsWeek(id, dateTime)
		if len(events) == 0 {
			h.succesResponse(w, "нет событий в запрошенную неделю", events)
			return
		}
		h.succesResponse(w, "получен список событий в указанную неделю", events)
		return
	}
	err := fmt.Errorf("использован неверный метод %s, необходио использовать метод GET", r.Method)
	h.errorResponse(w, err, http.StatusBadRequest)
}

// eventsForMonth - ендлер, показывающий события конкретного id в конкретный месяц
func (h *handler) eventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Получаем id пользователя и дату из запроса
		id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			h.errorResponse(w, err, http.StatusBadRequest)
			return
		}
		dateStr := r.URL.Query().Get("date")
		dateTime, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			dateTime, err = time.Parse("2006-01", dateStr)
			if err != nil {
				h.errorResponse(w, err, http.StatusBadRequest)
				return
			}
		}
		events := h.storage.EventsMonth(id, dateTime)
		if len(events) == 0 {
			h.succesResponse(w, "нет событий в запрошенном месяце", events)
			return
		}
		h.succesResponse(w, "получен список событий в заданнный месяц", events)
		return
	}
	err := fmt.Errorf("использован неверный метод %s, необходио использовать метод GET", r.Method)
	h.errorResponse(w, err, http.StatusBadRequest)
}
