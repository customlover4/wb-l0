package webapp

import (
	"first-task/internal/storage"
	"first-task/internal/web-app/handlers"
	"net/http"

	httpSwager "github.com/swaggo/http-swagger"
)

func Handle(str storage.Storager) http.Handler {
	mux := http.NewServeMux()

	h := handlers.NewHandler(str)

	// swagger
	mux.HandleFunc("/swagger/", httpSwager.WrapHandler)

	mux.HandleFunc("GET /order/{order_uid}", h.GetOrder)
	mux.HandleFunc("GET /find-order", h.HandleFindOrder)
	mux.HandleFunc("/", h.HandleMainPage)

	return mux
}
