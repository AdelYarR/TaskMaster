package handler

import (
	"TaskMaster/pkg/models"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

type Service interface {
	SignUp(models.User) (int, error)
	SignIn(http.ResponseWriter, *http.Request)
	Tasks(http.ResponseWriter, *http.Request, int)
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(serv Service, log *slog.Logger) *Handler {
	return &Handler{
		service: serv,
		logger:  log,
	}
}

func (h *Handler) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Failed to decode json while signing up", http.StatusBadRequest)
			return
		}

		err = validateUser(user)
		if err != nil {
			http.Error(w, "Failed to sign up: wrong email or password", http.StatusBadRequest)
			return
		}

		id, err := h.service.SignUp(user)
		if err != nil {
			http.Error(w, "Failed to sign up", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}

func (h *Handler) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.service.SignIn(w, r)
	}
}

func (h *Handler) Tasks(w http.ResponseWriter, r *http.Request, userID int) {
	h.service.Tasks(w, r, userID)
}

func validateUser(user models.User) error {
	if user.Email == "" {
		return errors.New("email is required")
	}

	if len(user.Password) < 6 {
		return errors.New("password is too short")
	}

	if !strings.Contains(user.Email, "@") {
		return errors.New("you must enter an email domain")
	}

	return nil
}
