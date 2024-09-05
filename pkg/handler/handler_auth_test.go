package handler

import (
	"TaskMaster/pkg/models"
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type mockService struct {
	signUpFunc func(models.User) (int, error)
	signInFunc func(http.ResponseWriter, *http.Request)
	tasksFunc  func(http.ResponseWriter, *http.Request, int)
}

func TestHandler_signUp(t *testing.T) {

	testCases := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Successful sign up",
			inputBody: `{"email":"testing@mail.ru", "password":"qwertytest"}`,
			inputUser: models.User{
				Email:    "testing@mail.ru",
				Password: "qwertytest",
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Decoding error",
			inputBody: `{"testerror":"testerror", "testerror":"testerror"}`,
			inputUser: models.User{
				Email:    "",
				Password: "",
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Validation error",
			inputBody: `{"email":"wrong_type", "password":"short"}`,
			inputUser: models.User{
				Email:    "wrong_type",
				Password: "short",
			},
			expectedStatusCode:   400,
			expectedResponseBody: "",
		},
	}

	h := NewHandler(&mockService{}, slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	hand := http.HandlerFunc(h.SignUp())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "localhost:8000/signup",
				bytes.NewBuffer([]byte(tc.inputBody)))
			hand.ServeHTTP(rec, req)
		})
	}
}
