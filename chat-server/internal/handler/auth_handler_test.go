package handler

import (
	"bytes"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_handler "github.com/vavelour/chat/internal/handler/mocks"
	"github.com/vavelour/chat/internal/handler/request"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthHandler_Register(t *testing.T) {
	type mockBehavior func(s *mock_handler.MockAuthService, username, password string)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           request.RegisterRequest
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "ok",
			inputBody: `{"username": "tester","password": "123"}`,
			inputUser: request.RegisterRequest{Username: "tester", Password: "123"},
			mockBehavior: func(s *mock_handler.MockAuthService, username, password string) {
				s.EXPECT().CreateUser(username, password).Return("sometoken", nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: `{"response":"user created","token":"sometoken"}`,
		},
		{
			name:                "empty_fields",
			inputBody:           `{"username": "","password": ""}`,
			inputUser:           request.RegisterRequest{Username: "", Password: ""},
			mockBehavior:        func(s *mock_handler.MockAuthService, username, password string) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"error":"Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag\nKey: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:      "service_failure_1",
			inputBody: `{"username": "user","password": "123"}`,
			inputUser: request.RegisterRequest{Username: "user", Password: "123"},
			mockBehavior: func(s *mock_handler.MockAuthService, username, password string) {
				s.EXPECT().CreateUser(username, password).Return("", errors.New("user already exists"))
			},
			expectedStatusCode:  http.StatusConflict,
			expectedRequestBody: `{"error":"user already exists"}`,
		},
		{
			name:                "json_decode_error",
			inputBody:           `{invalid JSON}`,
			mockBehavior:        func(s *mock_handler.MockAuthService, username, password string) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"error":"invalid character 'i' looking for beginning of object key string"}`,
		},
		{
			name:                "empty_request_body",
			inputBody:           ``,
			mockBehavior:        func(s *mock_handler.MockAuthService, username, password string) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"error":"EOF"}`,
		},
		{
			name:                "incorrect_json",
			inputBody:           `{"username": "vika" "password": "123"}`,
			mockBehavior:        func(s *mock_handler.MockAuthService, username, password string) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"error":"invalid character '\"' after object key:value pair"}`,
		},
		{
			name:                "incorrect_json_2",
			inputBody:           `{"username": "vika", "password": 123}`,
			mockBehavior:        func(s *mock_handler.MockAuthService, username, password string) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"error":"json: cannot unmarshal number into Go struct field RegisterRequest.password of type string"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := mock_handler.NewMockAuthService(ctrl)
			validate := validator.New()
			testCase.mockBehavior(auth, testCase.inputUser.Username, testCase.inputUser.Password)

			authHandler := NewAuthHandler(auth, validate)

			r := chi.NewRouter()
			r.Post("/register", authHandler.Register)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register",
				bytes.NewBufferString(testCase.inputBody))

			// Serve
			r.ServeHTTP(w, req)

			// Assert
			actualResponse := strings.TrimSpace(w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, actualResponse)
		})
	}
}
