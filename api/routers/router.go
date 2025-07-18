package routers

import (
	"net/http"

	"stakeholder-service/api/handlers"
	"stakeholder-service/middleware"

	"github.com/go-chi/chi/v5"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.LogrusMiddleware)

	r.Mount("/profiles", handlers.Routes())

	return r
}
