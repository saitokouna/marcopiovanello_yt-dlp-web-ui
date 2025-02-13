package configurator

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// App configurator REST handlers

func GetConfig(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(Instance().Config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SetConfig(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var req AppConfig

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Instance().setAppConfig(&req)

	if err := Instance().Persist(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("ok")
}

func setAppTitle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var req string

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	editField(w, func(c *AppConfig) {
		if req != "" {
			c.Title = req
		}
	})

	json.NewEncoder(w).Encode("ok")
}

func setBaseURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var req string

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	editField(w, func(c *AppConfig) {
		if req != "" {
			c.BaseURL = req
		}
	})

	json.NewEncoder(w).Encode("ok")
}

func editField(w http.ResponseWriter, editFunc func(c *AppConfig)) {
	editFunc(&Instance().Config)

	if err := Instance().Persist(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ApplyRouter() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", GetConfig)
		r.Post("/", SetConfig)
		r.Patch("/title", setAppTitle)
		r.Patch("/baseURL", setBaseURL)
	}
}
