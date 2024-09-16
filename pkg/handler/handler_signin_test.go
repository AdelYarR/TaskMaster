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

func TestHandler_signIn(t *testing.T) {
	testCases := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		expectedStatusCode   int
		expectedResponseBody string
		mockServiceResponse  func() *mockService
	}{
		{
			name:      "Successful sign in",
			inputBody: `{"email":"sergeev2004@mailbox.ru", "password":"@valeriy"}`,
			inputUser: models.User{
				Email:    "sergeev2004@mailbox.ru",
				Password: "@valeriy",
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"JWT Token":""}`,
			mockServiceResponse:  nil,
		},
		{
			name:      "Decoding error",
			inputBody: `{"testerror":"test", "testerror":"test"}`,
			inputUser: models.User{
				Email:    "test",
				Password: "test",
			},
			expectedStatusCode:   400,
			expectedResponseBody: "Failed to decode json while singning in",
			mockServiceResponse:  nil,
		},
		{
			name:      "Incorrect Email",
			inputBody: `{"email":"error", "password":"@valeriy"}`,
			inputUser: models.User{
				Email:    "error",
				Password: "@valeriy",
			},
			expectedStatusCode:   400,
			expectedResponseBody: "Failed to sign in: incorrect email",
			mockServiceResponse: func() *mockService {
				return &mockService{signInFunc: func(models.User) (string, error) {
					return "", fmt.Errorf("IncorrectEmail")
				}}
			},
		},
		{
			name:      "Incorrect Password",
			inputBody: `{"email":"sergeev2004@mailbox.ru", "password":"short"}`,
			inputUser: models.User{
				Email:    "sergeev2004@mailbox.ru",
				Password: "short",
			},
			expectedStatusCode:   400,
			expectedResponseBody: "Failed to sign in: incorrect password",
			mockServiceResponse: func() *mockService {
				return &mockService{signInFunc: func(models.User) (string, error) {
					return "", fmt.Errorf("IncorrectPassword")
				}}
			},
		},
		{
			name:      "Service error",
			inputBody: `{"email":"sergeev2004@mailbox.ru", "password":"@valeriy"}`,
			inputUser: models.User{
				Email:    "sergeev2004@mailbox.ru",
				Password: "@valeriy",
			},
			expectedStatusCode:   500,
			expectedResponseBody: "Failed to create JWT Token",
			mockServiceResponse: func() *mockService {
				return &mockService{signInFunc: func(models.User) (string, error) {
					return "", fmt.Errorf("service error")
				}}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "localhost:8000/signin",
				bytes.NewBuffer([]byte(tc.inputBody)))

			var svc *mockService
			if tc.mockServiceResponse != nil {
				svc = tc.mockServiceResponse()
			} else {
				svc = &mockService{}
			}

			h := NewHandler(svc, slog.New(slog.NewJSONHandler(os.Stdout, nil)))
			hand := http.HandlerFunc(h.SignIn())
			hand.ServeHTTP(rec, req)

			actualBody := strings.TrimSpace(rec.Body.String())

			assert.Equal(t, tc.expectedStatusCode, rec.Code, "Status codes should match")
			assert.Equal(t, tc.expectedResponseBody, actualBody, "Body values should match")
		})
	}
}
