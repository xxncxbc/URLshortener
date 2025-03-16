package di

import "URLshortener/internal/user"

type IStatRepository interface {
	AddClick(LinkId uint)
}

type IUserRepository interface {
	Create(user *user.User) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
}
