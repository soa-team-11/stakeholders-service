package main

import (
	"net/http"
	"stakeholder-service/api/routers"
	"stakeholder-service/utils"

	log "github.com/sirupsen/logrus"
)

func main() {
	router := routers.Router()

	port := utils.Getenv("PORT", "8080")

	log.Infof("Running services on PORT %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
