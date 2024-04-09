package repos

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vavelour/chat/internal/domain/entities"
	"github.com/vavelour/chat/internal/repository/inmemorydb/model"
	"github.com/vavelour/chat/internal/repository/inmemorydb/model/constant"
	mock_repos "github.com/vavelour/chat/internal/repository/inmemorydb/repos/mocks"
	"testing"
)

func TestPublicRepos_GetMessages(t *testing.T) {
	type mockBehavior func(m *mock_repos.MockMemoryDB, limit, offset int)

	testTable := []struct {
		name             string
		limit            int
		offset           int
		mockBehavior     mockBehavior
		expectedMessages []entities.Message
		expectedError    error
	}{
		{
			name:   "ok",
			limit:  2,
			offset: 0,
			mockBehavior: func(m *mock_repos.MockMemoryDB, limit, offset int) {
				m.EXPECT().Get(constant.PublicChatKey).Return(model.PublicChat{Messages: []entities.Message{
					{
						Sender:  "tester",
						Content: "hello, world!",
					},
					{
						Sender:  "user",
						Content: "hello, tester!",
					},
				}})
			},
			expectedMessages: []entities.Message{
				{
					Sender:  "tester",
					Content: "hello, world!",
				},
				{
					Sender:  "user",
					Content: "hello, tester!",
				},
			},
			expectedError: nil,
		},
		{
			name:   "incorrect_type",
			limit:  1,
			offset: 0,
			mockBehavior: func(m *mock_repos.MockMemoryDB, limit, offset int) {
				m.EXPECT().Get(constant.PublicChatKey).Return(errIncorrectType)
			},
			expectedMessages: nil,
			expectedError:    errIncorrectType,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_repos.NewMockMemoryDB(ctrl)
			repo := NewPublicRepos(mockDB)

			testCase.mockBehavior(mockDB, testCase.limit, testCase.offset)

			messages, err := repo.GetMessages(testCase.limit, testCase.offset)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedMessages, messages)
		})
	}
}

func TestPublicRepos_InsertMessage(t *testing.T) {
	type mockBehavior func(m *mock_repos.MockMemoryDB, mess entities.Message)

	testTable := []struct {
		name            string
		expectedMessage entities.Message
		mockBehavior    mockBehavior
		expectedError   error
	}{
		{
			name:            "ok",
			expectedMessage: entities.Message{Sender: "tester", Content: "hello, world!"},
			mockBehavior: func(m *mock_repos.MockMemoryDB, mess entities.Message) {
				m.EXPECT().Get(constant.PublicChatKey).Return(model.PublicChat{Messages: []entities.Message{{Sender: "tester", Content: "hello, world!"}}})
				m.EXPECT().Insert(constant.PublicChatKey, gomock.Any()).Do(func(key string, data interface{}) {
					messages, _ := data.(model.PublicChat)
					assert.Equal(t, mess, messages.Messages[0])
				})
			},
			expectedError: nil,
		},
		{
			name:            "incorrect_type",
			expectedMessage: entities.Message{Sender: "tester", Content: "hello, world!"},
			mockBehavior: func(m *mock_repos.MockMemoryDB, mess entities.Message) {
				m.EXPECT().Get(constant.PublicChatKey).Return("invalid type")
			},
			expectedError: errIncorrectType,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_repos.NewMockMemoryDB(ctrl)
			repo := NewPublicRepos(mockDB)

			testCase.mockBehavior(mockDB, testCase.expectedMessage)

			err := repo.InsertMessage(testCase.expectedMessage)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
