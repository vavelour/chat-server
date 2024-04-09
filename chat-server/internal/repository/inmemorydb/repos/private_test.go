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

func TestPrivateRepos_GetMessages(t *testing.T) {
	type mockBehavior func(m *mock_repos.MockMemoryDB, sender, recipient string, limit, offset int)

	testTable := []struct {
		name             string
		sender           string
		recipient        string
		limit            int
		offset           int
		mockBehavior     mockBehavior
		expectedMessages []entities.Message
		expectedError    error
	}{
		{
			name:      "ok",
			sender:    "sender_tester",
			recipient: "recipient_tester",
			limit:     2,
			offset:    0,
			mockBehavior: func(m *mock_repos.MockMemoryDB, sender, recipient string, limit, offset int) {
				m.EXPECT().Get(constant.PrivateChatKey).Return(model.PrivateChatTable{Table: map[model.MembersPrivateChatModel]model.PrivateChat{model.MembersPrivateChatModel{User1: recipient, User2: sender}: {Messages: []entities.Message{{
					Sender:    sender,
					Recipient: recipient,
					Content:   "hello, bro!",
				},
					{
						Sender:    sender,
						Recipient: recipient,
						Content:   "how are you?",
					}}}},
				})
			},
			expectedMessages: []entities.Message{
				{
					Sender:    "sender_tester",
					Recipient: "recipient_tester",
					Content:   "hello, bro!",
				},
				{
					Sender:    "sender_tester",
					Recipient: "recipient_tester",
					Content:   "how are you?",
				},
			},
			expectedError: nil,
		},
		{
			name:      "incorrect_type",
			sender:    "sender_tester",
			recipient: "recipient_tester",
			limit:     2,
			offset:    0,
			mockBehavior: func(m *mock_repos.MockMemoryDB, sender, recipient string, limit, offset int) {
				m.EXPECT().Get(constant.PrivateChatKey).Return(errIncorrectType)
			},
			expectedMessages: nil,
			expectedError:    errIncorrectType,
		},
		{
			name:      "chat_is_not_exists",
			sender:    "sender_tester",
			recipient: "tester",
			limit:     2,
			offset:    0,
			mockBehavior: func(m *mock_repos.MockMemoryDB, sender, recipient string, limit, offset int) {
				m.EXPECT().Get(constant.PrivateChatKey).Return(model.PrivateChatTable{Table: map[model.MembersPrivateChatModel]model.PrivateChat{}})
			},
			expectedMessages: nil,
			expectedError:    ErrChatIsNotExists,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_repos.NewMockMemoryDB(ctrl)
			repo := NewPrivateRepos(mockDB)

			testCase.mockBehavior(mockDB, testCase.sender, testCase.recipient, testCase.limit, testCase.offset)

			messages, err := repo.GetMessages(testCase.sender, testCase.recipient, testCase.limit, testCase.offset)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedMessages, messages)
		})
	}

}

