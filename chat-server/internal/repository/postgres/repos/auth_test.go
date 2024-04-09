package repos

import (
	"github.com/golang/mock/gomock"
	"github.com/vavelour/chat/internal/domain/entities"
	mock_repos "github.com/vavelour/chat/internal/repository/postgres/repos/mocks"
	"testing"
)

func TestAuthSqlRepos_GetUser(t *testing.T) {
	type mockBehavior func(m *mock_repos.MockPostgresDB, username string)

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
			mockBehavior: func(m *mock_repos.MockPostgresDB, username string) {
				m.EXPECT().Get(gomock.Any()).Return(nil, nil)
			},
			expectedUser:  entities.User{Username: "tester", Password: "123"},
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {

	}
}

func TestAuthSqlRepos_InsertUser(t *testing.T) {

}
