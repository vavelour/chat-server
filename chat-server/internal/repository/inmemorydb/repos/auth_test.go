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

func TestAuthRepos_GetUser(t *testing.T) {
	type mockBehavior func(m *mock_repos.MockMemoryDB, username string)

	testTable := []struct {
		name          string
		username      string
		mockBehavior  mockBehavior
		expectedUser  entities.User
		expectedError error
	}{
		{
			name:     "ok",
			username: "tester",
			mockBehavior: func(m *mock_repos.MockMemoryDB, username string) {
				m.EXPECT().Get(constant.UsersKey).Return(model.UsersTable{Table: map[string]entities.User{username: {Username: username, Password: "123"}}})
			},
			expectedUser:  entities.User{Username: "tester", Password: "123"},
			expectedError: nil,
		},
		{
			name:     "user_not_found",
			username: "tester",
			mockBehavior: func(m *mock_repos.MockMemoryDB, username string) {
				m.EXPECT().Get(constant.UsersKey).Return(model.UsersTable{Table: map[string]entities.User{}})
			},
			expectedUser:  entities.User{},
			expectedError: errUnregisteredUser,
		},
		{
			name:     "incorrect_type",
			username: "tester",
			mockBehavior: func(m *mock_repos.MockMemoryDB, username string) {
				m.EXPECT().Get(constant.UsersKey).Return("invalid type")
			},
			expectedUser:  entities.User{},
			expectedError: errIncorrectType,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_repos.NewMockMemoryDB(ctrl)
			repo := NewAuthRepos(mockDB)

			testCase.mockBehavior(mockDB, testCase.username)

			user, err := repo.GetUser(testCase.username)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedUser, user)
		})
	}
}

func TestAuthRepos_InsertUser(t *testing.T) {
	type mockBehavior func(m *mock_repos.MockMemoryDB, username, password string)

	testTable := []struct {
		name          string
		username      string
		password      string
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name:     "ok",
			username: "tester",
			password: "123",
			mockBehavior: func(m *mock_repos.MockMemoryDB, username, password string) {
				m.EXPECT().Get(constant.UsersKey).Return(model.UsersTable{Table: map[string]entities.User{}})
				m.EXPECT().Insert(constant.UsersKey, gomock.Any()).Do(func(key string, data interface{}) {
					users, _ := data.(model.UsersTable)
					assert.Equal(t, entities.User{Username: username, Password: password}, users.Table[username])
				})
			},
			expectedError: nil,
		},
		{
			name:     "user_already_exists",
			username: "tester",
			password: "123",
			mockBehavior: func(m *mock_repos.MockMemoryDB, username, password string) {
				m.EXPECT().Get(constant.UsersKey).Return(model.UsersTable{Table: map[string]entities.User{username: {Username: username, Password: password}}})
			},
			expectedError: errUserAlreadyExists,
		},
		{
			name:     "incorrect_type",
			username: "tester",
			password: "123",
			mockBehavior: func(m *mock_repos.MockMemoryDB, username, password string) {
				m.EXPECT().Get(constant.UsersKey).Return("invalid type")
			},
			expectedError: errIncorrectType,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_repos.NewMockMemoryDB(ctrl)
			repo := NewAuthRepos(mockDB)

			testCase.mockBehavior(mockDB, testCase.username, testCase.password)

			err := repo.InsertUser(testCase.username, testCase.password)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
