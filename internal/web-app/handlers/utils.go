package handlers

import (
	"bytes"
	order "first-task/internal/entities/Order"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func FoundOrderTmpl(w http.ResponseWriter, ord *order.Order) {
	const op = "internal.web-app.handlers.FoundOrderTmpl"
	
	buf := bytes.NewBuffer([]byte{})
	err := tpl.ExecuteTemplate(buf, "order.html", ord)
	if err != nil {
		zap.L().Error(fmt.Sprintf("%s: %s", op, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := buf.WriteTo(w); err != nil {
		zap.L().Error(fmt.Sprintf("%s: err on writing template to page", op))
	}
}

func NotFoundOrderTmpl(w http.ResponseWriter) {
	const op = "internal.web-app.handlers.NotFoundOrderTmpl"

	buf := bytes.NewBuffer([]byte{})
	err := tpl.ExecuteTemplate(buf, "not-found-order.html", nil)
	if err != nil {
		zap.L().Error(fmt.Sprintf("%s: %s", op, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	if _, err := buf.WriteTo(w); err != nil {
		zap.L().Error("on writing template to page")
	}
}
