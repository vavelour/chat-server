package middlewares

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_middlewares "github.com/vavelour/chat/internal/handler/middlewares/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJWTUserIdentity_Identify(t *testing.T) {
	type mockBehavior func(s *mock_middlewares.MockJWTBearerService, token string)

	testTable := []struct {
		name                string
		headerName          string
		headerValue         string
		token               string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_middlewares.MockJWTBearerService, token string) {
				s.EXPECT().UserIdentity("token").Return("user", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: ``,
		},
		{
			name:        "invalid_token",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_middlewares.MockJWTBearerService, token string) {
				s.EXPECT().UserIdentity("token").Return("", errors.New("invalid token"))
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"invalid token"}`,
		},
		{
			name:                "empty_header",
			headerName:          "",
			headerValue:         "Bearer token",
			token:               "token",
			mockBehavior:        func(s *mock_middlewares.MockJWTBearerService, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"authorization header is missing"}`,
		},
		{
			name:                "empty_token",
			headerName:          "Authorization",
			headerValue:         "Bearer ",
			token:               "",
			mockBehavior:        func(s *mock_middlewares.MockJWTBearerService, token string) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"login and password cannot be empty"}`,
		},
		{
			name:                "invalid_header",
			headerName:          "Authorization",
			headerValue:         "Basic Auth",
			token:               "token",
			mockBehavior:        func(s *mock_middlewares.MockJWTBearerService, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"invalid authorization header"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			jwt := mock_middlewares.NewMockJWTBearerService(ctrl)
			validate := validator.New()
			testCase.mockBehavior(jwt, testCase.token)

			mwJWT := NewJWTUserIdentity(jwt, validate)

			req := httptest.NewRequest(http.MethodGet, "/identity", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			w := httptest.NewRecorder()

			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				username := r.Context().Value("Sender")
				assert.NotNil(t, username)
				w.WriteHeader(http.StatusOK)
			})

			mwJWT.Identify(dummyHandler).ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			actualResponse := strings.TrimSpace(w.Body.String())
			assert.Equal(t, testCase.expectedRequestBody, actualResponse)
		})
	}
}
