package handlers

import (
	"encoding/json"
	"errors"
	"first-task/internal/storage"
	"first-task/internal/storage/postgres"
	"fmt"
	"net/http"
	"text/template"

	"go.uber.org/zap"
)

type Handler struct {
	str storage.Storager
}

func NewHandler(str storage.Storager) *Handler {
	return &Handler{str}
}

type ErrorResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

var StatusNotFound = "not found"
var StatusBadRequest = "bad request"
var StatusInternalServerError = "internal server error"

func (h *Handler) HandleMainPage(w http.ResponseWriter, r *http.Request) {
	const op = "internal.web-app.handlers.HandleMainPage"

	tmpl, err := template.ParseFiles("templates/main.html")
	if err != nil {
		zap.L().Error(fmt.Sprintf("%s: %s", op, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func (h *Handler) HandleFindOrder(w http.ResponseWriter, r *http.Request) {
	const op = "internal.web-app.handlers.HandleFindOrder"
	r.ParseForm()
	orderUID := r.FormValue("order_uid")

	ord, err := h.str.FindOrder(orderUID)
	if errors.Is(err, postgres.ErrNotFound) {
		NotFoundOrderTmpl(w)
		return
	} else if err != nil {
		zap.L().Error(fmt.Sprintf("%s: %s", op, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	FoundOrderTmpl(w, ord)
}

// API

// @Summary GetOrder
// @Tags Order
// @Description get order
// @Produce json
// @Param order_uid path string true "Уникальный номер заказа"
// @Success 200 {object} order.Order "Успешный запрос"
// @Failure 404 {object} ErrorResponse "Заказ не найден"
// @Failure 400 {object} ErrorResponse "Плохой запрос"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /order/{order_uid} [get]
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	const op = "internal.web-app.handlers.HandleOrderPage"

	w.Header().Set("Content-Type", "application/json")

	orderUID := r.PathValue("order_uid")
	if orderUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status: StatusBadRequest,
			Code:   http.StatusBadRequest,
		})
		return
	}
	ord, err := h.str.FindOrder(orderUID)
	if errors.Is(err, postgres.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status: StatusNotFound,
			Code:   http.StatusNotFound,
		})
		return
	} else if err != nil {
		zap.L().Error(fmt.Sprintf("%s: %s", op, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status: StatusInternalServerError,
			Code:   http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ord)
}
