package webapp

import (
	"first-task/internal/storage"
	"first-task/internal/web-app/handlers"
	"net/http"
)

func Handle(str storage.Storager) http.Handler {
	mux := http.NewServeMux()

	h := handlers.NewHandler(str)

	mux.HandleFunc("/", h.HandleMainPage)
	mux.HandleFunc("/order/{order_uid}", h.HandleOrderPage)
	mux.HandleFunc("/find-order", h.HandleFindOrder)

	return mux
}
