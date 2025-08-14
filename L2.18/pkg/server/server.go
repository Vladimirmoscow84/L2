package server

import (
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
func (h *handler) Routers(chi *chi.Mux){
	chi.HandleFunc("/create event", h.createEvent)
}

//Хендлер createEvent создает событие
func(h *handler) createEvent(w http.ResponseWriter, r *http.Request){
	if r.Method=="POST"{
		event,err:=
	}
}
