package routers

import (
	"net/http"

	"stakeholder-service/api/handlers"
	"stakeholder-service/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.Use(otelhttp.NewMiddleware("stakeholders-service"))
	r.Use(middleware.LogrusMiddleware)

	r.Mount("/profiles", handlers.Routes())

	// Metrics
	r.Handle("/metrics", promhttp.Handler())

	return r
}
