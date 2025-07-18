package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"stakeholder-service/models"
	"stakeholder-service/services"
)

var (
	profileService = services.NewProfileService()
)

func Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))

	r.Post("/", HandleCreate)
	r.Get("/getByUserId", HandleGetByUserId)

	return r
}

func HandleCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"message":"%s"}`, err.Error())
		return
	}

	var profile models.Profile
	err = json.Unmarshal(b, &profile)

	fmt.Println(profile.ToJSONString())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"%s"}`, err.Error())
		return
	}

	createdProfile, err := profileService.Create(profile)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"%s"}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(createdProfile.ToJSON())
}

func HandleGetByUserId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDParam := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"invalid user_id format"}`)
		return
	}

	profile, err := profileService.GetByUserId(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"message":"profile not found"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(profile.ToJSON())
}
