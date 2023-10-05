package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(username, password string) error
	Authenticate(username, password string) (user *User, authToken *token, err error)
	Validate(token string) (*User, error)
}

type service struct {
	users map[string]*User
}

type User struct {
	ID       uuid.UUID
	Username string
	password string
}

type token struct {
	Value      string
	Expiration time.Time
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

	admin := &User{ID: adminId, Username: "admin", password: string(password)}

	users := make(map[string]*User)
	users[admin.Username] = admin

	return &service{users: users}
}

func (s *service) Register(username, password string) error {
	return nil
}

func (s *service) Authenticate(username, password string) (user *User, authToken *token, err error) {
	u, ok := s.users[username]
	if !ok {
		return nil, nil, fmt.Errorf("incorrect username '%s'", username)
	}

	if bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password)) != nil {
		return nil, nil, fmt.Errorf("incorrect password")
	}

	return u, &token{
		Value:      "token",
		Expiration: time.Now().Add(30 * time.Minute),
	}, nil
}

func (s *service) Validate(token string) (*User, error) {
	if token == "token" {
		u, _ := s.users["admin"]
		return u, nil
	}
	return nil, fmt.Errorf("invalid token '%s'", token)
}
