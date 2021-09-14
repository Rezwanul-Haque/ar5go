package users_test

import (
	"fmt"
	"os"
	"boilerplate/app/domain"
	"boilerplate/app/svc/impl"
	"boilerplate/app/tests/mock/svc/users"
	"boilerplate/infra/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	config.LoadConfig()
	os.Exit(m.Run())
}

func Test_CreateUser(t *testing.T) {
	mockRepo := new(users.MockRepository)

	user := domain.User{UserName: "username", FirstName: "first_name", LastName: "last_name", Email: "user@gmail.com"}

	mockRepo.On("Save").Return(&user)

	testService := impl.NewUsersService(mockRepo, nil, )
	result, err := testService.CreateUser(user)

	mockRepo.AssertExpectations(t)

	fmt.Println(err)
	assert.Equal(t, "username", result.UserName)
	assert.Equal(t, "first_name", result.FirstName)
	assert.Equal(t, "last_name", result.LastName)
	assert.Equal(t, "user@gmail.com", result.Email)
}
