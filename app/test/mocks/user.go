package mocks

import (
	"github.com/stretchr/testify/mock"
	"server/app/model"
)

type MockDBUSer struct {
	mock.Mock
}

func (m *MockDBUSer) GetRoleByName(roleName string) (*model.Role, error) {
	args := m.Called(roleName)
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockDBUSer) Create(user *model.User) error {
	user.RoleId = 1
	// Mock
	args := m.Called(user)
	return args.Error(0)
}
