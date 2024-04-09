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

func TestPublicHandler_SendPublicMessage(t *testing.T) {
	type mockBehavior func(s *mock_handler.MockPublicService, m entities.Message)

	testTable := []struct {
		name                string
		inputBody           string
		inputMessage        request.SendPublicMessageRequest
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:         "ok",
			inputBody:    `{"content": "hello, world!"}`,
			inputMessage: request.SendPublicMessageRequest{Sender: "tester", Content: "hello, world!"},
			mockBehavior: func(s *mock_handler.MockPublicService, m entities.Message) {
				s.EXPECT().SendPublicMessage(m).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"response":"message sent"}`,
		},
		{
			name:                "invalid_input",
			inputBody:           `{"content": invalid}`,
			inputMessage:        request.SendPublicMessageRequest{Sender: "tester", Content: ""},
			mockBehavior:        func(s *mock_handler.MockPublicService, m entities.Message) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:                "empty_message",
			inputBody:           `{"content": ""}`,
			inputMessage:        request.SendPublicMessageRequest{Sender: "tester", Content: ""},
			mockBehavior:        func(s *mock_handler.MockPublicService, m entities.Message) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"Key: 'SendPublicMessageRequest.Content' Error:Field validation for 'Content' failed on the 'required' tag"}`,
		},
		{
			name:         "send_error",
			inputBody:    `{"content": "hello, world!"}`,
			inputMessage: request.SendPublicMessageRequest{Sender: "tester", Content: "hello, world!"},
			mockBehavior: func(s *mock_handler.MockPublicService, m entities.Message) {
				s.EXPECT().SendPublicMessage(m).Return(errors.New("send error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"send error"}`,
		},
		{
			name:                "failed_get_sender",
			inputBody:           `{"content": "hello, world!"}`,
			inputMessage:        request.SendPublicMessageRequest{Sender: "tester", Content: "hello, world!"},
			mockBehavior:        func(s *mock_handler.MockPublicService, m entities.Message) {},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"failed to get sender"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			public := mock_handler.NewMockPublicService(ctrl)
			validate := validator.New()
			testCase.mockBehavior(public, entities.Message{Sender: testCase.inputMessage.Sender, Content: testCase.inputMessage.Content})

			publicHandler := NewPublicHandler(public, validate)

			r := chi.NewRouter()
			r.Post("/messages", publicHandler.SendPublicMessage)

			// Request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/messages",
				bytes.NewBufferString(testCase.inputBody))

			if testCase.name != "failed_get_sender" {
				ctx := context.WithValue(context.Background(), "Sender", testCase.inputMessage.Sender)
				req = req.WithContext(ctx)
			}

			// Serve
			r.ServeHTTP(w, req)

			// Assert
			actualResponse := strings.TrimSpace(w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, actualResponse)
		})
	}
}

func TestPublicHandler_ShowPublicMessages(t *testing.T) {
	type mockBehavior func(s *mock_handler.MockPublicService, limit int, offset int)

	testTable := []struct {
		name                string
		inputBody           string
		inputParam          request.ShowPublicMessageRequest
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:       "ok",
			inputBody:  `{"limit": 1,"offset": 0}`,
			inputParam: request.ShowPublicMessageRequest{Limit: 1, Offset: 0},
			mockBehavior: func(s *mock_handler.MockPublicService, limit int, offset int) {
				s.EXPECT().GetPublicMessages(limit, offset).Return([]entities.Message{{Sender: "valera", Content: "hello, world!"}}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"response":"messages received","messages":["hello, world!"]}`,
		},
		{
			name:       "offset_out_of_range",
			inputBody:  `{"limit": 10, "offset": 1000000000}`,
			inputParam: request.ShowPublicMessageRequest{Limit: 10, Offset: 1000000000},
			mockBehavior: func(s *mock_handler.MockPublicService, limit int, offset int) {
				s.EXPECT().GetPublicMessages(limit, offset).Return([]entities.Message{}, pagination.ErrOffsetRange)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"response":"no messages found","messages":null}`,
		},
		{
			name:                "invalid_input",
			inputBody:           `{"limit": "invalid", "offset": "invalid"}`,
			inputParam:          request.ShowPublicMessageRequest{Limit: 0, Offset: 0},
			mockBehavior:        func(s *mock_handler.MockPublicService, limit int, offset int) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"json: cannot unmarshal string into Go struct field ShowPublicMessageRequest.limit of type int"}`,
		},
		{
			name:                "incorrect_input",
			inputBody:           `{"limit": 0,"offset": -1}`,
			inputParam:          request.ShowPublicMessageRequest{Limit: 0, Offset: 0},
			mockBehavior:        func(s *mock_handler.MockPublicService, limit int, offset int) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"Key: 'ShowPublicMessageRequest.Limit' Error:Field validation for 'Limit' failed on the 'min' tag\nKey: 'ShowPublicMessageRequest.Offset' Error:Field validation for 'Offset' failed on the 'min' tag"}`,
		},
		{
			name:       "service_error",
			inputBody:  `{"limit": 10, "offset": 0}`,
			inputParam: request.ShowPublicMessageRequest{Limit: 10, Offset: 0},
			mockBehavior: func(s *mock_handler.MockPublicService, limit int, offset int) {
				s.EXPECT().GetPublicMessages(limit, offset).Return(nil, errors.New("service error"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"service error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			public := mock_handler.NewMockPublicService(ctrl)
			validate := validator.New()
			testCase.mockBehavior(public, testCase.inputParam.Limit, testCase.inputParam.Offset)

			publicHandler := NewPublicHandler(public, validate)

			r := chi.NewRouter()
			r.Get("/messages", publicHandler.ShowPublicMessages)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/messages",
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
