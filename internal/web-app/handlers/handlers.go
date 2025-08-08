package handlers

import (
	"encoding/json"
	"errors"
	order "first-task/internal/entities/Order"
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

func NotFoundOrderTmpl(w http.ResponseWriter) {
	const op = "internal.web-app.handlers.NotFoundOrderTmpl"

	tmpl, err := template.ParseFiles("templates/not-found-order.html")
	if err != nil {
		zap.L().Error(fmt.Sprintf("%s: %s", op, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// API
type Answer struct {
	Status string `json:"status"`
	Code   int    `json:"code"`

	Order *order.Order `json:"order,omitempty"`
}

func ResponseWithSuccess(w http.ResponseWriter, ord *order.Order) {
	const op = "internal.web-app.handlers.ResponseWithSuccess"

	answ := Answer{}

	answ.Code = http.StatusOK
	answ.Status = "found"
	answ.Order = ord

	res, err := json.MarshalIndent(answ, " ", "  ")
	if err != nil {
		zap.L().Error(fmt.Sprintf("%s: %s", op, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func ResponseWithError(w http.ResponseWriter) {
	const op = "internal.web-app.handlers.ResponseWithError"

	answ := Answer{}

	answ.Code = http.StatusNotFound
	answ.Status = "can't find this order"

	res, err := json.MarshalIndent(answ, " ", "  ")
	if err != nil {
		zap.L().Error(fmt.Sprintf("%s: %s", op, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func (h *Handler) HandleOrderPage(w http.ResponseWriter, r *http.Request) {
	const op = "internal.web-app.handlers.HandleOrderPage"

	orderUID := r.PathValue("order_uid")
	if orderUID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	ord, err := h.str.FindOrder(orderUID)
	if errors.Is(err, postgres.ErrNotFound) {
		ResponseWithError(w)
		return
	} else if err != nil {
		zap.L().Error(fmt.Sprintf("%s: %s", op, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ResponseWithSuccess(w, ord)
}
