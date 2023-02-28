// Package handler - реализует настройку маршрутов и управление.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"forward/internal/repository"
)

// NodeJSON - сущность для ответа запроса GET.
type NodeJSON struct {
	Success bool   `json:"success"`
	Count   string `json:"count"`
}

const (
	ContentTypeApplicationJSON = "application/json"
)

// NewRouter - создает роутер,настраивает маршруты.
func NewRouter(d *repository.Repository) chi.Router {
	router := chi.NewRouter()

	controller := newController(d)

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/check/{INN}", controller.Get)

	router.NotFound(NotFound())
	router.MethodNotAllowed(NotAllowed())

	return router
}

type Controller struct {
	Database *repository.Repository
}

// newController - функция-конструктор контролера хэндлера.
func newController(s *repository.Repository) *Controller {
	return &Controller{Database: s}
}

// Get - обрабатывает запрос маршрута GET /checks/{INN}
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	inn := chi.URLParam(r, "INN")
	if inn == "" {
		http.Error(w, "ErrNoEmptyINNParam", http.StatusBadRequest)
		return
	}
	var result NodeJSON
	count, err := c.Database.Get(inn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch count {
	case "":
		result.Success = true
		result.Count = "0"
	case "0":
		result.Success = true
		result.Count = "0"
	default:
		result.Success = false
		result.Count = count
	}
	body, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// NotFound - обработчик неподдерживаемых маршрутов.
func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("route does not exist"))
	}
}

// NotAllowed - обработчик неподдерживаемых методов.
func NotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("method does not allowed"))
	}
}
