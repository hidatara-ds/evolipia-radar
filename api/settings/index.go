package settings

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hidatara-ds/evolipia-radar/pkg/api"
	"github.com/hidatara-ds/evolipia-radar/pkg/config"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
	"github.com/hidatara-ds/evolipia-radar/pkg/models"
)

// Handler handles the /v1/settings route on Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	api.EnableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	cfg := config.Load()
	database, err := db.New(cfg)
	if err != nil {
		log.Printf("[VERCEL SETTINGS] DB Connection failed: %v", err)
		http.Error(w, `{"error":"database connection failed"}`, http.StatusInternalServerError)
		return
	}
	defer database.Close()

	repo := db.NewSettingRepository(database)

	switch r.Method {
	case "GET":
		keys := []string{"x_api_key", "threads_api_key", "openrouter_api_key"}
		settings := make(map[string]string)
		for _, key := range keys {
			val, _ := repo.Get(r.Context(), key) // ignore error, return empty string
			settings[key] = val
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)

	case "POST":
		var settings map[string]string
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		for key, val := range settings {
			// Basic validation/filtering of keys
			if key == "x_api_key" || key == "threads_api_key" || key == "openrouter_api_key" {
				err := repo.Set(r.Context(), models.Setting{
					Key:   key,
					Value: val,
				})
				if err != nil {
					log.Printf("[VERCEL SETTINGS] Failed to set %s: %v", key, err)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "saved"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
