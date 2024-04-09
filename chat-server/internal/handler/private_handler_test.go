package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vavelour/chat/internal/domain/entities"
	mock_handler "github.com/vavelour/chat/internal/handler/mocks"
	"github.com/vavelour/chat/internal/handler/request"
	"github.com/vavelour/chat/pkg/pagination"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPrivateHandler_SendPrivateMessage(t *testing.T) {
	type mockBehavior func(s *mock_handler.MockPrivateService, m entities.Message)

	testTable := []struct {
		name                string
		inputBody           string
		inputMessage        request.SendPrivateMessageRequest
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:         "ok",
			inputBody:    `{"content": "hello, world!"}`,
			inputMessage: request.SendPrivateMessageRequest{Sender: "tester", Content: "hello, world!"},
			mockBehavior: func(s *mock_handler.MockPrivateService, m entities.Message) {
				s.EXPECT().SendPrivateMessage(m).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"response":"message sent"}`,
		},
		{
			name:                "empty_message",
			inputBody:           `{"content": ""}`,
			inputMessage:        request.SendPrivateMessageRequest{Sender: "tester", Content: "hello, world!"},
			mockBehavior:        func(s *mock_handler.MockPrivateService, m entities.Message) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"Key: 'SendPrivateMessageRequest.Content' Error:Field validation for 'Content' failed on the 'required' tag"}`,
		},
		{
			name:                "invalid_input",
			inputBody:           `{"content": invalid}`,
			inputMessage:        request.SendPrivateMessageRequest{Sender: "tester", Content: "hello, world!"},
			mockBehavior:        func(s *mock_handler.MockPrivateService, m entities.Message) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:         "send_error",
			inputBody:    `{"content": "hello, world!"}`,
			inputMessage: request.SendPrivateMessageRequest{Sender: "tester", Content: "hello, world!"},
			mockBehavior: func(s *mock_handler.MockPrivateService, m entities.Message) {
				s.EXPECT().SendPrivateMessage(m).Return(errors.New("send error"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"send error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			private := mock_handler.NewMockPrivateService(ctrl)
			validate := validator.New()

			privateHandler := NewPrivateHAndler(private, validate)

			r := chi.NewRouter()
			r.Post("/messages", privateHandler.SendPrivateMessage)

			// Request
			w := httptest.NewRecorder()

			ctx := context.WithValue(context.Background(), "Sender", testCase.inputMessage.Sender)

			req := httptest.NewRequest("POST", "/messages?username=recipient",
				bytes.NewBufferString(testCase.inputBody))
			req = req.WithContext(ctx)

			recipient := req.URL.Query().Get("username")

			testCase.mockBehavior(private, entities.Message{Sender: testCase.inputMessage.Sender, Recipient: recipient, Content: testCase.inputMessage.Content})

			// Serve
			r.ServeHTTP(w, req)

			// Assert
			actualResponse := strings.TrimSpace(w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, actualResponse)
		})
	}
}

func TestPrivateHandler_ShowPrivateMessages(t *testing.T) {
	type mockBehavior func(s *mock_handler.MockPrivateService, sender, recipient string, limit int, offset int)

	testTable := []struct {
		name                string
		inputBody           string
		inputParam          request.ShowPrivateMessageRequest
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:       "ok",
			inputBody:  `{"limit": 1,"offset": 0}`,
			inputParam: request.ShowPrivateMessageRequest{Sender: "tester", Limit: 1, Offset: 0},
			mockBehavior: func(s *mock_handler.MockPrivateService, sender, recipient string, limit int, offset int) {
				s.EXPECT().GetPrivateMessages(sender, recipient, limit, offset).Return([]entities.Message{{Sender: "vika", Recipient: "valera", Content: "hello, world!"}}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"response":"messages received","messages":["hello, world!"]}`,
		},
		{
			name:                "invalid_input",
			inputBody:           `{"limit":}`,
			inputParam:          request.ShowPrivateMessageRequest{Sender: "tester", Limit: 0, Offset: 0},
			mockBehavior:        func(s *mock_handler.MockPrivateService, sender, recipient string, limit int, offset int) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid character '}' looking for beginning of value"}`,
		},
		{
			name:                "invalid_param",
			inputBody:           `{"limit": 0,"offset": 0}`,
			inputParam:          request.ShowPrivateMessageRequest{Sender: "tester", Limit: 1, Offset: 0},
			mockBehavior:        func(s *mock_handler.MockPrivateService, sender, recipient string, limit int, offset int) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"Key: 'ShowPrivateMessageRequest.Limit' Error:Field validation for 'Limit' failed on the 'min' tag"}`,
		},
		{
			name:       "service_error",
			inputBody:  `{"limit": 5,"offset": 0}`,
			inputParam: request.ShowPrivateMessageRequest{Sender: "tester", Limit: 5, Offset: 0},
			mockBehavior: func(s *mock_handler.MockPrivateService, sender, recipient string, limit int, offset int) {
				s.EXPECT().GetPrivateMessages(sender, recipient, limit, offset).Return(nil, errors.New("service error"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"service error"}`,
		},
		{
			name:       "offset_out_of_range",
			inputBody:  `{"limit": 10, "offset": 1000000000}`,
			inputParam: request.ShowPrivateMessageRequest{Sender: "tester", Limit: 10, Offset: 1000000000},
			mockBehavior: func(s *mock_handler.MockPrivateService, sender, recipient string, limit int, offset int) {
				s.EXPECT().GetPrivateMessages(sender, recipient, limit, offset).Return(nil, pagination.ErrOffsetRange)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"response":"no messages found","messages":null}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			private := mock_handler.NewMockPrivateService(ctrl)
			validate := validator.New()

			privateHandler := NewPrivateHAndler(private, validate)

			r := chi.NewRouter()
			r.Get("/messages", privateHandler.ShowPrivateMessages)

			// Request
			w := httptest.NewRecorder()

			ctx := context.WithValue(context.Background(), "Sender", testCase.inputParam.Sender)

			req := httptest.NewRequest("GET", "/messages?username=recipient",
				bytes.NewBufferString(testCase.inputBody))
			req = req.WithContext(ctx)

			recipient := req.URL.Query().Get("username")

			testCase.mockBehavior(private, testCase.inputParam.Sender, recipient, testCase.inputParam.Limit, testCase.inputParam.Offset)

			// Serve
			r.ServeHTTP(w, req)

			// Assert
			actualResponse := strings.TrimSpace(w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, actualResponse)
		})
	}
}

