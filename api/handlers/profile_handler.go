package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

	"stakeholder-service/internal/providers/cloudinary"
	"stakeholder-service/models"
	"stakeholder-service/services"
)

var (
	profileService = services.NewProfileService()
)

func Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json", "multipart/form-data"))

	r.Post("/", HandleCreate)
	r.Put("/", HandleUpdate)
	r.Get("/getByUserId", HandleGetByUserId)
	r.Get("/getRecommendationsByUserId", HandleGetRecommendationsByUserId)

	return r
}

func HandleUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("stakeholders-service")

	_, span := tracer.Start(ctx, "HandleUpdate")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")

	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		span.RecordError(err)
		http.Error(w, fmt.Sprintf(`{"message":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	jsonData := r.FormValue("profile")
	var profile models.Profile

	span.AddEvent("Profile json unmarshal")
	err = json.Unmarshal([]byte(jsonData), &profile)

	if err != nil {
		span.RecordError(err)
		http.Error(w, fmt.Sprintf(`{"message":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		span.AddEvent("Profile image upload")
		imageURL, err := cloudinary.UploadImage(file)

		if err != nil {
			span.RecordError(err)
			http.Error(w, fmt.Sprintf(`{"message":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		log.Infof("Uploaded file: %s\n", header.Filename)
		log.Info(imageURL)

		profile.ImageURL = imageURL
	}

	updatedProfile, err := profileService.Update(profile)
	if err != nil {
		span.RecordError(err)
		http.Error(w, fmt.Sprintf(`{"message":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(updatedProfile.ToJSON())
}

func HandleCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("stakeholders-service")

	_, span := tracer.Start(ctx, "HandleCreate")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		span.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		fmt.Fprintf(w, `{"message":"%s"}`, err.Error())
		return
	}

	var profile models.Profile

	span.AddEvent("Profile json unmarshal")
	err = json.Unmarshal(b, &profile)

	if err != nil {
		span.RecordError(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"%s"}`, err.Error())
		return
	}

	createdProfile, err := profileService.Create(ctx, profile)

	if err != nil {
		span.RecordError(err)
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

func HandleGetRecommendationsByUserId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDParam := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"invalid user_id format"}`)
		return
	}

	profiles, err := profileService.GetRecommendations(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"message":"failed to get recommendations: %s"}`, err.Error())
		return
	}

	// Convert profiles to JSON
	resp, err := json.Marshal(profiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"message":"failed to marshal response: %s"}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
