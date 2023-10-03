package auth

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(username, password string) error
	Authenticate(username, password string) (string, error)
}

type service struct {
	users map[uuid.UUID]*user
}

type user struct {
	ID       uuid.UUID
	Username string
	Password string
}

func NewService() Service {
	password, err := bcrypt.GenerateFromPassword([]byte("admin"), 8)
	if err != nil {
		panic("Failed to has admin pass")
	}

	adminId, err := uuid.NewUUID()
	if err != nil {
		panic("Failed to create admin ID")
	}

	admin := &user{ID: adminId, Username: "admin", Password: string(password)}

	users := make(map[uuid.UUID]*user)
	users[adminId] = admin

	return &service{users: users}
}

func (s *service) Register(username, password string) error {
	return nil
}

func (s *service) Authenticate(username, password string) (string, error) {
	return "", nil
}
