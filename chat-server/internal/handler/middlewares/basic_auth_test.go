package middlewares

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vavelour/chat/internal/domain/entities"
	mock_middlewares "github.com/vavelour/chat/internal/handler/middlewares/mocks"
	"github.com/vavelour/chat/internal/handler/request"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBasicUserIdentity_Identify(t *testing.T) {
	type mockBehavior func(s *mock_middlewares.MockBasicAuthService, u entities.User)

	testTable := []struct {
		name                string
		headerName          string
		headerValue         string
		inputUser           request.BasicAuthLogInRequest
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "ok",
			headerName:  "Authorization",
			headerValue: "Basic dGVzdGVyOjEyMw==",
			inputUser:   request.BasicAuthLogInRequest{Username: "tester", Password: "123"},
			mockBehavior: func(s *mock_middlewares.MockBasicAuthService, u entities.User) {
				s.EXPECT().UserIdentity(u).Return("tester", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: ``,
		},
		{
			name:        "failed_service",
			headerName:  "Authorization",
			headerValue: "Basic dGVzdGVyOjEyMw==",
			inputUser:   request.BasicAuthLogInRequest{Username: "tester", Password: "123"},
			mockBehavior: func(s *mock_middlewares.MockBasicAuthService, u entities.User) {
				s.EXPECT().UserIdentity(u).Return("", errors.New("incorrect login or password"))
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"incorrect login or password"}`,
		},
		{
			name:                "failed_header",
			headerName:          "Nothing",
			headerValue:         "Basic dGVzdGVyOjEyMw==",
			inputUser:           request.BasicAuthLogInRequest{Username: "tester", Password: "123"},
			mockBehavior:        func(s *mock_middlewares.MockBasicAuthService, u entities.User) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"authorization header is missing"}`,
		},
		{
			name:                "failed_base64",
			headerName:          "Authorization",
			headerValue:         "Basic JWT",
			inputUser:           request.BasicAuthLogInRequest{Username: "tester", Password: "123"},
			mockBehavior:        func(s *mock_middlewares.MockBasicAuthService, u entities.User) {},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"illegal base64 data at input byte 0"}`,
		},
		{
			name:                "failed_empty",
			headerName:          "Authorization",
			headerValue:         "Basic",
			inputUser:           request.BasicAuthLogInRequest{Username: "", Password: ""},
			mockBehavior:        func(s *mock_middlewares.MockBasicAuthService, u entities.User) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"invalid authorization header"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			basic := mock_middlewares.NewMockBasicAuthService(ctrl)
			validate := validator.New()
			testCase.mockBehavior(basic, entities.User{Username: testCase.inputUser.Username, Password: testCase.inputUser.Password})

			mwBasic := NewBasicUserIdentity(basic, validate)

			req := httptest.NewRequest(http.MethodGet, "/identity", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			w := httptest.NewRecorder()

			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				username := r.Context().Value("Sender")
				assert.NotNil(t, username)
				w.WriteHeader(http.StatusOK)
			})

			mwBasic.Identify(dummyHandler).ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			actualResponse := strings.TrimSpace(w.Body.String())
			assert.Equal(t, testCase.expectedRequestBody, actualResponse)
		})
	}
}
