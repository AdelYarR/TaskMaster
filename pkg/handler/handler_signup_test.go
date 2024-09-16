package handler

import (
	"TaskMaster/pkg/models"
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockService struct {
	signUpFunc func(models.User) (int, error)
	signInFunc func(models.User) (string, error)
	tasksFunc  func(http.ResponseWriter, *http.Request, int)
}

func (m *mockService) SignUp(user models.User) (int, error) {
	if m.signUpFunc != nil {
		return m.signUpFunc(user)
	}
	return 0, nil
}

func (m *mockService) SignIn(user models.User) (string, error) {
	if m.signInFunc != nil {
		return m.signInFunc(user)
	}
	return "", nil
}

func (m *mockService) Tasks(w http.ResponseWriter, r *http.Request, id int) {
	if m.tasksFunc != nil {
		m.tasksFunc(w, r, id)
	}
}

func TestHandler_signUp(t *testing.T) {

	testCases := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		expectedStatusCode   int
		expectedResponseBody string
		mockServiceResponse  func() *mockService
	}{
		{
			name:      "Successful sign up",
			inputBody: `{"email":"testing@mail.ru", "password":"qwertytest"}`,
			inputUser: models.User{
				Email:    "testing@mail.ru",
				Password: "qwertytest",
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":0}`,
			mockServiceResponse: func() *mockService {
				return &mockService{signUpFunc: func(models.User) (int, error) {
					return 0, nil
				}}
			},
		},
		{
			name:      "Decoding error",
			inputBody: `{"testerror":"testerror", "testerror":"testerror"}`,
			inputUser: models.User{
				Email:    "",
				Password: "",
			},
			expectedStatusCode:   400,
			expectedResponseBody: "Failed to decode json while signing up",
			mockServiceResponse:  nil,
		},
		{
			name:      "Validation error",
			inputBody: `{"email":"wrong_type", "password":"short"}`,
			inputUser: models.User{
				Email:    "wrong_type",
				Password: "short",
			},
			expectedStatusCode:   400,
			expectedResponseBody: "Failed to sign up: wrong email or password",
			mockServiceResponse:  nil,
		},
		{
			name:      "Service error",
			inputBody: `{"email":"testing@mail.ru", "password":"qwertytest"}`,
			inputUser: models.User{
				Email:    "testing@mail.ru",
				Password: "qwertytest",
			},
			expectedStatusCode:   500,
			expectedResponseBody: "Failed to sign up",
			mockServiceResponse: func() *mockService {
				return &mockService{signUpFunc: func(models.User) (int, error) {
					return 0, fmt.Errorf("service error")
				}}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "localhost:8000/signup",
				bytes.NewBuffer([]byte(tc.inputBody)))

			var svc *mockService
			if tc.mockServiceResponse != nil {
				svc = tc.mockServiceResponse()
			} else {
				svc = &mockService{}
			}

			h := NewHandler(svc, slog.New(slog.NewJSONHandler(os.Stdout, nil)))
			hand := http.HandlerFunc(h.SignUp())
			hand.ServeHTTP(rec, req)

			actualBody := strings.TrimSpace(rec.Body.String())

			assert.Equal(t, tc.expectedStatusCode, rec.Code, "Status codes should match")
			assert.Equal(t, tc.expectedResponseBody, actualBody, "Body values should match")
		})
	}
}
