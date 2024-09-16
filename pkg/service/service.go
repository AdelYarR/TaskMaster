package service

import ("TaskMaster/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type PGRepo interface {
	SignUp(string, string) (int, error)
	SignIn(string) (int, string, error)
	GetTasks(int) ([]models.Task, error)
	AddTask(models.Task, int) (error)
}

type Service struct {
	repository PGRepo
	logger     *slog.Logger
}

func NewService(repo PGRepo, log *slog.Logger) *Service {
	return &Service{
		repository: repo,
		logger:     log,
	}
}

func (s *Service) SignUp(user models.User) (int, error) {

	hashed_Password, err := s.HashPassword(user.Password)
	if err != nil {
		s.logger.Error("Failed to hash password while signing up")
		return 0, err
	}

	id, err := s.repository.SignUp(user.Email, hashed_Password)
	if err != nil {
		s.logger.Error("Failed to sign up")
		return 0, err
	}

	s.logger.Info(
		"New account has been created successfully",
		"UserID", id,
	)

	return id, nil
}

func (s *Service) SignIn(user models.User) (string, error) {

	user_id, hash, err := s.repository.SignIn(user.Email)
	if err != nil {
		s.logger.Error("Failed to sign in: email is incorrect")
		return "", errors.New("IncorrectEmail")
	}

	if !s.CheckPasswordHash(hash, user.Password) {
		s.logger.Error("Failed to sign in: password is incorrect")
		return "", errors.New("IncorrectPassword")
	}

	token, err := s.CreateJWTToken(user_id)
	if err != nil {
		s.logger.Error("Failed to create token")
		return "", err
	}

	s.logger.Info(
		"Signing in has been done successfully",
		"JWT Token: ", token,
	)

	return token, nil
}

func (s *Service) Tasks(w http.ResponseWriter, r *http.Request, userID int) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := s.repository.GetTasks(userID)
		if err != nil {
			s.logger.Error("Failed to get tasks")
		}

		err = json.NewEncoder(w).Encode(tasks)
		if err != nil {
			fmt.Fprintf(w, "Failed to encode while getting tasks: %v", err)
			return
		}
	case http.MethodPost:
		var task models.Task

		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			fmt.Fprintf(w, "Failed to encode while getting tasks: %v", err)
			return
		}
		
		err = s.repository.AddTask(task, userID)
		if err != nil {
			s.logger.Error("Failed to add new task")
			return
		}

		fmt.Fprintf(w, "Successfully added new task.")
	}
}

func (s *Service) HashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	return string(bytes), err
}

func (s *Service) CheckPasswordHash(hash string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

var jwtKey = []byte("ekjgerj5et45F@wEDFRge$*riwe934urHsajd*W!@ffklgjmklVWbdsklcnkJBJKFSWGVWEF")

func (s *Service) CreateJWTToken(userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, err
}
