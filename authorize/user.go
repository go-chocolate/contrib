package authorize

import (
	"context"
)

type User interface {
	GetID() string
	GetSecret() string
	GetUsername() string
	GetPassword() string
	GetClaims() map[string]string
}

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (User, error)
}

type SimpleUser struct {
	ID       string
	Username string
	Secret   string
	Password string
	Claims   map[string]string
}

func (u *SimpleUser) GetID() string                { return u.ID }
func (u *SimpleUser) GetSecret() string            { return u.Secret }
func (u *SimpleUser) GetUsername() string          { return u.Username }
func (u *SimpleUser) GetPassword() string          { return u.Password }
func (u *SimpleUser) GetClaims() map[string]string { return u.Claims }

type SimpleUserRepository struct {
	users map[string]*SimpleUser
}

func NewSimpleUserRepository() *SimpleUserRepository {
	return &SimpleUserRepository{users: make(map[string]*SimpleUser)}
}

func (rep *SimpleUserRepository) GetByUsername(ctx context.Context, username string) (User, error) {
	return rep.users[username], nil
}

func (rep *SimpleUserRepository) Add(u *SimpleUser) {
	rep.users[u.Username] = u
}