func TestPrivateHandler_ViewUserList(t *testing.T) {
	type mockBehavior func(s *mock_handler.MockPrivateService, user string)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           request.ViewUserListRequest
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "ok",
			inputBody: "",
			inputUser: request.ViewUserListRequest{Username: "tester"},
			mockBehavior: func(s *mock_handler.MockPrivateService, user string) {
				s.EXPECT().ViewUsers(user).Return([]string{"valera"}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"response":"users received","users":["valera"]}`,
		},
		{
			name:      "service_error",
			inputBody: "",
			inputUser: request.ViewUserListRequest{Username: "tester"},
			mockBehavior: func(s *mock_handler.MockPrivateService, user string) {
				s.EXPECT().ViewUsers(user).Return(nil, errors.New("no users who have written to you"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"no users who have written to you"}`,
		},
		{
			name:                "empty_user",
			inputBody:           "",
			inputUser:           request.ViewUserListRequest{Username: ""},
			mockBehavior:        func(s *mock_handler.MockPrivateService, user string) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"Key: 'ViewUserListRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag"}`,
		},
		{
			name:                "failed_get_sender",
			inputBody:           "",
			inputUser:           request.ViewUserListRequest{Username: ""},
			mockBehavior:        func(s *mock_handler.MockPrivateService, user string) {},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"failed to get sender"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			private := mock_handler.NewMockPrivateService(ctrl)
			validate := validator.New()

			privateHandler := NewPrivateHAndler(private, validate)

			r := chi.NewRouter()
			r.Get("/users", privateHandler.ViewUserList)

			// Request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/users",
				bytes.NewBufferString(testCase.inputBody))

			if testCase.name != "failed_get_sender" {
				ctx := context.WithValue(context.Background(), "Sender", testCase.inputUser.Username)
				req = req.WithContext(ctx)
			}

			testCase.mockBehavior(private, testCase.inputUser.Username)

			// Serve
			r.ServeHTTP(w, req)

			// Assert
			actualResponse := strings.TrimSpace(w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, actualResponse)
		})
	}
}
