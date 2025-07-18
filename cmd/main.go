package main

import (
	"net/http"
	"stakeholder-service/api/routers"
	"stakeholder-service/utils"

	cld "stakeholder-service/internal/providers/cloudinary"

	log "github.com/sirupsen/logrus"
)

func main() {
	router := routers.Router()

	port := utils.Getenv("PORT", "8080")
	cloudinary_api := utils.Getenv("CLOUDINARY_APIKEY", "669436399163959")
	cloudinary_secret := utils.Getenv("CLOUDINARY_SECRET", "zIztHxdCJLJN1Onr2bO74dgQeEw")
	cloudinary_name := utils.Getenv("CLOUDINARY_NAME", "dslchettz")

	cld.Init(cloudinary_api, cloudinary_secret, cloudinary_name)

	log.Infof("Running services on PORT %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