func TestPrivateRepos_InsertMessage(t *testing.T) {
	type mockBehavior func(m *mock_repos.MockMemoryDB, mess entities.Message)

	testTable := []struct {
		name            string
		expectedMessage entities.Message
		mockBehavior    mockBehavior
		expectedError   error
	}{
		{
			name: "ok",
			expectedMessage: entities.Message{
				Sender:    "sender_sender",
				Recipient: "tester",
				Content:   "hello, tester!",
			},
			mockBehavior: func(m *mock_repos.MockMemoryDB, mess entities.Message) {
				m.EXPECT().Get(constant.UsersKey).Return(model.UsersTable{Table: map[string]entities.User{"tester": {Username: "tester", Password: "123"}}})
				m.EXPECT().Get(constant.PrivateChatKey).Return(model.PrivateChatTable{Table: map[model.MembersPrivateChatModel]model.PrivateChat{
					model.MembersPrivateChatModel{User1: "sender_sender", User2: "tester"}: {
						Messages: []entities.Message{{
							Sender:    mess.Sender,
							Recipient: mess.Recipient,
							Content:   mess.Content,
						}}}}})
				m.EXPECT().Insert(constant.PrivateChatKey, gomock.Any()).Do(func(key string, data interface{}) {
					messages, _ := data.(model.PrivateChatTable)
					assert.Equal(t, mess, messages.Table[model.MembersPrivateChatModel{User1: "sender_sender", User2: "tester"}].Messages[0])
				})
			},
			expectedError: nil,
		},
		{
			name:            "user_is_not_exists",
			expectedMessage: entities.Message{},
			mockBehavior: func(m *mock_repos.MockMemoryDB, mess entities.Message) {
				m.EXPECT().Get(constant.UsersKey).Return(model.UsersTable{Table: map[string]entities.User{}})
			},
			expectedError: ErrUserIsNotExists,
		},
		{
			name:            "incorrect_type_user",
			expectedMessage: entities.Message{},
			mockBehavior: func(m *mock_repos.MockMemoryDB, mess entities.Message) {
				m.EXPECT().Get(constant.UsersKey).Return(errIncorrectType)
			},
			expectedError: errIncorrectType,
		},
		{
			name: "incorrect_type_message",
			expectedMessage: entities.Message{
				Sender:    "sender_sender",
				Recipient: "tester",
				Content:   "hello, tester!",
			},
			mockBehavior: func(m *mock_repos.MockMemoryDB, mess entities.Message) {
				m.EXPECT().Get(constant.UsersKey).Return(model.UsersTable{Table: map[string]entities.User{"tester": {Username: "tester", Password: "123"}}})
				m.EXPECT().Get(constant.PrivateChatKey).Return(errIncorrectType)
			},
			expectedError: errIncorrectType,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_repos.NewMockMemoryDB(ctrl)
			repo := NewPrivateRepos(mockDB)

			testCase.mockBehavior(mockDB, testCase.expectedMessage)

			err := repo.InsertMessage(testCase.expectedMessage)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestPrivateRepos_GetUsers(t *testing.T) {
	type mockBehavior func(m *mock_repos.MockMemoryDB, user string)

	testTable := []struct {
		name          string
		user          string
		mockBehavior  mockBehavior
		expectedUsers []string
		expectedError error
	}{
		{
			name: "ok",
			user: "tester",
			mockBehavior: func(m *mock_repos.MockMemoryDB, user string) {
				m.EXPECT().Get(constant.PrivateChatKey).Return(model.PrivateChatTable{Table: map[model.MembersPrivateChatModel]model.PrivateChat{
					model.MembersPrivateChatModel{User1: "tester_1", User2: "tester"}: {Messages: nil},
					model.MembersPrivateChatModel{User1: "tester_2", User2: "tester"}: {Messages: nil},
				}})
			},
			expectedUsers: []string{"tester_1", "tester_2"},
			expectedError: nil,
		},
		{
			name: "incorrect_type",
			user: "tester",
			mockBehavior: func(m *mock_repos.MockMemoryDB, user string) {
				m.EXPECT().Get(constant.PrivateChatKey).Return(errIncorrectType)
			},
			expectedUsers: nil,
			expectedError: errIncorrectType,
		},
		{
			name: "non_users",
			user: "tester",
			mockBehavior: func(m *mock_repos.MockMemoryDB, user string) {
				m.EXPECT().Get(constant.PrivateChatKey).Return(model.PrivateChatTable{Table: map[model.MembersPrivateChatModel]model.PrivateChat{}})
			},
			expectedUsers: nil,
			expectedError: ErrNonUsers,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_repos.NewMockMemoryDB(ctrl)
			repo := NewPrivateRepos(mockDB)

			testCase.mockBehavior(mockDB, testCase.user)

			users, err := repo.GetUsers(testCase.user)
			assert.Equal(t, testCase.expectedUsers, users)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
