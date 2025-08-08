package handlers

import (
	"errors"
	order "first-task/internal/entities/Order"
	"first-task/internal/storage/postgres"
	"fmt"
	"net/http"
	"text/template"

	"go.uber.org/zap"
)

func FoundOrderTmpl(w http.ResponseWriter, ord *order.Order) {
	const op = "internal.web-app.handlers.NotFoundOrderTmpl"

	tmpl, err := template.ParseFiles("templates/order.html")
	if err != nil {
		zap.L().Error(fmt.Sprintf("%s: %s", op, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, ord)
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
