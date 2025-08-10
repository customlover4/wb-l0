package handlers

import (
	order "first-task/internal/entities/Order"
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
