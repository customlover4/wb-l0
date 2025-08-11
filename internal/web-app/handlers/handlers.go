package handlers

import (
	"encoding/json"
	"errors"
	order "first-task/internal/entities/Order"
	"first-task/internal/storage"
	"first-task/internal/templates"
	"fmt"
	"net/http"
	"text/template"

	"go.uber.org/zap"
)

type OrderGetter interface {
	FindOrder(orderUID string) (*order.Order, error)
}

type ErrorResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

var StatusNotFound = "not found"
var StatusBadRequest = "bad request"
var StatusInternalServerError = "internal server error"

var tpl = template.Must(templates.LoadTemplates())

func MainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.web-app.handlers.HandleMainPage"

		err := tpl.ExecuteTemplate(w, "main.html", nil)
		if err != nil {
			zap.L().Error(fmt.Sprintf("%s: %s", op, err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func FindOrder(str OrderGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.web-app.handlers.HandleFindOrder"
		
		r.ParseForm()
		orderUID := r.FormValue("order_uid")

		ord, err := str.FindOrder(orderUID)
		if errors.Is(err, storage.ErrNotFound) {
			NotFoundOrderTmpl(w)
			return
		} else if err != nil {
			zap.L().Error(fmt.Sprintf("%s: %s", op, err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		FoundOrderTmpl(w, ord)
	}
}

// API

// @Summary FindOrderAPI
// @Tags Order
// @Description get order
// @Produce json
// @Param order_uid path string true "Уникальный номер заказа"
// @Success 200 {object} order.Order "Успешный запрос"
// @Failure 404 {object} ErrorResponse "Заказ не найден"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /order/{order_uid} [get]
func FindOrderAPI(str OrderGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.web-app.handlers.HandleOrderPage"

		w.Header().Set("Content-Type", "application/json")

		orderUID := r.PathValue("order_uid")
		ord, err := str.FindOrder(orderUID)
		if errors.Is(err, storage.ErrNotFound) {
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

}
