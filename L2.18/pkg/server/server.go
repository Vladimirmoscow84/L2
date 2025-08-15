package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"L2.18/pkg/service"
	"github.com/go-chi/chi"
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

// Routers - решистрация хэндлеров в роутере
func (h *handler) Routers(chi *chi.Mux) {
	chi.HandleFunc("/create event", h.createEvent)
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
	err := fmt.Errorf("использован неверный метод %s, необходим метод POST", r.Method)
	h.errorResponse(w, err, http.StatusBadRequest)
}
