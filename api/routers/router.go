package routers

import (
	"net/http"

	"stakeholder-service/api/handlers"
	"stakeholder-service/middleware"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.Use(otelhttp.NewMiddleware("stakeholders-service"))
	r.Use(middleware.LogrusMiddleware)

	r.Mount("/profiles", handlers.Routes())

	return r
}
